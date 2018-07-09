package store

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/vla47/course-api/model"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DB struct {
	db *mgo.Database
}

func NewDB(config model.Config) (*DB, error) {
	session, err := mgo.Dial(config.Database.Host)
	if err != nil {
		fmt.Printf("Something went wrong: db %s", err)
		return &DB{session.DB(config.Database.Name)}, err
	}

	return &DB{session.DB(config.Database.Name)}, nil
}

type Store interface {
	Register(*model.Credentials) error
	Login(*model.Credentials) (map[string]interface{}, error)
	GetAccount(string) (*model.Profile, error)
	GetCourses(string) ([]*model.Course, error)
	SearchCourses(string) ([]*model.Course, error)
	AddCourse(*model.Course) error
	UpdateCourse(*model.Course) error
	GetCourse(bson.ObjectId) (*model.Course, error)
	DeleteCourse(bson.ObjectId) error
	IsUserAuthenticated(string) (*model.Session, error)
}

func (db *DB) Register(cred *model.Credentials) error {
	ID, _ := uuid.NewV4()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(cred.Password), 14)
	account := model.Account{
		Type:     "account",
		Pid:      ID.String(),
		Email:    cred.Email,
		Password: string(passwordHash),
	}
	profile := model.Profile{
		Type:      "profile",
		Firstname: cred.Firstname,
		Lastname:  cred.Lastname,
	}
	fmt.Println(account.Email, profile.Firstname)
	_, err := db.db.C("user").UpsertId(cred.Email, account)
	_, err = db.db.C("user").UpsertId(ID.String(), profile)

	return err
}

func (db *DB) Login(cred *model.Credentials) (map[string]interface{}, error) {
	var account *model.Account
	err := db.db.C("user").FindId(cred.Email).One(&account)
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(cred.Password))

	session := model.Session{
		Type: "session",
		Pid:  account.Pid,
	}
	var result map[string]interface{}
	result = make(map[string]interface{})
	ID, err := uuid.NewV4()
	result["sid"] = ID.String()

	_, err = db.db.C("user").UpsertId(ID.String(), session)

	return result, err
}

func (db *DB) GetAccount(pid string) (*model.Profile, error) {

	var profile *model.Profile
	err := db.db.C("user").FindId(bson.M{"pid": pid}).One(&profile)

	return profile, err
}

func (db *DB) GetCourses(pid string) ([]*model.Course, error) {
	var courses []*model.Course
	err := db.db.C("courses").Find(bson.M{"pid": pid}).All(&courses)

	return courses, err
}

func (db *DB) SearchCourses(term string) ([]*model.Course, error) {
	var courses []*model.Course
	err := db.db.C("courses").Find(bson.M{"name": bson.M{"$regex": term}}).All(&courses)

	return courses, err
}

func (db *DB) AddCourse(course *model.Course) error {
	err := db.db.C("courses").Insert(course)
	return err
}

func (db *DB) UpdateCourse(course *model.Course) error {
	course.Timestamp = int(time.Now().Unix())
	_, err := db.db.C("courses").UpsertId(course.ID, course)
	return err
}

func (db *DB) GetCourse(id bson.ObjectId) (*model.Course, error) {
	var course *model.Course
	err := db.db.C("courses").FindId(id).One(&course)
	return course, err
}

func (db *DB) DeleteCourse(id bson.ObjectId) error {
	err := db.db.C("courses").RemoveId(id)
	return err
}

func (db *DB) IsUserAuthenticated(token string) (*model.Session, error) {
	var session *model.Session
	err := db.db.C("user").FindId(token).One(&session)

	return session, err
}
