package main

import (
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"grpc-gateway-example/config"
	"grpc-gateway-example/handler"
	api "grpc-gateway-example/proto"
	"net"
	"time"
)

func main() {
	cfg := config.Init()
	h, err := handler.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make new handler")
	}
	grpcServer, err := NewGRPCServer(h, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("NewGRPCServer failed")
	}
	grpcServer.ServeGRPC()
}

type GRPCServer struct {
	Server   *grpc.Server
	listener net.Listener
}

func NewGRPCServer(h *handler.Handler, cfg config.DatabaseConfig) (*GRPCServer, error) {
	server := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle:     15 * time.Second,
				MaxConnectionAge:      30 * time.Second,
				MaxConnectionAgeGrace: 15 * time.Second,
				Time:                  15 * time.Second,
				Timeout:               10 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             3 * time.Second,
				PermitWithoutStream: true,
			},
		),
		grpcmiddleware.WithUnaryServerChain(
			grpcrecovery.UnaryServerInterceptor(
				grpcrecovery.WithRecoveryHandler(handleRecoveryGRPC),
			),
		),
	)

	api.RegisterUserServiceServer(server, h)

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}

	return &GRPCServer{
		Server:   server,
		listener: lis,
	}, nil
}

func handleRecoveryGRPC(p interface{}) error {
	if err, ok := p.(error); ok {
		log.Err(err)
		return err
	}
	return nil
}

func (server *GRPCServer) ServeGRPC() {
	if err := server.Server.Serve(server.listener); err != nil {
		log.Fatal().Err(err).Msg("GRPCServer.ServeGRPC() failed")
	}
}
