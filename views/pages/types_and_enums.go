package pages

import (
	"time"

	"github.com/google/uuid"
)


type TrackProgress struct {
	Id        uuid.UUID
	CreatedAt time.Time
  WeightProgress float32
  TimeFrameProgress float32
}
type DailyLogs struct {
	Id        uuid.UUID
	CreatedAt time.Time
	Calories  float32
	Carbs     float32
	Protien   float32
	Fat       float32
	Fiber     float32
}
