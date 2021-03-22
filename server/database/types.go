package database

import "time"

type Model struct {
	ID   ID     `bson:"_id_fake"`
	Name string `bson:"name"`
}

type Project struct {
	ID        ID        `bson:"_id_fake"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
}
