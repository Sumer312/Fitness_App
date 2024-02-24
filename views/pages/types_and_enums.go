package pages

import (
	"time"

	"github.com/google/uuid"
)

const (
	program_fatLoss    = "fatloss"
	program_muscleGain = "musclegain"
	program_maintain   = "maintaince"
	sex_male           = "male"
	sex_female         = "female"
	sex_none           = "none"
)

type TrackProgress struct {
	Id        uuid.UUID
	CreatedAt time.Time
  WeightProgress float64
  TimeFrameProgress float64
  ProgramSelected bool
  Program string
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
