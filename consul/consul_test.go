package consul_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"yimcom/consul"
	"yimcom/rpcserver"

	"github.com/stretchr/testify/require"
)

func TestConsul(t *testing.T) {
	serviceName := "TestService"
	servicePort := 8080
	serviceIP, err := rpcserver.LocalAddr()
	require.NoError(t, err)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		err := http.ListenAndServe(fmt.Sprintf(":%d", servicePort), nil)
		require.NoError(t, err)
	}()

	tags := []string{"tag1", "tag2"}
	consulClient, err := consul.NewClient(consul.DefaultConfig())
	require.NoError(t, err)

	err = consulClient.ServiceRegister(serviceName, serviceIP, servicePort, tags...)
	require.NoError(t, err)
	defer func() {
		err = consulClient.ServiceDeregister(fmt.Sprintf("%s-%s:%d", serviceName, serviceIP, servicePort))
		require.NoError(t, err)
		time.Sleep(1 * time.Second) // 等待 3 秒，确保服务已经注册到 consul 上
	}()

	time.Sleep(3 * time.Second) // 等待 3 秒，确保服务已经注册到 consul 上

	addrs, lastIndex, err := consulClient.ServiceAddress(serviceName)
	require.NoError(t, err)
	require.NotEqual(t, 0, len(addrs))

	go func() {
		for {
			_, lastIndex, err = consulClient.ServiceAddressWatch(serviceName, lastIndex)
			fmt.Println(lastIndex)
			require.NoError(t, err)
			require.NotEqual(t, 0, lastIndex)
		}
	}()
}
