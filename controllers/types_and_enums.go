package controllers

import "github.com/sumer312/Health-App-Backend/internal/database"

type Api struct {
	DB *database.Queries
}

type nutritionParams struct {
	calories float64
	carbs    float64
	protien  float64
	fat      float64
	fiber    float64
}

type totalCalorieIntakeParams struct {
	calories_you_ate               float64
	calories_you_should_have_eaten float64
	deficit_for_the_day            float64
	surplus_the_day                float64
}

type api_parameters struct {
	access_point string
	app_key      string
	app_id       string
}

type edamam_response_total_nutrients_element struct {
	Label    string  `json:"label"`
	Unit     string  `json:"unit"`
	Quantity float64 `json:"quantity"`
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
	TotalWeight    float64         `json:"totalWeight"`
}

const (
	program_fatLoss    = "fatloss"
	program_muscleGain = "musclegain"
	program_maintain   = "maintaince"
	sex_male           = "male"
	sex_female         = "female"
	sex_none           = "none"
)

const (
  access_token_cookie_name = "access-token"
  refresh_token_cookie_name = "refresh-token"
  user_id_cookie_name = "user-id"
)
