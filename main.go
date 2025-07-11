package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sumer312/Health-App-Backend/controllers"
	"github.com/sumer312/Health-App-Backend/internal/database"
	"github.com/sumer312/Health-App-Backend/views/pages"
	"github.com/sumer312/Health-App-Backend/views/partials"
)

func main() {
	godotenv.Load()
	port := "5000"
	router := chi.NewRouter()
	viewRouter := chi.NewRouter()
	serverRouter := chi.NewRouter()
	dbConnString := os.Getenv("DB_URL")
	conn, connerr := sql.Open("postgres",
		dbConnString,
	)
	if connerr != nil {
		log.Fatalln("error connecting", connerr)
	}
	db := database.New(conn)
	apiCfg := controllers.Api{DB: db}
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))
	viewRouter.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if partials.DrawerAuthFlag {
			pages.Random("Already Logged In").Render(r.Context(), w)
		} else {
			pages.Login().Render(r.Context(), w)
		}
	})
	viewRouter.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if partials.DrawerAuthFlag {
			pages.Random("Already Logged In").Render(r.Context(), w)
		} else {
			pages.Signup().Render(r.Context(), w)
		}
	})
	viewRouter.HandleFunc("/user-input/fatloss", func(w http.ResponseWriter, r *http.Request) {
		if partials.DrawerAuthFlag {
			pages.UserInputFatloss().Render(r.Context(), w)
		} else {
			pages.Login().Render(r.Context(), w)
		}
	})
	viewRouter.HandleFunc("/user-input/muscle", func(w http.ResponseWriter, r *http.Request) {
		if partials.DrawerAuthFlag {
			pages.UserInputMuscle().Render(r.Context(), w)
		} else {
			pages.Login().Render(r.Context(), w)
		}
	})
	viewRouter.HandleFunc("/user-input/maintain", func(w http.ResponseWriter, r *http.Request) {
		if partials.DrawerAuthFlag {
			pages.UserInputMaintain().Render(r.Context(), w)
		} else {
			pages.Login().Render(r.Context(), w)
		}
	})
	viewRouter.HandleFunc(
		"/logs",
		viewRenderInControllerMiddleware(func(w http.ResponseWriter, r *http.Request) {
			apiCfg.LogsRender(w, r)
		}),
	)
	viewRouter.HandleFunc(
		"/daily-input",
		viewRenderInControllerMiddleware(func(w http.ResponseWriter, r *http.Request) {
			apiCfg.DailyNutritionRender(w, r)
		}),
	)
	viewRouter.HandleFunc(
		"/profile",
		viewRenderInControllerMiddleware(func(w http.ResponseWriter, r *http.Request) {
			apiCfg.ProfileRender(w, r)
		}),
	)
	viewRouter.HandleFunc(
		"/kcal-calc",
		viewRenderInControllerMiddleware(func(w http.ResponseWriter, r *http.Request) {
			pages.KcalCalc().Render(r.Context(), w)
		}),
	)

	serverRouter.Post("/login", apiCfg.LoginHandler)
	serverRouter.Post("/signup", apiCfg.SignupHandler)
	serverRouter.Post("/logout", apiCfg.LogoutHandler)
	serverRouter.Post("/user-input", controllerMiddleware(apiCfg.InputHandler))
	serverRouter.Post("/nutrition-api-request", controllerMiddleware(apiCfg.ApiRequest))
	serverRouter.Post("/daily-input", controllerMiddleware(apiCfg.DailyNutritionInputHandler))
	serverRouter.Post(
		"/daily-input-delete",
		controllerMiddleware(apiCfg.DailyNutritionDeleteRowById),
	)
	serverRouter.Delete("/change-program", controllerMiddleware(apiCfg.ChangeProgram))
	serverRouter.Delete("/delete-user", controllerMiddleware(apiCfg.DeleteUser))

	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if partials.DrawerAuthFlag {
			apiCfg.HomeRender(w, r)
		} else {
			pages.Programs().Render(r.Context(), w)
		}
	})
	viewRouter.Handle("/*", templ.Handler(pages.Random("404 not found")))
	router.Mount("/view", viewRouter)
	router.Mount("/server", serverRouter)

	router.Handle("/*", templ.Handler(pages.Random("404 not found")))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	fmt.Printf("Using chi \nServer starting on port %s\n", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln("OOPs something went wrong")
	}
}
