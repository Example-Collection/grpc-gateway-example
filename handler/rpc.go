package handler

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	return handler.convertUserToUserResponse(user), nil
}

func (handler *Handler) convertSaveUserRequestToUserModel(request *api.SaveUserRequest) *model.User {
	return &model.User{
		Nickname: request.Nickname,
		Name:     request.Name,
	}
}

func (handler *Handler) convertUserToUserResponse(user *model.User) *api.UserResponse {
	return &api.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Nickname:  user.Nickname,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func (handler *Handler) GetUser(_ context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
	user, err := handler.UserService.GetUserByID(request.Id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFoundById) {
			return nil, status.Errorf(codes.NotFound, "user not found(user_id: %s)", request.Id)
		} else {
			return nil, status.Errorf(codes.Internal, "failed to get user by id(user_id: %s)", request.Id)
		}
	}
	return &api.GetUserResponse{
		User:    handler.convertUserToUserResponse(user),
		Message: fmt.Sprintf("Successfully found user with user_id %s", request.Id),
	}, nil
}

func (handler *Handler) GetUsers(_ context.Context, request *api.GetUsersRequest) (*api.GetUsersResponse, error) {
	users, err := handler.UserService.GetUsersByNickname(request.Nickname, request.Page, request.Size, request.Sort)
	if err != nil {
		return nil, err
	}
	return handler.convertUsersToGetUsersResponse(
			users,
			fmt.Sprintf("Successfully found all users with nickname %s", request.Nickname)),
		nil
}

func (handler *Handler) convertUsersToGetUsersResponse(users []*model.User, message string) *api.GetUsersResponse {
	var userResponse = make([]*api.UserResponse, 0)
	for _, user := range users {
		userResponse = append(userResponse, handler.convertUserToUserResponse(user))
	}
	return &api.GetUsersResponse{
		Users:   userResponse,
		Message: message,
	}

}
