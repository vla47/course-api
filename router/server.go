package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/vla47/course-api/model"
	"github.com/vla47/course-api/store"
	"gopkg.in/mgo.v2/bson"
)

// LoadConfiguration load the cofig file
func LoadConfiguration(file string) *model.Config {
	// var config *model.Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

// Init initialize env vars and port
func Init() {
	if os.Getenv("ENV") == "dev" {
		config = LoadConfiguration("config.development.json")
	}
	if os.Getenv("ENV") == "prod" {
		config = LoadConfiguration("config.production.json")
		config.Database.Host = os.Getenv("MONGOLAB_URI")
		config.Port = os.Getenv("PORT")
	}

	if config.Port == "" {
		log.Fatal("$PORT must be set")
	}
}

// RegisterHandler register a new user
func (env *Env) RegisterHandler(w http.ResponseWriter, req *http.Request) {
	var creds *model.Credentials
	_ = json.NewDecoder(req.Body).Decode(&creds)

	err := store.Store.Register(env.db, creds)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(creds)
}

// LoginHandler is login a registered user
func (env *Env) LoginHandler(w http.ResponseWriter, req *http.Request) {
	var creds *model.Credentials
	var result map[string]interface{}
	_ = json.NewDecoder(req.Body).Decode(&creds)

	result, err := store.Store.Login(env.db, creds)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(result)
}

// AccountHandler shows the account of the user
func (env *Env) AccountHandler(w http.ResponseWriter, req *http.Request) {
	var profile *model.Profile

	profile, err := store.Store.GetAccount(env.db, context.Get(req, "pid").(string))
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(profile)
}

// GetCoursesHandler gets the courses
func (env *Env) GetCoursesHandler(w http.ResponseWriter, req *http.Request) {
	var courses []*model.Course

	courses, err := store.Store.GetCourses(env.db, context.Get(req, "pid").(string))
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	if courses == nil {
		courses = make([]*model.Course, 0)
	}
	encodeJSON(w, courses, env.Logger)
	// json.NewEncoder(w).Encode()
}

// SearchCoursesHandler is searching for a course
func (env *Env) SearchCoursesHandler(w http.ResponseWriter, req *http.Request) {
	var courses []*model.Course
	params := mux.Vars(req)
	courses, err := store.Store.SearchCourses(env.db, strings.ToLower(params["term"]))
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	if courses == nil {
		courses = make([]*model.Course, 0)
	}
	json.NewEncoder(w).Encode(courses)
}

func (env *Env) AddCourseHandler(w http.ResponseWriter, req *http.Request) {
	// rewrite all the encoding values to a message that it was successfull
	var course *model.Course
	_ = json.NewDecoder(req.Body).Decode(&course)
	fmt.Println(context.Get(req, "pid"))
	course.Type = "course"
	course.Pid = context.Get(req, "pid").(string)
	course.Timestamp = int(time.Now().Unix())
	course.ID = bson.NewObjectId()
	err := store.Store.AddCourse(env.db, course)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(course)
}

func (env *Env) UpdateCourseHandler(w http.ResponseWriter, req *http.Request) {
	// rewrite all the encoding values to a message that it was successfull
	var course *model.Course
	_ = json.NewDecoder(req.Body).Decode(&course)
	err := store.Store.UpdateCourse(env.db, course)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(course)
}

func (env *Env) GetCourseHandler(w http.ResponseWriter, req *http.Request) {
	var course *model.Course
	params := mux.Vars(req)
	id := params["id"]
	if !bson.IsObjectIdHex(id) {
		return
	}
	course, err := store.Store.GetCourse(env.db, bson.ObjectIdHex(id))
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(course)
}

func (env *Env) DeleteCourseHandler(w http.ResponseWriter, req *http.Request) {
	var course *model.Course
	params := mux.Vars(req)
	id := params["id"]
	err := store.Store.DeleteCourse(env.db, bson.ObjectIdHex(id))
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(course)
}

func (env *Env) Validate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				var session *model.Session
				session, err := store.Store.IsUserAuthenticated(env.db, bearerToken[1])
				if err != nil {
					w.WriteHeader(401)
					w.Write([]byte(err.Error()))
					return
				}
				context.Set(req, "pid", session.Pid)
				next.ServeHTTP(w, req)
			}
		} else {
			w.WriteHeader(401)
			w.Write([]byte("An authorization header is required"))
			return
		}
	})
}

// Error writes an API error message to the response and logger.
func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	// Log error.
	logger.Printf("http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		// err = wtf.ErrInternal
	}

	// Write generic error response.
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&model.ErrorResponse{Err: err.Error()})
}

// NotFoundHandler overrides the default not found handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// You can use the serve file helper to respond to 404 with
	// your request file.
	fmt.Println("sdadsd")
	http.ServeFile(w, r, "public/index.html")
}

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}
