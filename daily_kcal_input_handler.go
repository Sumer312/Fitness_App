package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sumer312/Health-App-Backend/internal/database"
)

func (apicfg *apiConfig) calorie_input_handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		UserID   uuid.UUID `json:"user_id"`
		Calories int       `json:"calories"`
		Program  string    `json:"program"`
	}
	params := &parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatal("err parsing json", err)
	}
	_, daily_kcal_create_err := apicfg.DB.CreateDailyCalorieIntake(r.Context(), database.CreateDailyCalorieIntakeParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Calories:  int32(params.Calories),
		UserID:    params.UserID,
	})
	if err != nil {
		log.Fatal(daily_kcal_create_err)
	}
	var most_recent database.TotalCalorieIntake
	user_total, err := apicfg.DB.GetMostRecentUserKcal(r.Context(), params.UserID)
	if err == sql.ErrNoRows {
		user_create_total, err := apicfg.DB.CreateTotalCalorieIntake(r.Context(), database.CreateTotalCalorieIntakeParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    params.UserID,
			Calories:  0,
			Program:   params.Program,
		})
		fmt.Println("hi from line 48 ", time.Now().Unix())
		if err != nil {
			log.Fatal(err)
		} else {
			most_recent = user_create_total
			fmt.Println("hi from line 51", user_create_total.CreatedAt.Unix())
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		most_recent = user_total
		fmt.Println("hi from line 57")
	}

	fmt.Println(time.Now().Unix())
	fmt.Println(most_recent.CreatedAt.Unix())
	if time.Now().Unix() >= most_recent.CreatedAt.Unix()+24 {
		fmt.Println("hi from line 62")
		total_kcal := 0
		user_daily, err := apicfg.DB.GetDailyCalories(r.Context(), params.UserID)
		if err != nil {
			log.Fatal("falied to fetch data", err)
		}
		for i := 0; i < len(user_daily); i++ {
			total_kcal += int(user_daily[i].Calories)
		}
		apicfg.DB.CreateTotalCalorieIntake(r.Context(), database.CreateTotalCalorieIntakeParams{ID: uuid.New(), Calories: int32(total_kcal), UserID: params.UserID, Program: params.Program, CreatedAt: time.Now(), UpdatedAt: time.Now()})
		apicfg.DB.DeleteDailyCalories(r.Context(), params.UserID)
		apicfg.DB.DeleteRedundantRows(r.Context())
	}
}
