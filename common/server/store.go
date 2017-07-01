package server

import (
	"gopkg.in/mgo.v2"
)

// Connection to the mongo database
var MongoConn *mgo.Session

var Db *mgo.Database

func ConnectToDatabase() (err error) {
	MongoConn, err = mgo.Dial("localhost")
	if err != nil {
		return
	}
	Db = MongoConn.DB("neo")

	return
}
