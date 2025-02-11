package grpcapp

import (
	"fmt"
	"log"
	"net"

	usersgrpc "github.com/al3ksus/messengerusers/internal/grpc/users"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	log        *log.Logger
	grpcServer *grpc.Server
	port       int
}

func New(log *log.Logger, port int) *GRPCServer {
	grpcServer := grpc.NewServer()
	usersgrpc.Register(grpcServer)
	return &GRPCServer{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *GRPCServer) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *GRPCServer) Run() error {
	log.Print("starting grpc server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", "grpcapp.Run", err)
	}

	log.Printf("grpc server is running. addr=%s", l.Addr().String())

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", "grpcapp.Run", err)
	}

	return nil
}

func (a *GRPCServer) Stop() {
	const op = "grpcapp.Stop"

	a.log.Printf("stopping grpc server. op: %s, port:%d", op, a.port)

	a.grpcServer.GracefulStop()
}
