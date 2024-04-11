package statev1

import (
	"context"

	"connectrpc.com/connect"
	statev1 "github.com/ophum/tfstate-manager/gen/api/state/v1"
	"github.com/ophum/tfstate-manager/pkg/middlewares"
	"github.com/ophum/tfstate-manager/pkg/models"
)

func (s *StateServer) List(
	ctx context.Context,
	req *connect.Request[statev1.ListRequest],
) (*connect.Response[statev1.ListResponse], error) {
	userID, err := middlewares.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	var states []*models.State
	if err := s.db.Where("user_id = ?", userID).
		Find(&states).Error; err != nil {
		return nil, err
	}
	res := statev1.ListResponse{
		States: make([]*statev1.State, 0, len(states)),
	}
	for _, s := range states {
		res.States = append(res.States, &statev1.State{
			Id:          s.ID,
			Name:        s.Name,
			Description: s.Description,
		})
	}

	return &connect.Response[statev1.ListResponse]{
		Msg: &res,
	}, nil
}
