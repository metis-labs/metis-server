package server

import (
	"oss.navercorp.com/metis/metis-server/server/database/mongodb"
	"oss.navercorp.com/metis/metis-server/server/rpc"
)

// The following are the defaults for the Server config.
const (
	DefaultRPCPort = 10118

	DefaultMongoConnectionURI        = "mongodb://localhost:27017"
	DefaultMongoConnectionTimeoutSec = 5
	DefaultMongoPingTimeoutSec       = 5
	DefaultMongoDatabase             = "metis"
)

// Config is the configuration for creating a Server instance.
type Config struct {
	RPC   *rpc.Config     `json:"RPC"`
	Mongo *mongodb.Config `json:"Mongo"`
}

// NewConfig returns a Config struct that contains reasonable defaults
// for most of the configurations.
func NewConfig() *Config {
	return &Config{
		RPC: &rpc.Config{
			Port: DefaultRPCPort,
		},
		Mongo: &mongodb.Config{
			ConnectionURI:        DefaultMongoConnectionURI,
			ConnectionTimeoutSec: DefaultMongoConnectionTimeoutSec,
			PingTimeoutSec:       DefaultMongoPingTimeoutSec,
			Database:             DefaultMongoDatabase,
		},
	}
}
