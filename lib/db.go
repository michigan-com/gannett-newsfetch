package lib

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func DBConnect(uri string) *mgo.Session {
	// TODO read this from config
	session, err := mgo.Dial(uri)
	if err != nil {
		log.Printf("Failed to connect to '%s': %v", uri, err)
		os.Exit(1)
	}

	session.SetMode(mgo.Monotonic, true)
	return session
}

func DBClose(session *mgo.Session) {
	session.Close()
}
