package client

import (
	"log"
	"os"
	"testing"

	"oss.navercorp.com/metis/metis-server/server"
	"oss.navercorp.com/metis/metis-server/server/database/mongodb"
	"oss.navercorp.com/metis/metis-server/server/rpc"
)

var portOffset = 10000

func TestMain(m *testing.M) {
	s, err := server.New(&server.Config{
		RPC: &rpc.Config{
			Port: server.DefaultRPCPort + portOffset,
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

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := s.Shutdown(true); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
