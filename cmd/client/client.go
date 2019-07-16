package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/rfyiamcool/istio_grpc/pool"
	"github.com/rfyiamcool/istio_grpc/proto"
)

var (
	debug  bool
	target string
	method string
	name   = "proxy"
)

func init() {
	flag.StringVar(&target, "target", "proxy", "target name")
	flag.StringVar(&method, "method", "GetName", "method name")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	if debug {
		pool.SetDebugAddr()
	}
}

func makeRandomValue() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%v", rand.Int63())
}

func makeKey(id int) string {
	return fmt.Sprintf("key-id-%d", id)
}

func getProxyName() {
	client, err := pool.GetProxyClient()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	caller := call.NewCallProxyClient(client)
	resp, err := caller.GetName(context.TODO(), &call.NameReq{})
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(resp.String())
}

func getCacheName() {
	client, err := pool.GetCacheClient()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	caller := call.NewCallCacheClient(client)
	resp, err := caller.GetName(context.TODO(), &call.NameReq{})
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(resp.String())
}

func getStoreName() {
	client, err := pool.GetStoreClient()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	caller := call.NewCallStoreClient(client)
	resp, err := caller.GetName(context.TODO(), &call.NameReq{})
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(resp.String())
}

func testGetNames() {
	getProxyName()
	getStoreName()
	getCacheName()
}

func testOperateData() {
	client, err := pool.GetProxyClient()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer client.Close()

	caller := call.NewCallProxyClient(client)
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"id": "biss",
	}))

	var (
		key string
		da  *call.DataModel
	)

	for i := 0; i < 10; i++ {
		key = makeKey(i)
		da, err = caller.SetData(ctx, &call.DataModel{
			Key:   key,
			Value: makeRandomValue(),
		})
		if err != nil {
			log.Printf("SetData %v\n", err)
			continue
		}

		log.Printf("resp %v\n", da.String())

		da, err = caller.GetData(ctx, &call.DataModel{
			Key: key,
		})
		if err != nil {
			log.Printf("GetData %v\n", err)
			continue
		}

		log.Println(da.String())
	}
}

func main() {
	testGetNames()
	testOperateData()
}
