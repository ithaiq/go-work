package rpcclient

import (
	"context"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"time"

	"yimcom/consul"
	"yimcom/seelog"

	"github.com/pkg/errors"
	pool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	initialCapacity = 5
	maximumCapacity = 30
	idleTimeout     = time.Hour
	maxLifeDuration = 24 * time.Hour
)

var poolMap struct {
	pools map[uint32]*connPool
	mutex sync.RWMutex
}

func init() {
	log.Println("init yimcom/rpcclient")
	poolMap.pools = make(map[uint32]*connPool)
}

// poolConfig 配置信息
type poolConfig struct {
	InitialCapacity int           `json:"initial_capacity"`  // 初始化连接数
	MaximumCapacity int           `json:"maximum_capacity"`  // 最大容量
	IdleTimeout     time.Duration `json:"idle_timeout"`      // 空闲连接释放超时
	MaxLifeDuration time.Duration `json:"max_life_duration"` // 连接最大生命时长
}

// defaultPoolConfig 返回默认配置信息
func defaultPoolConfig() *poolConfig {
	return &poolConfig{
		InitialCapacity: initialCapacity,
		MaximumCapacity: maximumCapacity,
		IdleTimeout:     idleTimeout,
		MaxLifeDuration: maxLifeDuration,
	}
}

// Acquire 从连接池获取一个连接
func Acquire(bigCmd uint32) (*pool.ClientConn, error) {
	var err error

	poolMap.mutex.RLock()
	p, ok := poolMap.pools[bigCmd]
	poolMap.mutex.RUnlock()

	if !ok || len(p.address) == 0 {
		p, err = NewConnPool(bigCmd)
		if err != nil {
			return nil, errors.Wrap(err, "NewConnPool error")
		}

		poolMap.mutex.Lock()
		poolMap.pools[bigCmd] = p
		poolMap.mutex.Unlock()
	}

	return p.acquire(context.Background())
}

// connPool 连接池
type connPool struct {
	pool    *pool.Pool
	bigCmd  uint32
	address []string
	mutex   sync.RWMutex
}

// Acquire 从连接池获取一个连接
func (p *connPool) acquire(ctx context.Context) (*pool.ClientConn, error) {
get:
	conn, err := p.pool.Get(ctx)
	if err != nil || conn == nil {
		return nil, errors.Wrap(err, "Get error")
	}

	switch conn.GetState() {
	case connectivity.TransientFailure, connectivity.Shutdown:
		conn.Unhealthy()
		// Close will return an error after first time
		// But It is safe to call multiple time, So just ignore error
		seelog.Infof("unhealthy conn %s", conn.GetState().String())
		_ = conn.Close()
		time.Sleep(100 * time.Millisecond)
		goto get
	default:
		return conn, nil
	}
}

// Close 关闭连接池并关闭所有连接
func (p *connPool) Close() {
	p.pool.Close()
}

// NewConnPool 创建客户端连接池
func NewConnPool(bigCmd uint32) (*connPool, error) {
	p := &connPool{bigCmd: bigCmd}
	if err := p.watch(); err != nil {
		return nil, errors.Wrap(err, "watch error")
	}

	factory := func() (*grpc.ClientConn, error) { return p.createConn() }

	grpcpoolConf := defaultPoolConfig()
	gp, err := pool.New(factory, grpcpoolConf.InitialCapacity, grpcpoolConf.MaximumCapacity,
		grpcpoolConf.IdleTimeout, grpcpoolConf.MaxLifeDuration)
	if err != nil {
		return nil, errors.Wrap(err, "New error")
	}

	p.pool = gp

	return p, nil
}

// updateAddress 更新服务IP地址列表
func (p *connPool) updateAddress(address []string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !reflect.DeepEqual(p.address, address) {
		seelog.Infof("address updated (old: %v, new: %v, bigCmd: %d)",
			p.address, address, p.bigCmd)
		p.address = address
	}
}

// createConn 创建客户端连接
func (p *connPool) createConn() (*grpc.ClientConn, error) {
	p.mutex.RLock()
	length := len(p.address)
	if length == 0 {
		return nil, errors.Errorf("bigCmd %d is unavailable", p.bigCmd)
	}
	randAddress := p.address[rand.Intn(length)]
	p.mutex.RUnlock()

	conn, err := grpc.Dial(randAddress, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "Dial error")
	}

	return conn, nil
}

// watch 监听服务地址变动情况
func (p *connPool) watch() error {
	if consul.DefaultClient == nil {
		return errors.New("you need to init consul first")
	}

	defer func() {
		if err := recover(); err != nil {
			seelog.Error("watch panic", err)
		}
	}()

	bigCmdString := strconv.FormatUint(uint64(p.bigCmd), 10)
	address, waitIndex, err := consul.DefaultClient.ServiceAddress(bigCmdString)
	if err != nil {
		return errors.Wrap(err, "ServiceAddress error")
	}
	p.updateAddress(address)

	go func() {
		for {
			latestAddrs, lastIndex, err := consul.DefaultClient.ServiceAddressWatch(bigCmdString, waitIndex)
			if err != nil {
				seelog.Error("bigCmd entries watch error:", err)
				time.Sleep(30 * time.Second)
				continue
			}
			waitIndex = lastIndex
			p.updateAddress(latestAddrs)
		}
	}()

	return nil
}
