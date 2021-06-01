package test

import (
	"log"
	"os"
	"testing"

	"oss.navercorp.com/metis/metis-server/server"
	"oss.navercorp.com/metis/metis-server/server/database/mongodb"
	"oss.navercorp.com/metis/metis-server/server/rpc"
	"oss.navercorp.com/metis/metis-server/server/web"
	"oss.navercorp.com/metis/metis-server/server/yorkie"
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
			Addr:         server.DefaultYorkieAddr,
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
