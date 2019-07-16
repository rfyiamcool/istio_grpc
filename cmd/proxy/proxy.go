package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/rfyiamcool/istio_grpc/constant"
	"github.com/rfyiamcool/istio_grpc/pool"
	"github.com/rfyiamcool/istio_grpc/proto"
)

var (
	debug bool
	port  int

	name = "proxy"
)

func init() {
	flag.IntVar(&port, "port", 33331, "The server listening port")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.Parse()

	if debug {
		pool.SetDebugAddr()
	}
}

type Server struct {
}

func newServer() *Server {
	return &Server{}
}

func (s *Server) GetName(ctx context.Context, req *call.NameReq) (*call.NameResp, error) {
	return &call.NameResp{
		Name: name,
		Host: constant.HostName,
	}, nil
}

func (s *Server) GetDataFromCache(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	client, err := pool.GetCacheClient()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	caller := call.NewCallCacheClient(client)
	return caller.GetData(ctx, req)
}

func (s *Server) SetDataToCache(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	client, err := pool.GetCacheClient()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	caller := call.NewCallCacheClient(client)
	return caller.SetData(ctx, req)
}

func (s *Server) GetDataFromStore(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	client, err := pool.GetStoreClient()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	caller := call.NewCallStoreClient(client)
	return caller.GetData(ctx, req)
}

func (s *Server) SetDataToStore(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	client, err := pool.GetStoreClient()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	caller := call.NewCallStoreClient(client)
	return caller.SetData(ctx, req)
}

func (s *Server) GetData(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	client, err := pool.GetCacheClient()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	caller := call.NewCallCacheClient(client)
	resp, err := caller.GetDataTrace(ctx, req)
	if err != nil {
		log.Printf("GetDataTrace failed, err: %s \n", err.Error())
	}

	return resp, err
}

func (s *Server) SetData(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	// first purge, after set
	client, err := pool.GetCacheClient()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	cacheCaller := call.NewCallCacheClient(client)
	_, err = cacheCaller.DeleteData(ctx, req)
	if err != nil {
		log.Printf("delete cache failed, err: %s \n", err.Error())
	}

	storeClient, _ := pool.GetStoreClient()
	storeCaller := call.NewCallStoreClient(storeClient)
	resp, err := storeCaller.SetData(ctx, req)
	if err != nil {
		log.Printf("set data to store failed, err: %s \n", err.Error())
	}

	return resp, err
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var (
		s   = grpc.NewServer()
		srv = newServer()
	)

	// register
	call.RegisterCallProxyServer(s, srv)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down gRPC server...")
			s.GracefulStop()
		}
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
