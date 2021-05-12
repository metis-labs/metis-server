package client

import (
	"log"
	"os"
	"testing"

	"oss.navercorp.com/metis/metis-server/server"
)

func TestMain(m *testing.M) {
	s, err := server.New(server.NewConfig())
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
