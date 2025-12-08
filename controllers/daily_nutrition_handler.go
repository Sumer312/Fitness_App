package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
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
	for _, ele := range user_daily {
		curr.carbs += float64(ele.Carbohydrates)
		curr.protein += float64(ele.Protein)
		curr.fat += float64(ele.Fat)
		curr.calories += float64(ele.Calories)
		curr.fiber += float64(ele.Fiber)
	}
	switch user.Sex {
	case sex_male:
		total.fiber = 31
	case sex_female:
		total.fiber = 21
	default:
		total.fiber = 26
	}

	program := user.Program
	switch program {
	case program_fatLoss:
		total.protein = float64(user.Weight)
		total.calories = float64(user.CurrKcal - user.Deficit.Float64)
		total.carbs = 0.45 * float64(total.calories/4)
		total.fat = 0.2 * float64(total.calories/9)
	case program_muscleGain:
		total.protein = 1.2 * float64(user.Weight)
		total.calories = float64(user.CurrKcal) + 200
		total.carbs = 0.4 * float64(total.calories/4)
		total.fat = 0.2 * float64(total.calories/9)
	default:
		total.protein = 0.8 * float64(user.Weight)
		total.calories = float64(user.CurrKcal)
		total.carbs = 0.6 * float64(total.calories/4)
		total.fat = 0.2 * float64(total.calories/9)
	}
	carbsPercent := (curr.carbs / total.carbs) * 100
	caloriesPercent := (curr.calories / total.calories) * 100
	fatPercent := float64(curr.fat/total.fat) * 100
	proteinPercent := float64(curr.protein/total.protein) * 100
	fiberPercent := (curr.fiber / total.fiber) * 100
	pages.DailyInput(caloriesPercent, carbsPercent, proteinPercent, fatPercent, fiberPercent).Render(r.Context(), w)
}

func (apiCfg *Api) DailyNutritionInputHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not parse data" }`)
		w.WriteHeader(500)
		return
	}
	userId, err := uuid.Parse(cookieVal.Value)
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not parse data" }`)
		w.WriteHeader(500)
		return
	}
	carbs, err := strconv.ParseFloat(r.FormValue("carbohydrates"), 10)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not parse data" }`)
		w.WriteHeader(500)
		return
	}
	protein, err := strconv.ParseFloat(r.FormValue("protein"), 10)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not parse data" }`)
		w.WriteHeader(500)
		return
	}
	fat, err := strconv.ParseFloat(r.FormValue("fat"), 10)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not parse data" }`)
		w.WriteHeader(500)
		return
	}
	fiber, err := strconv.ParseFloat(r.FormValue("fiber"), 10)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not parse data" }`)
		w.WriteHeader(500)
		return
	}
	kCal, err := strconv.ParseFloat(r.FormValue("calories"), 10)
	if err != nil && r.Form.Has("calories") {
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not parse data" }`)
		w.WriteHeader(500)
		return
	} else if !r.Form.Has("calories") {
		kCal = (carbs * 4) + (protein * 4) + (fat * 9)
	}

	user, err := apiCfg.DB.GetUserInputById(r.Context(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Add("HX-Trigger", `{ "errorToast" : "No program selected" }`)
			w.WriteHeader(500)
			return
		}
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "DB error" }`)
		w.WriteHeader(500)
		return
	} else if len(user.Program) == 0 {
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Choose a program" }`)
		w.WriteHeader(400)
		return

	}

	var most_recent database.TotalCalorieIntake
	user_total, err := apiCfg.DB.GetMostRecentUserKcalByUserId(r.Context(), userId)
	if err == sql.ErrNoRows {
		user_create_total, err := apiCfg.DB.CreateTotalCalorieIntake(r.Context(), database.CreateTotalCalorieIntakeParams{
			ID:           uuid.New(),
			CreatedAt:    time.Now().UTC(),
			UserID:       userId,
			Program:      user.Program,
			Calories:     0,
			TotalDeficit: 0,
			TotalSurplus: 0,
		})
		if err != nil {
			log.Println(err)
			w.Header().Add("HX-Trigger", `{ "errorToast" : "DB error" }`)
			w.WriteHeader(500)
			return
		} else {
			most_recent = user_create_total
		}
	} else if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "DB error" }`)
		w.WriteHeader(500)
		return
	} else {
		most_recent = user_total
	}
	if time.Now().Unix() >= most_recent.CreatedAt.Unix()+(24*60*60) {
		fmt.Println("about to write to total nutrition database")
		var curr totalCalorieIntakeParams
		user_daily, err := apiCfg.DB.GetDailyNutritionOfUserByUserId(r.Context(), userId)
		if err != nil {
			log.Println(err)
			w.Header().Add("HX-Trigger", `{ "errorToast" : "DB error" }`)
			w.WriteHeader(500)
			return
		}
		curr.calories_you_should_have_eaten = float64(user.CurrKcal)
		for _, ele := range user_daily {
			curr.calories_you_ate += float64(ele.Calories)
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
			CreatedAt:    time.Now(),
			UserID:       userId,
			Program:      user.Program,
			Calories:     (curr.calories_you_ate),
			TotalDeficit: (curr.deficit_for_the_day),
			TotalSurplus: (curr.surplus_the_day),
		})
		apiCfg.DB.DeleteDailyNutritionOfUserByUserId(r.Context(), userId)
	}
	_, daily_create_err := apiCfg.DB.CreateDailyNutrition(r.Context(), database.CreateDailyNutritionParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		UserID:        userId,
		Program:       user.Program,
		Calories:      kCal,
		Carbohydrates: carbs,
		Protein:       protein,
		Fat:           fat,
		Fiber:         fiber,
	})
	if daily_create_err != nil {
		log.Println(daily_create_err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not save to DB" }`)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(200)
}

func (apiCfg *Api) DailyNutritionDeleteRowById(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	base_url := os.Getenv("BASE_URL")
	r.ParseForm()
	var temp string = ""
	for k := range r.Form {
		if strings.Contains(k, "rowId") {
			temp = r.FormValue(k)
			break
		}
	}
	if len(temp) == 0 {
		w.Header().Add("HX-Redirect", base_url)
		w.WriteHeader(400)
	} else {
		rowId, err := uuid.Parse(temp)
		if err != nil {
			log.Println(err)
			return
		}
		err = apiCfg.DB.DeleteRowFromDailyNutritionTableById(r.Context(), rowId)
		if err != nil {
			log.Println(err)
			w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not delete row" }`)
			w.WriteHeader(500)
		}
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(200)
	}
}
