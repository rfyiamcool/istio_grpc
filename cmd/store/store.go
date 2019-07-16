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
	"github.com/rfyiamcool/istio_grpc/proto"
)

var (
	debug bool
	port  int

	name = "store"
)

func init() {
	flag.IntVar(&port, "port", 33333, "The server listening port")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.Parse()
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
	s.RLock()
	defer s.RUnlock()

	value, ok := s.store[req.Key]
	if ok {
		return &call.DataModel{
			Key:    req.Key,
			Value:  value,
			Status: constant.OK,
			Where:  call.WhereType_STORE,
		}, nil
	}

	return &call.DataModel{
		Status: constant.Error,
		Where:  call.WhereType_STORE,
	}, nil
}

func (s *Server) SetData(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	s.Lock()
	defer s.Unlock()

	s.store[req.Key] = req.Value
	log.Println("set key to store ok")

	return &call.DataModel{
		Key:    "",
		Value:  "",
		Status: constant.OK,
		Where:  call.WhereType_STORE,
	}, nil
}

func (s *Server) DeleteData(ctx context.Context, req *call.DataModel) (*call.DataModel, error) {
	s.Lock()
	defer s.Unlock()

	delete(s.store, req.Key)
	return &call.DataModel{
		Key:    "",
		Value:  "",
		Status: constant.OK,
		Where:  call.WhereType_STORE,
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
	call.RegisterCallStoreServer(s, srv)

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
