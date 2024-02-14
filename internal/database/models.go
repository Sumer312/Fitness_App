// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type DailyNutritionIntake struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UserID        uuid.UUID
	Calories      int32
	Carbohydrates int32
	Protien       int32
	Fat           int32
	Fiber         int32
}

type TotalCalorieIntake struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	Calories     int32
	TotalDeficit int32
	TotalSurplus int32
	UserID       uuid.UUID
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     sql.NullString
	Password  string
}

type UserInput struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Program       string
	Sex           string
	Height        int32
	Weight        int32
	DesiredWeight sql.NullInt32
	TimeFrame     sql.NullInt32
	Bmi           float64
	CurrKcal      int32
	Deficit       sql.NullInt32
}
