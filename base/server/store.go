package server

import (
	"gopkg.in/mgo.v2"
)

type InitializeCollectionsFunc func(*mgo.Database)

var initializeCollectionsFuncs []InitializeCollectionsFunc

func RegisterInitializeCollections(f InitializeCollectionsFunc) {
	initializeCollectionsFuncs = append(initializeCollectionsFuncs, f)
}

// Connection to the mongo database
var MongoConn *mgo.Session

var Db *mgo.Database

func ConnectToDatabase() (err error) {
	MongoConn, err = mgo.Dial("localhost")
	if err != nil {
		return
	}
	Db = MongoConn.DB("neo")

	for _, initializeCollections := range initializeCollectionsFuncs {
		initializeCollections(Db)
	}

	return
}
