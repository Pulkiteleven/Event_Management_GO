package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	User_name     string             `json:"user_name" validate:"required"`
	Avatar        string             `json:"avatar"`
	City          string             `json:"city"`
	Email         string             `json:"email" validate:"required"`
	Categories    *[]string           `json:"categories"`
	Password      *string            `json:"password" validate:"required"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}
