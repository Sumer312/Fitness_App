package main

import (
	"database/sql"
	"math"
)

func deficit_calc(current_weight int, desired_weight int, time_frame int) sql.NullInt32 {
	var difference int = int(math.Floor(math.Abs(float64(current_weight) - float64(desired_weight))))
	diff_in_kcal := difference * 7716
	deficit := (diff_in_kcal / time_frame) / 7
	return sql.NullInt32{
		Int32: int32(deficit),
		Valid: true,
	}
}
