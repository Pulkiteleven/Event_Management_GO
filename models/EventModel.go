package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID         primitive.ObjectID `bson:"_id"`
	Owner      string             `json:"owner"`
	Title      string             `json:"title"`
	Desc       string             `json:"desc"`
	Capacity   int                `json:"capacity"`
	Category   string             `json:"category"`
	City       string             `json:"city"`
	Venu       string             `json:"venu"`
	Index      string             `json:"index"`
	Online     bool               `json:"online"`
	Approve    bool               `json:"approve"`
	Link       string             `json:"link"`
	Multiday   bool               `json:"multiday"`
	Start      string             `json:"start"`
	End        string             `json:"end"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Event_id   string             `json:"event_id"`
}
