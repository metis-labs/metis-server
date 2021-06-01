package types

import "time"

// TemplateInfo represents the template.
type TemplateInfo struct {
	ID        ID        `bson:"_id_fake"`
	Name      string    `bson:"name"`
	Owner     string    `bson:"owner"`
	Contents  string    `bson:"contents"`
	CreatedAt time.Time `bson:"created_at"`
}
