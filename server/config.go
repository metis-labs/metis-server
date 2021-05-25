package server

import (
	"fmt"

	"oss.navercorp.com/metis/metis-server/server/database/mongodb"
	"oss.navercorp.com/metis/metis-server/server/rpc"
	"oss.navercorp.com/metis/metis-server/server/web"
)

// The following are the defaults for the Server config.
const (
	DefaultRPCPort = 10118

	DefaultWebPort = 10119

	DefaultMongoConnectionURI        = "mongodb://localhost:27017"
	DefaultMongoConnectionTimeoutSec = 5
	DefaultMongoPingTimeoutSec       = 5
	DefaultMongoDatabase             = "metis"
)

// Config is the configuration for creating a Server instance.
type Config struct {
	RPC   *rpc.Config     `json:"RPC"`
	Web   *web.Config     `json:"Web"`
	Mongo *mongodb.Config `json:"Mongo"`
}

// RPCAddr returns the RPC address.
func (c Config) RPCAddr() string {
	return fmt.Sprintf("localhost:%d", c.RPC.Port)
}

// NewConfig returns a Config struct that contains reasonable defaults
// for most of the configurations.
func NewConfig() *Config {
	return &Config{
		RPC: &rpc.Config{
			Port: DefaultRPCPort,
		},
		Web: &web.Config{
			Port: DefaultWebPort,
		},
		Mongo: &mongodb.Config{
			ConnectionURI:        DefaultMongoConnectionURI,
			ConnectionTimeoutSec: DefaultMongoConnectionTimeoutSec,
			PingTimeoutSec:       DefaultMongoPingTimeoutSec,
			Database:             DefaultMongoDatabase,
		},
	}
}
