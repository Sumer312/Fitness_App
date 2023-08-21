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
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatal("line 41", err)
	}
	if params.Email == "" {
		log.Fatalln("enter a valid email")
		return
	}
	user, err := apiCfg.DB.GetUserByEmail(
		r.Context(),
		sql.NullString{String: params.Email, Valid: true},
	)
	if err != nil {
		log.Fatalln("line 52", err)
	}
	fmt.Println(user)
	passwordCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if passwordCheck != nil {
		log.Fatalln("error storing passoword")
		return
	}
	accesstoken, err := createJWT(time.Minute*5, user.Name)
	if err != nil {
		log.Fatalln("line 62", err)
	}
	refreshToken, err := createJWT(time.Hour*840, user.Name)
	if err != nil {
		log.Fatalln("line 66", err)
	}
	cookie := http.Cookie{Name: "access-token", Value: accesstoken, HttpOnly: true}
	w.Header().Add("refresh-token", refreshToken)
	w.Write([]byte(user.ID.String()))
	http.SetCookie(w, &cookie)
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
		log.Fatalln("line 84", err)
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
		log.Fatalln("line 106", err)
	}
	accesstoken, err := createJWT(time.Minute*5, user.ID.String())
	if err != nil {
		log.Fatalln("line 110", err)
	}
	refreshToken, err := createJWT(time.Hour*840, user.ID.String())
	if err != nil {
		log.Fatalln("line 114", err)
	}
	fmt.Println(user)
	cookie := http.Cookie{Name: "access-token", Value: accesstoken, HttpOnly: true}
	w.Header().Add("refresh-token", refreshToken)
	w.Write([]byte(user.ID.String()))
	http.SetCookie(w, &cookie)
}
