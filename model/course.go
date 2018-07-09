package model

import "gopkg.in/mgo.v2/bson"

type Course struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Type      string        `json:"type,omitempty"`
	Pid       string        `json:"pid,omitempty"`
	Name      string        `json:"name,omitempty"`
	Code      string        `json:"code,omitempty"`
	Start     int           `json:"start,omitempty"`
	End       int           `json:"end,omitempty"`
	Grade     float64       `json:"grade,omitempty"`
	Timestamp int           `json:"timestamp,omitempty"`
}
