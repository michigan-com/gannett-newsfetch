package brvtycore

import (
	"gopkg.in/mgo.v2"

	"github.com/michigan-com/brvty-api/mongoqueue"
)

func New(db *mgo.Database) (*Core, error) {
	db.Session.EnsureSafe(&mgo.Safe{WMode: "majority", FSync: true})
	coll := db.C("resources")

	// err = collection.EnsureIndex(Index{
	// 	Key: []string{"ns"},
	// })
	// if err != nil {
	// 	return err
	// }

	queue := mongoqueue.New(db.C("queue"), mongoqueue.Params{})

	err := queue.Migrate()
	if err != nil {
		return nil, err
	}

	return &Core{
		DB:         db,
		Queue:      queue,
		collection: coll,
	}, nil
}
