package router

import (
	"log"
	"net/http"

	"github.com/vla47/go-api-mongo/store"

	"github.com/gorilla/mux"
	"github.com/vla47/go-api-mongo/middleware"
	"github.com/vla47/go-api-mongo/model"
)

var config *model.Config

type Env struct {
	db     store.Store
	Logger *log.Logger
	config *model.Config
}

// LoadRoutes loads all the routes and middlewares and creating the db
func LoadRoutes() error {
	router := mux.NewRouter()
	Init()

	db, _ := store.NewDB(*config)

	env := &Env{db: db}

	router.Use(env.Validate)
	router.Use(header.HandleCors)

	router.HandleFunc("/api/register", env.RegisterEndpoint).Methods("POST")
	router.HandleFunc("/api/login", env.LoginEndpoint).Methods("POST")
	router.HandleFunc("/api/account", env.AccountEndpoint).Methods("GET")

	router.HandleFunc("/api/courses", env.GetCoursesEndpoint).Methods("GET")
	router.HandleFunc("/api/courses", env.AddCourseHandler).Methods("POST")
	router.HandleFunc("/api/course/{id}", env.GetCourseHandler).Methods("GET")
	router.HandleFunc("/api/course", env.UpdateCourseHandler).Methods("PUT")
	router.HandleFunc("/api/course/{id}", env.DeleteCourseHandler).Methods("DELETE")
	router.HandleFunc("/api/courses/search/{term}", env.SearchCoursesEndpoint).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	// router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	return http.ListenAndServe(":"+config.Port, router)
}
