package service

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	savedUser, err := service.DB.CreateUser(ctx, user.New())
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

func (service *Service) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	if strings.Trim(userId, " ") == "" {
		return nil, errors.New("empty userId")
	}
	user, err := service.DB.GetUserByID(ctx, userId)
	if err != nil {
		return nil, ErrUserNotFoundById
	}
	return user, nil
}

func (service *Service) GetUsersByNickname(ctx context.Context, name string, page int64, size int64, sort string) ([]*model.User, error) {
	if strings.Trim(name, " ") == "" {
		return nil, errors.New("empty name")
	}
	users, err := service.DB.GetUsersByNickname(ctx, name, sort, page, size)
	if err != nil {
		if errors.Is(err, userdb.ErrWrongSortValue) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid sort value %s", sort)
		}
		return nil, status.Errorf(codes.Internal, "Error occurred in while processing GetUsers(), %v", err)
	}
	return users, nil
}
