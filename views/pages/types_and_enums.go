package pages

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"os"
	"time"
)

var err error = godotenv.Load()
var base_url string = os.Getenv("BASE_URL")

const (
	program_fatLoss    = "fatloss"
	program_muscleGain = "musclegain"
	program_maintain   = "maintenance"
	sex_male           = "male"
	sex_female         = "female"
	sex_none           = "none"
)

type TrackProgress struct {
	Id                uuid.UUID
	CreatedAt         time.Time
	WeightProgress    float64
	TimeFrameProgress float64
	ProgramSelected   bool
	Program           string
	ProgramDisplay    string
}
type DailyLogs struct {
	Id        uuid.UUID
	CreatedAt time.Time
	Calories  float32
	Carbs     float32
	Protein   float32
	Fat       float32
	Fiber     float32
}
