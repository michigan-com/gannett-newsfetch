package model

import (
	"gopkg.in/mgo.v2"
)

type Model interface {
	Save(session *mgo.Session)
}
