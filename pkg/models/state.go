package models

type State struct {
	Model
	Name        string
	Description string
	State       string `gorm:"type:MEDIUMTEXT"`
	UserID      uint64

	User *User
}
