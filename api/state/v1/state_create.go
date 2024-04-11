package statev1

import (
	"context"

	"connectrpc.com/connect"
	statev1 "github.com/ophum/tfstate-manager/gen/api/state/v1"
	"github.com/ophum/tfstate-manager/pkg/middlewares"
	"github.com/ophum/tfstate-manager/pkg/models"
)

func (s *StateServer) Create(
	ctx context.Context,
	req *connect.Request[statev1.CreateRequest],
) (*connect.Response[statev1.CreateResponse], error) {
	userID, err := middlewares.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	state := models.State{
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
		State:       "",
		UserID:      userID,
	}

	if err := s.db.Create(&state).Error; err != nil {
		return nil, err
	}

	return &connect.Response[statev1.CreateResponse]{
		Msg: &statev1.CreateResponse{
			Id:          state.ID,
			Name:        state.Name,
			Description: state.Description,
		},
	}, nil
}
