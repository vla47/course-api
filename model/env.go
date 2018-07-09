package model

import mgo "gopkg.in/mgo.v2"

type Env struct {
	db *mgo.Database
	// config *model.Config
}
