package controllers

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/sumer312/Health-App-Backend/internal/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (apiCfg *Api) InputHandler(w http.ResponseWriter, r *http.Request) {
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
	r.ParseForm()
	DesiredWeightIsEmpty := false
	TimeFrameIsEmpty := false
	for k, vs := range r.Form {
		for _, v := range vs {
			fmt.Printf("%s => %s\n", k, v)
		}
	}
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		log.Fatalln(err)
	}
	userId, err := uuid.Parse(cookieVal.Value)
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
		if r.Form.Has("desired_weight") {
			log.Fatalln(err)
		}
		DesiredWeightIsEmpty = true
	}
	timeFrame, err := strconv.ParseInt(r.FormValue("time_frame"), 10, 32)
	if err != nil {
		if r.Form.Has("time_frame") {
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
	if DesiredWeightIsEmpty == true {
		if TimeFrameIsEmpty == true {
			TempChan := make(chan sql.NullInt32)
			go func(w int, dw int, tf int) {
				TempChan <- DeficitCalc(w, dw, tf)
			}(int(weight), int(desiredWeight), int(timeFrame))
			deficit := <-TempChan
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
				TimeFrame: sql.NullInt32{Int32: int32(timeFrame), Valid: true},
				CurrKcal:  int32(currKcal),
				Bmi:       int32(bmi),
				Program:   program,
			})
		}
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
	w.Header().Add("Hx-Redirect", "http://localhost:5000")
  w.WriteHeader(200)
}
