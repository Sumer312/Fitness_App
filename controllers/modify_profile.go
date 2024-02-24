package controllers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sumer312/Health-App-Backend/views/pages"
)

func (apiCfg *Api) ChangeProgram(w http.ResponseWriter, r *http.Request) {
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Redirect", "http://localhost:5000/view/login")
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
