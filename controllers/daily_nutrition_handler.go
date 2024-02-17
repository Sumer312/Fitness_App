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
	"github.com/sumer312/Health-App-Backend/views/pages"
)

func (apiCfg *Api) DailyNutritionRender(w http.ResponseWriter, r *http.Request) {
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		log.Println(err)
	}
	userId, err := uuid.Parse(cookieVal.Value)
	if err != nil {
		log.Println(err)
	}
	user, err := apiCfg.DB.GetUserInputById(r.Context(), userId)
	if err != nil {
		log.Println(err)
	}

	user_daily, err := apiCfg.DB.GetDailyNutritionOfUserByUserId(r.Context(), userId)
	if err != nil {
		log.Println(err)
	}
	var total, curr nutritionParams
	for i := 0; i < len(user_daily); i++ {
		curr.carbs += float32(user_daily[i].Carbohydrates)
		curr.protien += float32(user_daily[i].Protien)
		curr.fat += float32(user_daily[i].Fat)
		curr.calories += float32(user_daily[i].Calories)
		curr.fiber += float32(user_daily[i].Fiber)
	}
	if user.Sex == sex_male {
		total.fiber = 31
	} else if user.Sex == sex_female {
		total.fiber = 21
	} else {
		total.fiber = 26
	}

	program := user.Program
	if program == program_fatloss {
		total.protien = float32(user.Weight)
		total.calories = float32(user.CurrKcal - user.Deficit.Int32)
		total.carbs = 0.45 * float32(total.calories/4)
		total.fat = 0.2 * float32(total.calories/9)
	} else if program == program_muscleGain {
		total.protien = 1.2 * float32(user.Weight)
		total.calories = float32(user.CurrKcal) + 200
		total.carbs = 0.4 * float32(total.calories/4)
		total.fat = 0.2 * float32(total.calories/9)
	} else {
		total.protien = 0.8 * float32(user.Weight)
		total.calories = float32(user.CurrKcal)
		total.carbs = 0.6 * float32(total.calories/4)
		total.fat = 0.2 * float32(total.calories/9)
	}
	carbsPercent := fmt.Sprintf("%f", (curr.carbs/total.carbs)*100)
	caloriesPercent := fmt.Sprintf("%f", (curr.calories/total.calories)*100)
	fatPercent := fmt.Sprintf("%f", float32(curr.fat/total.fat)*100)
	protienPercent := fmt.Sprintf("%f", float32(curr.protien/total.protien)*100)
	fiberPercent := fmt.Sprintf("%f", (curr.fiber/total.fiber)*100)
	fmt.Printf("%s\t%s\t%s\t%s\t%s", carbsPercent, caloriesPercent, fatPercent, protienPercent, fiberPercent)
	pages.DailyInput(caloriesPercent, carbsPercent, protienPercent, fatPercent, fiberPercent).Render(r.Context(), w)
}

func (apiCfg *Api) DailyNutritionInputHandler(w http.ResponseWriter, r *http.Request) {
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

	var most_recent database.TotalCalorieIntake
	user_total, err := apiCfg.DB.GetMostRecentUserKcal(r.Context(), userId)
	if err == sql.ErrNoRows {
		user_create_total, err := apiCfg.DB.CreateTotalCalorieIntake(r.Context(), database.CreateTotalCalorieIntakeParams{
			ID:           uuid.New(),
			UserID:       userId,
			Calories:     0,
			TotalDeficit: 0,
			TotalSurplus: 0,
			CreatedAt:    time.Now().UTC(),
		})
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
		fmt.Println(most_recent)
	}
	fmt.Println(most_recent.CreatedAt.Unix())
	if time.Now().Unix() >= most_recent.CreatedAt.Unix()+(24*60*60) {
		fmt.Println("about to write to total nutrition database")
		var curr totalCalorieIntakeParams
		user_daily, err := apiCfg.DB.GetDailyNutritionOfUserByUserId(r.Context(), userId)
		if err != nil {
			log.Println(err)
		}
		user, err := apiCfg.DB.GetUserInputById(r.Context(), userId)
		if err != nil {
			log.Println(err)
		}
		curr.calories_you_should_have_eaten = float32(user.CurrKcal)
		for i := 0; i < len(user_daily); i++ {
			curr.calories_you_ate += float32(user_daily[i].Calories)
		}
		if curr.calories_you_should_have_eaten > curr.calories_you_ate {
			curr.deficit_for_the_day = curr.calories_you_should_have_eaten - curr.calories_you_ate
			curr.surplus_the_day = 0
		} else {
			curr.surplus_the_day = curr.calories_you_ate - curr.calories_you_should_have_eaten
			curr.deficit_for_the_day = 0
		}
		apiCfg.DB.CreateTotalCalorieIntake(r.Context(), database.CreateTotalCalorieIntakeParams{
			ID:           uuid.New(),
			Calories:     int32(curr.calories_you_ate),
			TotalDeficit: int32(curr.deficit_for_the_day),
			TotalSurplus: int32(curr.surplus_the_day),
			UserID:       userId,
			CreatedAt:    time.Now(),
		})
		apiCfg.DB.DeleteDailyNutritionOfUserByUserId(r.Context(), userId)
	}
	_, daily_create_err := apiCfg.DB.CreateDailyNutrition(r.Context(), database.CreateDailyNutritionParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		UserID:        userId,
		Calories:      int32(kCal),
		Carbohydrates: int32(carbs),
		Protien:       int32(protien),
		Fat:           int32(fat),
		Fiber:         int32(fiber),
	})
	if daily_create_err != nil {
		log.Println(daily_create_err)
	}
}
