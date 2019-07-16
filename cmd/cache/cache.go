package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"google.golang.org/grpc"

	"github.com/rfyiamcool/istio_grpc/constant"
	"github.com/rfyiamcool/istio_grpc/pool"
	"github.com/rfyiamcool/istio_grpc/proto"
)

var (
	debug bool
	port  int

	name = "cache"
)

func init() {
	flag.IntVar(&port, "port", 33332, "The server listening port")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.Parse()

	if debug {
		pool.SetDebugAddr()
	}
}

type Server struct {
	sync.RWMutex
	store map[string]string
}

func newServer() *Server {
	return &Server{
		store: make(map[string]string, 100),
	}
}

func (s *Server) GetName(ctx context.Context, req *call.NameReq) (*call.NameResp, error) {
	return &call.NameResp{
		Name: name,
		Host: constant.HostName,
	}, nil
}

func (s *Server) GetData(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	s.Lock()
	defer s.Unlock()

	value, ok := s.store[req.Key]
	if ok {
		return &call.DataModel{
			Key:    req.Key,
			Value:  value,
			Status: constant.OK,
			Where:  call.WhereType_CACHE,
		}, nil
	}

	return &call.DataModel{}, constant.ErrNotFound
}

func (s *Server) GetDataTrace(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	s.Lock()
	value, ok := s.store[req.Key]
	s.Unlock()

	if ok {
		log.Printf("get key ok, key: %s", req.Key)
		return &call.DataModel{
			Key:    req.Key,
			Value:  value,
			Status: constant.OK,
			Where:  call.WhereType_CACHE,
		}, nil
	}

	client, err := pool.GetStoreClient()
	if err != nil {
		log.Println(90)
		return &call.DataModel{}, err
	}

	caller := call.NewCallStoreClient(client)
	resp, err := caller.GetData(ctx, req)
	if err != nil {
		log.Printf("GetData failed, err: %s", err.Error())
		return &call.DataModel{}, err
	}

	s.Lock()
	s.store[resp.Key] = resp.Value
	s.Unlock()

	return resp, nil
}

func (s *Server) SetData(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	s.Lock()
	defer s.Unlock()

	log.Println(111)

	s.store[req.Key] = req.Value
	return &call.DataModel{
		Key:    "",
		Value:  "",
		Status: constant.OK,
		Where:  call.WhereType_CACHE,
	}, nil
}

func (s *Server) DeleteData(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	s.Lock()
	defer s.Unlock()

	delete(s.store, req.Key)
	log.Println("delete key from cache")
	return &call.DataModel{
		Key:    "",
		Value:  "",
		Status: constant.OK,
		Where:  call.WhereType_CACHE,
	}, nil
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
	call.RegisterCallCacheServer(s, srv)

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
