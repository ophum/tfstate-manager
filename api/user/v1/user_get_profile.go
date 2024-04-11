package userv1

import (
	"context"

	"connectrpc.com/connect"
	userv1 "github.com/ophum/tfstate-manager/gen/api/user/v1"
	"github.com/ophum/tfstate-manager/pkg/middlewares"
	"github.com/ophum/tfstate-manager/pkg/models"
)

func (s *UserServer) GetProfile(
	ctx context.Context,
	req *connect.Request[userv1.GetProfileRequest],
) (*connect.Response[userv1.GetProfileResponse], error) {
	userID, err := middlewares.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &connect.Response[userv1.GetProfileResponse]{
		Msg: &userv1.GetProfileResponse{
			Data: &userv1.User{
				Id:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
		},
	}, nil
}
