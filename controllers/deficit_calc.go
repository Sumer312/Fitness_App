package controllers

import (
	"database/sql"
	"math"
)

func DeficitCalc(current_weight int, desired_weight int, time_frame int) sql.NullFloat64 {
	var difference int = int(math.Floor(math.Abs(float64(current_weight) - float64(desired_weight))))
	diff_in_kcal := difference * 7716
	deficit := (diff_in_kcal / time_frame) / 7
	return sql.NullFloat64{
		Float64: float64(deficit),
		Valid: true,
	}
}
