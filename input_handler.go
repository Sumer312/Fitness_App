package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sumer312/Health-App-Backend/internal/database"
)

func (apiCfg *apiConfig) input_handler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserId         uuid.UUID `json:"user_id"`
		Height         int       `json:"height"`
		Weight         int       `json:"weight"`
		Desired_Weight *int      `json:"desired_weight"`
		TimeFrame      *int      `json:"time-frame"`
		Bmi            int       `json:"bmi"`
		Program        string    `json:"program"`
		Curr_Kcal      int       `json:"curr_kcal"`
	}
	var DesiredWeightIsEmpty bool = false
	var TimeFrameIsEmpty bool = false
	r.ParseForm()
	userId, err := uuid.Parse(r.FormValue("userId"))
	if err != nil {
		log.Fatalln(err)
	}
	height, err := strconv.ParseInt(r.FormValue("height"), 10, 32)
	if err != nil {
		log.Fatalln(err)
	}
	weight, err := strconv.ParseInt(r.FormValue("weight"), 10, 32)
	if err != nil {
		log.Fatalln(err)
	}
	desiredWeight, err := strconv.ParseInt(r.FormValue("desired_weight"), 10, 32)
	if err != nil {
		if err != strconv.ErrSyntax {
			log.Fatalln(err)
		}
		DesiredWeightIsEmpty = true
	}
	timeFrame, err := strconv.ParseInt(r.FormValue("time_frame"), 10, 32)
	if err != nil {
		if err != strconv.ErrSyntax {
			log.Fatalln(err)
		}
		TimeFrameIsEmpty = true
	}
	bmi, err := strconv.ParseInt(r.FormValue("bmi"), 10, 32)
	if err != nil {
		log.Fatalln(err)
	}
	currKcal, err := strconv.ParseInt(r.FormValue("curr_kcal"), 10, 32)
	if err != nil {
		log.Fatalln(err)
	}
	program := r.FormValue("program")
	if DesiredWeightIsEmpty == true && TimeFrameIsEmpty == true {
		deficit := deficit_calc(int(weight), int(desiredWeight), int(timeFrame))
		apiCfg.DB.CreateUserInput(r.Context(), database.CreateUserInputParams{
			ID:            uuid.New(),
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
			UserID:        userId,
			Height:        int32(height),
			Weight:        int32(weight),
			DesiredWeight: sql.NullInt32{Int32: int32(desiredWeight), Valid: true},
			TimeFrame:     sql.NullInt32{Int32: int32(timeFrame), Valid: true},
			CurrKcal:      int32(currKcal),
			Bmi:           int32(bmi),
			Program:       program,
			Deficit:       deficit,
		})
	} else {
		apiCfg.DB.CreateUserInput(r.Context(), database.CreateUserInputParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    userId,
			Height:    int32(height),
			Weight:    int32(weight),
			CurrKcal:  int32(currKcal),
			Bmi:       int32(bmi),
			Program:   program,
		})
	}
}
