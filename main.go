package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-gateway-example/config"
	"grpc-gateway-example/handler"
	gw "grpc-gateway-example/proto"
	"grpc-gateway-example/server"
	"os"
	"os/signal"
	"syscall"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8080", "gRPC server endpoint")
)

func run() {
	cfg := config.Init()
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := gw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to perform RegisterUserServiceHandlerFromEndpoint()")
	}

	h, err := handler.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make new handler")
	}

	grpcServer, err := server.NewGRPCServer(h)
	if err != nil {
		log.Fatal().Err(err).Msg("NewGRPCServer failed")
	}

	go grpcServer.ServeGRPC()
	go server.ServeHTTP(mux)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	go stop(quit, done)
	log.Info().Msg("starting server...")
	<-done
}

func stop(quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.Info().Msg("stopping server...")
	close(done)
}

func main() {
	run()
}
