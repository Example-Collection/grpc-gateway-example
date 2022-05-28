package server

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gw "grpc-gateway-example/proto"
	"net/http"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8080", "gRPC server endpoint")
)

func ServeHTTP() {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := gw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to perform RegisterUserServiceHandlerFromEndpoint()")
	}

	if err := http.ListenAndServe("localhost:8081", mux); err != nil {
		log.Fatal().Err(err).Msg("ServeHTTP() failed.")
	}
}
