package server

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"net/http"
)

func ServeHTTP(mux *runtime.ServeMux) {
	if err := http.ListenAndServe("localhost:8081", mux); err != nil {
		log.Fatal().Err(err).Msg("ServeHTTP() failed.")
	}
}
