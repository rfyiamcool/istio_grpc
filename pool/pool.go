package pool

import (
	"log"

	"google.golang.org/grpc"
)

var (
	cacheClient *grpc.ClientConn
	proxyClient *grpc.ClientConn
	storeClient *grpc.ClientConn
)

var (
	ProxyServerAddr = "proxy:33331"
	CacheServerAddr = "cache:33332"
	StoreServerAddr = "store:33333"
)

func SetDebugAddr() {
	ProxyServerAddr = "127.0.0.1:33331"
	CacheServerAddr = "127.0.0.1:33332"
	StoreServerAddr = "127.0.0.1:33333"
}

func GetProxyClient() (*grpc.ClientConn, error) {
	if proxyClient != nil {
		return proxyClient, nil
	}

	client, err := grpc.Dial(ProxyServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	proxyClient = client
	return client, nil
}

func GetCacheClient() (*grpc.ClientConn, error) {
	if cacheClient != nil {
		return cacheClient, nil
	}

	client, err := grpc.Dial(CacheServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	cacheClient = client
	return client, nil
}

func GetStoreClient() (*grpc.ClientConn, error) {
	if storeClient != nil {
		return storeClient, nil
	}

	client, err := grpc.Dial(StoreServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	storeClient = client
	return client, nil
}
