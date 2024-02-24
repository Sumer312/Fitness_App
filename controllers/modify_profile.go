package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sumer312/Health-App-Backend/views/pages"
)

func (apiCfg *Api) ChangeProgram(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	base_url := os.Getenv("BASE_URL")
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Redirect", fmt.Sprintf("%s/view/login", base_url))
		w.WriteHeader(500)
		return
	}
	userId, err := uuid.Parse(cookieVal.Value)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not parse data" }`)
		w.WriteHeader(500)
		return
	}
	err = apiCfg.DB.DeleteDailyNutritionOfUserByUserId(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	err = apiCfg.DB.DeleteUserRecord(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	err = apiCfg.DB.DeleteUserInputByUserId(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	pages.Programs().Render(r.Context(), w)
	return
}

func (apiCfg *Api) DeleteUser(w http.ResponseWriter, r *http.Request) {
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		fmt.Println(err)
	}
	userId, err := uuid.Parse(cookieVal.Value)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not parse data" }`)
		w.WriteHeader(500)
		return
	}
	err = apiCfg.DB.DeleteDailyNutritionOfUserByUserId(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	err = apiCfg.DB.DeleteUserRecord(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	err = apiCfg.DB.DeleteUserInputByUserId(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	err = apiCfg.DB.DeleteUserById(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	apiCfg.LogoutHandler(w, r)
}
