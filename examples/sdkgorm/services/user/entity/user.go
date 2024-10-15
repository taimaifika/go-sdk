package entity

import (
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity/enum"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string           `json:"name"`
	Email  string           `json:"email"`
	Gender *enum.GenderEnum `json:"Gender" gorm:"type:gender"`
}
