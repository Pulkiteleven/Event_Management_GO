package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Enroll struct {
	ID         primitive.ObjectID `bson:"_id"`
	Event_id   string             `json:"event_id"`
	User_id    string             `json:"user_id"`
	Approved   bool               `json:"approved"`
	Enroll_id  string             `json:"enroll_id"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}
