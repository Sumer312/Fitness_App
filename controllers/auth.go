package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sumer312/Health-App-Backend/internal/database"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type Api struct {
	DB *database.Queries
}

func CreateJWT(expiresIn time.Duration, subject string) (string, error) {
	godotenv.Load()
	var SECRET = []byte(os.Getenv("JWT_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claim := token.Claims.(jwt.MapClaims)
	claim["exp"] = time.Now().Add(expiresIn).Unix()
	claim["sub"] = subject
	tokenStr, err := token.SignedString(SECRET)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return tokenStr, nil
}

func (apiCfg *Api) LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for k, vs := range r.Form {
		for _, v := range vs {
			fmt.Printf("%s => %s\n", k, v)
		}
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := apiCfg.DB.GetUserByEmail(
		r.Context(),
		sql.NullString{String: email, Valid: true},
	)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(user)
	passwordCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if passwordCheck != nil {
		log.Printf("wrong password")
		return
	}
	accessToken, err := CreateJWT(time.Minute*5, user.Name)
	if err != nil {
		log.Fatalln(err)
	}
	refreshToken, err := CreateJWT(time.Hour*840, user.Name)
	if err != nil {
		log.Fatalln(err)
	}
	access_cookie := http.Cookie{Name: "access-token", Path: "/", Value: accessToken, HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
	refresh_cookie := http.Cookie{Name: "refresh-token", Path: "/", Value: refreshToken, HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
	user_id := http.Cookie{Name: "user-id", Path: "/", Value: user.ID.String(), HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	http.SetCookie(w, &access_cookie)
	http.SetCookie(w, &refresh_cookie)
	http.SetCookie(w, &user_id)
}

func (apiCfg *Api) SignupHandler(w http.ResponseWriter, r *http.Request) {
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
	accessToken, err := CreateJWT(time.Minute*5, user.Name)
	if err != nil {
		log.Fatalln("access token signup error ", err)
	}
	refreshToken, err := CreateJWT(time.Hour*840, user.Name)
	if err != nil {
		log.Fatalln("refresh token signup err ", err)
		access_cookie := http.Cookie{Name: "access-token", Path: "/", Value: accessToken, HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
		refresh_cookie := http.Cookie{Name: "refresh-token", Path: "/", Value: refreshToken, HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
		user_id := http.Cookie{Name: "user_id", Path: "/", Value: user.ID.String(), HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		http.SetCookie(w, &access_cookie)
		http.SetCookie(w, &refresh_cookie)
		http.SetCookie(w, &user_id)
	}
}
