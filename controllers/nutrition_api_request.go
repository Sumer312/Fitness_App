package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type api_parameters struct {
	access_point string
	app_key      string
	app_id       string
}

type edamam_response_total_nutrients_element struct {
	Label    string  `json:"label"`
	Unit     string  `json:"unit"`
	Quantity float32 `json:"quantity"`
}

type total_nutrients struct {
	Enengc_Kcal edamam_response_total_nutrients_element `json:"ENENGC_KCAL"`
	Fat         edamam_response_total_nutrients_element `json:"FAT"`
	Fasat       edamam_response_total_nutrients_element `json:"FASAT"`
	Fatrn       edamam_response_total_nutrients_element `json:"FATRN"`
	Fibtg       edamam_response_total_nutrients_element `json:"FIBTG"`
	Chocdf      edamam_response_total_nutrients_element `json:"CHOCDF"`
	Sugar       edamam_response_total_nutrients_element `json:"SUGAR"`
	Procnt      edamam_response_total_nutrients_element `json:"PROCNT"`
}

type edamam_response struct {
	Calories       int             `json:"calories"`
	TotalNutrients total_nutrients `json:"totalNutrients"`
	TotalWeight    float32         `json:"totalWeight"`
}

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
	fmt.Printf("%s\n", json_obj)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewReader(json_obj))
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(response.StatusCode)
	if response.StatusCode == http.StatusOK {
		var response_variable edamam_response
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		err = json.Unmarshal(body, &response_variable)
		if err != nil {
			log.Fatalln(err)
		}
		htmx_response, err := json.Marshal(response_variable)
		if err != nil {
			log.Fatalln(err)
		}
		w.Header().Add("HX-Trigger", `{ "infoToast" : "Click to copy values" }`)
		w.Write(htmx_response)
		w.WriteHeader(200)
	} else {
		w.Header().Add("HX-Trigger", `{ "warnToast" : "Input not valid" }`)
		w.WriteHeader(400)
	}
}
