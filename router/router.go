package router

import (
	"log"
	"net/http"

	"github.com/vla47/course-api/store"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/vla47/course-api/model"
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
	var headers = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	var methods = handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"})
	var origins = handlers.AllowedOrigins([]string{"*"})

	db, _ := store.NewDB(*config)

	env := &Env{db: db}

	router.HandleFunc("/api/register", env.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/login", env.LoginHandler).Methods("POST")
	router.HandleFunc("/api/account", env.Validate(env.AccountHandler)).Methods("GET")

	router.HandleFunc("/api/courses", env.Validate(env.GetCoursesHandler)).Methods("GET")
	router.HandleFunc("/api/courses", env.Validate(env.AddCourseHandler)).Methods("POST")
	router.HandleFunc("/api/course/{id}", env.Validate(env.GetCourseHandler)).Methods("GET")
	router.HandleFunc("/api/course", env.Validate(env.UpdateCourseHandler)).Methods("PUT")
	router.HandleFunc("/api/course/{id}", env.Validate(env.DeleteCourseHandler)).Methods("DELETE")
	router.HandleFunc("/api/courses/search/{term}", env.Validate(env.SearchCoursesHandler)).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	return http.ListenAndServe(":"+config.Port, handlers.CORS(headers, methods, origins)(router))

}
