package mgo

import (
	"time"
)

type User struct {
	Name      string    `bson:"name"`
	Age       int       `bson:"age"`
	JonedAt   time.Time `bson:"joned_at"`
	Interests []string  `bson:"interests"`
}
