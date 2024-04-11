package statev1

import (
	"github.com/ophum/tfstate-manager/gen/api/state/v1/statev1connect"
	"gorm.io/gorm"
)

type StateServer struct {
	db *gorm.DB
}

var _ statev1connect.StateServiceHandler = (*StateServer)(nil)

func NewStateServer(db *gorm.DB) *StateServer {
	return &StateServer{
		db: db,
	}
}
