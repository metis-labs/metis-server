package database

import "time"

// Project represents the metadata of the project of Metis.
type Project struct {
	ID        ID        `bson:"_id_fake"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
}
