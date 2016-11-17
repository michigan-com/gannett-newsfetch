package main

import (
	"gopkg.in/mgo.v2"
)

func SetupMongoSession(uri string) (*mgo.Session, error) {
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}

	return session, nil
}
