package userv1

import (
	"github.com/ophum/tfstate-manager/gen/api/user/v1/userv1connect"
	"gorm.io/gorm"
)

type UserServer struct {
	db *gorm.DB
}

var _ userv1connect.UserServiceHandler = (*UserServer)(nil)

func NewUserServer(db *gorm.DB) *UserServer {
	return &UserServer{
		db: db,
	}
}
