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

package test

import (
	"log"
	"os"
	"testing"

	"github.com/metis-labs/metis-server/server"
	"github.com/metis-labs/metis-server/server/database/mongodb"
	"github.com/metis-labs/metis-server/server/rpc"
	"github.com/metis-labs/metis-server/server/web"
	"github.com/metis-labs/metis-server/server/yorkie"
)

var testServer *server.Server

const (
	testUserA = "KR18401"
	testUserB = "KR18817"
)

func TestMain(m *testing.M) {
	s, err := server.New(&server.Config{
		RPC: &rpc.Config{
			Port: server.DefaultRPCPort,
		},
		Web: &web.Config{
			Port: server.DefaultWebPort,
		},
		Yorkie: &yorkie.Config{
			RPCAddr:      server.DefaultYorkieRPCAddr,
			WebhookToken: server.DefaultYorkieWebhookToken,
			Collection:   server.DefaultYorkieCollection,
		},
		Mongo: &mongodb.Config{
			ConnectionURI:        server.DefaultMongoConnectionURI,
			ConnectionTimeoutSec: server.DefaultMongoConnectionTimeoutSec,
			PingTimeoutSec:       server.DefaultMongoPingTimeoutSec,
			Database:             server.DefaultMongoDatabase,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	testServer = s

	if err := testServer.Start(); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := testServer.Shutdown(true); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
