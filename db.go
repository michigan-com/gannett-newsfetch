package main

import (
	"gopkg.in/mgo.v2"
)

func SetupMongoSession(uri string) (*mgo.Session, error) {
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

func closeSessionIfNotNil(session *mgo.Session) {
	if session != nil {
		session.Close()
	}
}
