package handler

import (
	"grpc-gateway-example/config"
	api "grpc-gateway-example/proto"
	"grpc-gateway-example/service"
	"grpc-gateway-example/userdb"
)

type Handler struct {
	api.UnimplementedUserServiceServer
	UserService *service.Service
}

func New(cfg config.DatabaseConfig) (*Handler, error) {
	userDB, err := userdb.New(cfg)
	if err != nil {
		return nil, err
	}
	userService := service.New(cfg, userDB)
	return &Handler{
		UserService: userService,
	}, nil
}
