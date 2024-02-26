package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *Api) HomeRender(w http.ResponseWriter, r *http.Request) {
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Not logged in" }`)
		w.WriteHeader(400)
		return
	}
	userId, err := uuid.Parse(cookieVal.Value)
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Not logged in" }`)
		w.WriteHeader(400)
		return
	}
	_, err = apiCfg.DB.GetUserInputById(r.Context(), userId)
	if err == sql.ErrNoRows {
		apiCfg.ProfileRender(w, r)
		return
	} else if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Not logged in" }`)
		w.WriteHeader(400)
		return
	}
	apiCfg.LogsRender(w, r)
}
