package grpcapp

import (
	"fmt"
	"net"

	usersgrpc "github.com/al3ksus/messengerusers/internal/grpc/users"
	"github.com/al3ksus/messengerusers/internal/logger"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	log        logger.Logger
	grpcServer *grpc.Server
	port       int
}

func New(log logger.Logger, port int, users usersgrpc.Users) *GRPCServer {
	grpcServer := grpc.NewServer()
	usersgrpc.Register(grpcServer, users)
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
	const op = "grpcapp.Run"
	a.log.Infof("starting grpc server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Infof("grpc server is running. addr=%s", l.Addr().String())

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *GRPCServer) Stop() {
	a.log.Infof("stopping grpc server")

	a.grpcServer.GracefulStop()
}
