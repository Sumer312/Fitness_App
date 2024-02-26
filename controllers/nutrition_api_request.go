package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func getEnv() api_parameters {
	godotenv.Load()
	accessPoint := os.Getenv("API_ACCESS_POINT")
	apiAppKey := os.Getenv("API_APP_KEY")
	apiAppId := os.Getenv("API_APP_ID")
	return api_parameters{access_point: accessPoint, app_key: apiAppKey, app_id: apiAppId}
}

func (apiCfg *Api) ApiRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := getEnv().access_point + "?app_id=" + getEnv().app_id + "&app_key=" + getEnv().app_key
	ingredients := r.FormValue("ingredients")
	ingredients = strings.Replace(ingredients, "\n", "", -1)
	ingredients = strings.Replace(ingredients, "\t", "", -1)
	ingredients_array := strings.Split(ingredients, ",")
	obj := map[string][]string{
		"ingr": ingredients_array,
	}
	json_obj, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not marshal json" }`)
		w.WriteHeader(500)
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewReader(json_obj))
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Server unable to send request" }`)
		w.WriteHeader(500)
		return
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Server error" }`)
		w.WriteHeader(500)
		return
	}
	if response.StatusCode == http.StatusOK {
		var response_variable edamam_response
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not read response" }`)
			w.WriteHeader(500)
			return
		}
		err = json.Unmarshal(body, &response_variable)
		if err != nil {
			log.Println(err)
			w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not unmarshal json" }`)
			w.WriteHeader(500)
			return
		}
		htmx_response, err := json.Marshal(response_variable)
		if err != nil {
			log.Println(err)
			w.Header().Add("HX-Trigger", `{ "errorToast" : "Could not marshal json" }`)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(htmx_response)
	} else {
		w.Header().Add("HX-Trigger", `{ "warnToast" : "Input not valid" }`)
		w.WriteHeader(400)
	}
}
