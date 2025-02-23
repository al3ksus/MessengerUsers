package grpcapp

import (
	"fmt"
	"net"

	friendrequestsgrpc "github.com/al3ksus/messengerusers/internal/grpc/friendrequests"
	usersgrpc "github.com/al3ksus/messengerusers/internal/grpc/users"
	"github.com/al3ksus/messengerusers/internal/logger"

	"google.golang.org/grpc"
)

// GRPCServer представляет собой grpc приложение.
type GRPCServer struct {
	log        logger.Logger
	grpcServer *grpc.Server
	port       int
}

// New - контсруктор для типа *GRPCServer.
func New(log logger.Logger, port int, users usersgrpc.Users, friendRequests friendrequestsgrpc.FriendRequests) *GRPCServer {
	grpcServer := grpc.NewServer()
	usersgrpc.Register(grpcServer, users)
	friendrequestsgrpc.Register(grpcServer, friendRequests)
	return &GRPCServer{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

// MustRun точно запускает gprc приложение.
// Паникует в случае ошибки.
func (a *GRPCServer) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

// Run создает tcp соединение по заданному порту.
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

// Stop реализует безопасное завершение работы.
func (a *GRPCServer) Stop() {
	a.log.Infof("stopping grpc server")

	a.grpcServer.GracefulStop()
}
