package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"grpc-gateway-example/config"
	"grpc-gateway-example/model"
	"grpc-gateway-example/userdb"
	"strings"
)

type Service struct {
	config config.DatabaseConfig
	DB     *userdb.DB
}

func (service *Service) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := service.validateUser(user); err != nil {
		return nil, errors.Wrapf(err, "validateUser failed: name=%s, nickname=%s", user.Name, user.Nickname)
	}
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrapf(err, "uuid.NewUUID() failed. %v", err)
	}
	user.ID = userID.String()
	savedUser, err := service.DB.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "DB.CreateUser() failed")
	}
	return savedUser, nil
}

func (service *Service) validateUser(user *model.User) error {
	if strings.Trim(user.Nickname, " ") == "" {
		return ErrNicknameEmpty
	}
	if strings.Trim(user.Name, " ") == "" {
		return ErrNameEmpty
	}
	return nil
}

func New(cfg config.DatabaseConfig, db *userdb.DB) *Service {
	return &Service{
		config: cfg,
		DB:     db,
	}
}
