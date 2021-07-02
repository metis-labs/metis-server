/*
 * Copyright 2021-present NAVER Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"fmt"

	"github.com/metis-labs/metis-server/server/database/mongodb"
	"github.com/metis-labs/metis-server/server/rpc"
	"github.com/metis-labs/metis-server/server/web"
	"github.com/metis-labs/metis-server/server/yorkie"
)

// The following are the defaults for the Server config.
const (
	DefaultRPCPort = 10118

	DefaultWebPort = 10119

	DefaultMongoConnectionURI        = "mongodb://localhost:27017"
	DefaultMongoConnectionTimeoutSec = 5
	DefaultMongoPingTimeoutSec       = 5
	DefaultMongoDatabase             = "metis"

	DefaultYorkieRPCAddr      = "localhost:11101"
	DefaultYorkieWebhookToken = "metis-server"
	DefaultYorkieCollection   = "projects"
)

// Config is the configuration for creating a Server instance.
type Config struct {
	RPC    *rpc.Config     `json:"RPC"`
	Web    *web.Config     `json:"Web"`
	Mongo  *mongodb.Config `json:"Mongo"`
	Yorkie *yorkie.Config  `json:"Yorkie"`
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
		Yorkie: &yorkie.Config{
			RPCAddr:      DefaultYorkieRPCAddr,
			WebhookToken: DefaultYorkieWebhookToken,
			Collection:   DefaultYorkieCollection,
		},
	}
}
