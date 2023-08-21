package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sumer312/Health-App-Backend/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := "5000"
	router := chi.NewRouter()
	fmt.Println("using chi")
	authRouter := chi.NewRouter()
	dbConnString := os.Getenv("DB_URL")
	log.Println(dbConnString)
	conn, connerr := sql.Open("postgres",
		dbConnString,
	)
	if connerr != nil {
		log.Fatalln("error conncting", connerr)
	}
	db := database.New(conn)
	apiCfg := apiConfig{DB: db}
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	authRouter.HandleFunc("/login", apiCfg.loginHandler)
	router.HandleFunc("/user-input", validateJWT(apiCfg.input_handler))
	router.HandleFunc("/profile", validateJWT(apiCfg.profile))
	authRouter.HandleFunc("/signup", apiCfg.signupHandler)
	authRouter.HandleFunc("/calorie-tracker", validateJWT(apiCfg.calorie_input_handler))

	router.Mount("/auth", authRouter)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Println("Server starting on port:" + port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("OOPs something went wrong")
	}
}
