package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attendance struct {
	ID              primitive.ObjectID `bson:"_id"`
	Event_id        string             `json:"event_id"`
	User_id         string             `json:"user_id"`
	Attendance_id   string             `json:"enroll_id"`
	Attendance_date time.Time          `json:"attendance_date"`
	Created_at      time.Time          `json:"created_at"`
	Updated_at      time.Time          `json:"updated_at"`
}
