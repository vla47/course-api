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

	db, _ := store.NewDB(*config)

	env := &Env{db: db}

	router.Use(env.Validate)
	router.Use(HandleCors)

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
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})
	// handler := cors.Default().Handler(router)
	// return http.ListenAndServe(":"+config.Port, handler)

	return http.ListenAndServe(":"+config.Port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))

}

// HandleCors is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func HandleCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
