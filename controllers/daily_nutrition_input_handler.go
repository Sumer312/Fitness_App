package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sumer312/Health-App-Backend/internal/database"
)

func (apicfg *Api) CalorieInputHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		log.Println(err)
	}
	userId, err := uuid.Parse(cookieVal.Value)
	if err != nil {
		log.Println(err)
	}
	kCal, err := strconv.ParseInt(r.FormValue("calories"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	carbs, err := strconv.ParseInt(r.FormValue("carbohydrates"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	protien, err := strconv.ParseInt(r.FormValue("protien"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	fat, err := strconv.ParseInt(r.FormValue("fat"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	fiber, err := strconv.ParseInt(r.FormValue("fiber"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	program := r.FormValue("program")

	_, daily_create_err := apicfg.DB.CreateDailyNutrition(r.Context(), database.CreateDailyNutritionParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		UserID:        userId,
		Calories:      int32(kCal),
		Carbohydrates: int32(carbs),
		Protien:       int32(protien),
		Fat:           int32(fat),
		Fiber:         int32(fiber),
	})
	if err != nil {
		log.Println(daily_create_err)
	}
	var most_recent database.TotalCalorieIntake
	user_total, err := apicfg.DB.GetMostRecentUserKcal(r.Context(), userId)
	if err == sql.ErrNoRows {
		user_create_total, err := apicfg.DB.CreateTotalCalorieIntake(r.Context(), database.CreateTotalCalorieIntakeParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    userId,
			Calories:  0,
			Program:   program,
		})
		fmt.Println(time.Now().Unix())
		if err != nil {
			log.Println(err)
		} else {
			most_recent = user_create_total
			fmt.Println(user_create_total.CreatedAt.Unix())
		}
	} else if err != nil {
		log.Println(err)
	} else {
		most_recent = user_total
		fmt.Println("hi from line 83")
	}

	fmt.Println(time.Now().Unix())
	fmt.Println(most_recent.CreatedAt.Unix())
	if time.Now().Unix() >= most_recent.CreatedAt.Unix()+24 {
		fmt.Println("hi from line 89")
		total_kcal := 0
		user_daily, err := apicfg.DB.GetDailyNutritionOfUserByUserId(r.Context(), userId)
		if err != nil {
      log.Println(err)
		}
		for i := 0; i < len(user_daily); i++ {
			total_kcal += int(user_daily[i].Calories)
		}
		apicfg.DB.CreateTotalCalorieIntake(r.Context(), database.CreateTotalCalorieIntakeParams{
			ID:        uuid.New(),
			Calories:  int32(total_kcal),
			UserID:    userId,
			Program:   program,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now()})
		apicfg.DB.DeleteDailyNutritionOfUserByUserId(r.Context(), userId)
		apicfg.DB.DeleteRedundantRows(r.Context())
	}
}
