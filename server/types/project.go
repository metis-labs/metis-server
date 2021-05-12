package types

import "time"

// Project represents the metadata of the project of Metis.
type Project struct {
	ID        ID        `bson:"_id_fake"`
	Name      string    `bson:"name"`
	Owner     string    `bson:"owner"`
	CreatedAt time.Time `bson:"created_at"`
}
