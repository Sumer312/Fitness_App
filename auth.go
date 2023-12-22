package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sumer312/Health-App-Backend/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func createJWT(expiresIn time.Duration, subject string) (string, error) {
	godotenv.Load()
	var SECRET = []byte(os.Getenv("JWT_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claim := token.Claims.(jwt.MapClaims)
	claim["exp"] = time.Now().Add(expiresIn).Unix()
	claim["sub"] = subject
	tokenStr, err := token.SignedString(SECRET)
	if err != nil {
		log.Fatal("line 26", err)
		return "", err
	}
	return tokenStr, nil
}

func (apiCfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := apiCfg.DB.GetUserByEmail(
		r.Context(),
		sql.NullString{String: email, Valid: true},
	)
	if err != nil {
		log.Fatalln("line 54", err)
	}
	fmt.Println(user)
	passwordCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if passwordCheck != nil {
		log.Printf("wrong password")
		return
	}
	accessToken, err := createJWT(time.Minute*5, user.Name)
	if err != nil {
		log.Fatalln("line 64", err)
	}
	refreshToken, err := createJWT(time.Hour*840, user.Name)
	if err != nil {
		log.Fatalln("line 68", err)
	}
	access_cookie := http.Cookie{Name: "access-token", Value: accessToken, HttpOnly: true, SameSite: http.SameSiteNoneMode}
	refresh_cookie := http.Cookie{Name: "refresh-token", Value: refreshToken, HttpOnly: true, SameSite: http.SameSiteNoneMode}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	http.SetCookie(w, &access_cookie)
	http.SetCookie(w, &refresh_cookie)
	values := map[string]string{"userID": user.ID.String()}
	json_values, err := json.Marshal(values)
	if err != nil {
		fmt.Println("line 80", err)
	}
	w.Write(json_values)
}

func (apiCfg *apiConfig) signupHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatalln("line 95", err)
		return
	}

	if params.Email == "" {
		log.Fatalln("enter a valid email")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("error storing password")
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Email:     sql.NullString{String: params.Email, Valid: true},
		Password:  string(hashedPassword),
	})
	if err != nil {
		log.Fatalln("line 117", err)
	}
	accessToken, err := createJWT(time.Minute*5, user.Name)
	if err != nil {
		log.Fatalln("line 64", err)
	}
	refreshToken, err := createJWT(time.Hour*840, user.Name)
	if err != nil {
		log.Fatalln("line 68", err)
		access_cookie := http.Cookie{Name: "access-token", Value: accessToken, HttpOnly: true, Secure: true, SameSite: http.SameSiteNoneMode}
		refresh_cookie := http.Cookie{Name: "refresh-token", Value: refreshToken, HttpOnly: true, Secure: true, SameSite: http.SameSiteNoneMode}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		http.SetCookie(w, &access_cookie)
		http.SetCookie(w, &refresh_cookie)
		values := map[string]string{"userID": user.ID.String()}
		json_values, err := json.Marshal(values)
		if err != nil {
			fmt.Println("line 80", err)
		}
		w.Write(json_values)
	}
}
