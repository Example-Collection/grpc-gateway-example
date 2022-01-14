package handler

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-gateway-example/model"
	api "grpc-gateway-example/proto"
	"grpc-gateway-example/service"
)

func (handler *Handler) SaveUser(ctx context.Context, request *api.SaveUserRequest) (*api.UserResponse, error) {
	user := handler.convertSaveUserRequestToUserModel(request)
	user, err := handler.UserService.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, service.ErrNicknameEmpty) || errors.Is(err, service.ErrNameEmpty) {
			return nil, status.Errorf(codes.InvalidArgument, "failed to save user(invalid argument)")
		}
		return nil, status.Errorf(codes.Internal, "failed to save user(internal error)")
	}
	return handler.convertUserToSaveUserResponse(user), nil
}

func (handler *Handler) convertSaveUserRequestToUserModel(request *api.SaveUserRequest) *model.User {
	return &model.User{
		Nickname: request.Nickname,
		Name:     request.Name,
	}
}

func (handler *Handler) convertUserToSaveUserResponse(user *model.User) *api.UserResponse {
	return &api.UserResponse{
		Id:       user.ID,
		Name:     user.Name,
		Nickname: user.Nickname,
	}
}
