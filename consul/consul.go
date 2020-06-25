package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

var DefaultClient *Client

func InitConsul(address string) (err error) {
	conf := DefaultConfig()
	if address != "" {
		conf.Address = address
	}
	DefaultClient, err = NewClient(conf)
	return err
}

// Config 配置
type Config struct {
	Address                string        // Address is the address of the Consul server
	Scheme                 string        // Scheme is the URI scheme for the Consul server
	Datacenter             string        // Datacenter to use. If not provided, the default agent datacenter is used.
	CheckInterval          string        // Interval 间隔多长时间请求一次，例如："3s"
	CheckTimeout           string        // Timeout 请求超时时间，例如："5s"
	CheckDeregisterTimeout string        // DeregisterTimeout 请求超时多长时间后将服务取消注册，例如："30s"
	WaitTime               time.Duration // WaitTime 客户端监听服务地址等待时长
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Address:                "127.0.0.1:8500",
		Scheme:                 "http",
		Datacenter:             "dc1",
		CheckInterval:          "3s",
		CheckTimeout:           "5s",
		CheckDeregisterTimeout: "30s",
		WaitTime:               10 * time.Minute,
	}
}

// Client 是 consul 客户端
type Client struct {
	*api.Client
	config *Config
}

// NewClient 创建客户端
func NewClient(config *Config) (*Client, error) {
	consulConfig := &api.Config{
		Address:    config.Address,
		Scheme:     config.Scheme,
		Datacenter: config.Datacenter,
	}

	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client, config: config}, nil
}

// ServiceRegister 注册一个服务
func (c *Client) ServiceRegister(name, ip string, port int, tags ...string) error {
	registration := &api.AgentServiceRegistration{
		Name:    name,
		ID:      fmt.Sprintf("%s-%s:%d", name, ip, port),
		Address: ip,
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%d", ip, port),
			Interval:                       c.config.CheckInterval,
			Timeout:                        c.config.CheckTimeout,
			DeregisterCriticalServiceAfter: c.config.CheckDeregisterTimeout,
		},
	}

	return c.Agent().ServiceRegister(registration)
}

// ServiceDeregister 注销一个服务
func (c *Client) ServiceDeregister(id string) error {
	return c.Agent().ServiceDeregister(id)
}

// ServiceAddress 返回一个服务的地址列表
func (c *Client) ServiceAddress(name string) ([]string, uint64, error) {
	entries, meta, err := c.Health().Service(name, "", true, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, 0, err
	}

	return convertEntries(entries), meta.LastIndex, nil
}

// ServiceAddressWatch 阻塞监听一个服务的地址列表
func (c *Client) ServiceAddressWatch(name string, waitIndex uint64) ([]string, uint64, error) {
	entries, meta, err := c.Client.Health().Service(name, "", true,
		&api.QueryOptions{RequireConsistent: true, WaitIndex: waitIndex, WaitTime: c.config.WaitTime})
	if err != nil {
		return nil, waitIndex, err
	}

	return convertEntries(entries), meta.LastIndex, nil
}

func convertEntries(entries []*api.ServiceEntry) []string {
	var list []string
	for _, e := range entries {
		list = append(list, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return list
}
