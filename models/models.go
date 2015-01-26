package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Header struct {
	Field string `bson:"field" json:"field"`
	Value string `bson:"value" json:"value"`
}

type Task struct {
	Method  string   `bson:"method" json:"method"`
	Path    string   `bson:"path" json:"path"`
	Headers []Header `bson:"headers" json:"headers"`
	RawBody string   `bson:"raw_body" json:"rawBody"`
}

type Test struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	BaseUrl     string        `bson:"base_url" json:"baseUrl"`
	Requests    int           `bson:"requests" json:"requests"`
	Concurrency int           `bson:"concurrency" json:"concurrency"`
	Tasks       []Task        `bson:"tasks" json:"tasks"`
}
