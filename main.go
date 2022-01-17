package main

import (
	"github.com/rs/zerolog/log"
	"grpc-gateway-example/config"
	"grpc-gateway-example/handler"
	"grpc-gateway-example/server"
	"os"
	"os/signal"
	"syscall"
)

func run() {
	cfg := config.Init()

	h, err := handler.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make new handler")
	}
	grpcServer, err := server.NewGRPCServer(h)
	if err != nil {
		log.Fatal().Err(err).Msg("NewGRPCServer failed")
	}

	go grpcServer.ServeGRPC()
	go server.ServeHTTP()

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
