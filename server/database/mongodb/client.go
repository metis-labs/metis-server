package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"oss.navercorp.com/metis/metis-server/internal/log"
	"oss.navercorp.com/metis/metis-server/server/database"
	"oss.navercorp.com/metis/metis-server/server/types"
)

// Config is the configuration for creating a Client instance.
type Config struct {
	ConnectionTimeoutSec time.Duration `json:"ConnectionTimeoutSec"`
	ConnectionURI        string        `json:"ConnectionURI"`
	Database             string        `json:"Database"`
	PingTimeoutSec       time.Duration `json:"PingTimeoutSec"`
}

// Client is a client that connects to Mongo DB and reads or saves Metis data.
type Client struct {
	config *Config
	client *mongo.Client
}

// NewClient creates a new instance of Client.
func NewClient(conf *Config) *Client {
	return &Client{
		config: conf,
	}
}

// Dial creates an instance of Client and dials the given MongoDB.
func (c *Client) Dial(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.config.ConnectionTimeoutSec*time.Second)
	defer cancel()

	log.Logger.Info("Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.config.ConnectionURI))
	if err != nil {
		return err
	}

	ctxPing, cancel := context.WithTimeout(ctx, c.config.PingTimeoutSec*time.Second)
	defer cancel()

	if err := client.Ping(ctxPing, readpref.Primary()); err != nil {
		log.Logger.Errorf("Could not connect to MongoDB: %s\n", err.Error())
		return err
	}
	log.Logger.Info("Connected to MongoDB")

	c.client = client
	return nil
}

// Close all resources of this client.
func (c *Client) Close(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

// CreateProject creates a new project of the given name.
func (c *Client) CreateProject(ctx context.Context, name string) (*types.ProjectInfo, error) {
	now := time.Now()
	result, err := c.client.Database(c.config.Database).Collection("projects").InsertOne(ctx, bson.M{
		"name":       name,
		"owner":      types.UserIDFromCtx(ctx),
		"status":     "created",
		"created_at": now,
	})
	if err != nil {
		return nil, err
	}

	return &types.ProjectInfo{
		ID:        types.ID(result.InsertedID.(primitive.ObjectID).Hex()),
		Name:      name,
		CreatedAt: now,
	}, nil
}

// FindProject returns the project of the given ID.
func (c *Client) FindProject(ctx context.Context, id types.ID) (*types.ProjectInfo, error) {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", id, database.ErrInvalidID)
	}

	result := c.client.Database(c.config.Database).Collection("projects").FindOne(ctx, bson.M{
		"_id":    objectID,
		"owner":  types.UserIDFromCtx(ctx),
		"status": "created",
	}, options.FindOne())

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s: %w", id, database.ErrNotFound)
		}
		return nil, result.Err()
	}

	project := &types.ProjectInfo{}
	idHolder := struct {
		ID primitive.ObjectID `bson:"_id"`
	}{}
	if err := result.Decode(&idHolder); err != nil {
		return nil, err
	}
	if err := result.Decode(project); err != nil {
		return nil, err
	}
	project.ID = types.ID(idHolder.ID.Hex())
	return project, nil
}

// ListProjects returns the list of projects.
func (c *Client) ListProjects(ctx context.Context) ([]*types.ProjectInfo, error) {
	cursor, err := c.client.Database(c.config.Database).Collection("projects").Find(ctx, bson.M{
		"owner":  types.UserIDFromCtx(ctx),
		"status": "created",
	}, options.Find())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Logger.Error(err)
		}
	}()

	var projects []*types.ProjectInfo
	for cursor.Next(ctx) {
		var project types.ProjectInfo
		idHolder := struct {
			ID primitive.ObjectID `bson:"_id"`
		}{}
		if err := cursor.Decode(&idHolder); err != nil {
			return nil, err
		}
		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}
		project.ID = types.ID(idHolder.ID.Hex())
		projects = append(projects, &project)
	}

	return projects, nil
}

// UpdateProject updates the given project.
func (c *Client) UpdateProject(ctx context.Context, id types.ID, name string) error {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return fmt.Errorf("%s: %w", id, database.ErrInvalidID)
	}

	result := c.client.Database(c.config.Database).Collection("projects").FindOneAndUpdate(
		ctx,
		bson.M{
			"_id":    objectID,
			"owner":  types.UserIDFromCtx(ctx),
			"status": "created",
		},
		bson.M{
			"$set": bson.M{
				"name": name,
			},
		},
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return fmt.Errorf("%s: %w", id, database.ErrNotFound)
		}
		return result.Err()
	}

	return nil
}

// DeleteProject deletes the given project.
func (c *Client) DeleteProject(ctx context.Context, id types.ID) error {
	objectID, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return err
	}

	_, err = c.client.Database(c.config.Database).Collection("projects").UpdateOne(ctx, bson.M{
		"_id":   objectID,
		"owner": types.UserIDFromCtx(ctx),
	}, bson.M{
		"$set": bson.M{
			"status":     "deleted",
			"deleted_at": time.Now(),
		},
	})

	return err
}
