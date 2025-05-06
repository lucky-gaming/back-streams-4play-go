package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Status string

const (
	StatusScheduled Status = "scheduled"
	StatusInProgress Status = "in_progress"
	StatusExpired    Status = "expired"
	StatusCancelled  Status = "cancelled"
)

type Live struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Banner      string             `json:"banner" bson:"banner"`
	Streamer    string             `json:"streamer" bson:"streamer"`
	StartDate   int64              `json:"startDate" bson:"startDate"` // UNIX timestamp (seconds)
	EndDate     int64              `json:"endDate" bson:"endDate"`     // UNIX timestamp (seconds)
	Status      Status             `json:"status" bson:"status"`
	StreamKey   string             `json:"-" bson:"streamKey"`
}
