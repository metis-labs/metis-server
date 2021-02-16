package client_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"oss.navercorp.com/metis/metis-server/client"
	"oss.navercorp.com/metis/metis-server/server"
)

func TestMain(m *testing.M) {
	s, err := server.New()
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

func TestNewAndClose(t *testing.T) {
	c, err := client.New()
	assert.NoError(t, err)

	err = c.Close()
	assert.NoError(t, err)
}

func TestCreateDiagram(t *testing.T) {
	c, err := client.New()
	assert.NoError(t, err)
	defer func() {
		err = c.Close()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	testDiagramName := "HelloWorld"

	diagram, err := c.CreateDiagram(ctx, testDiagramName)
	assert.NoError(t, err)

	assert.Equal(t, testDiagramName, diagram.Name)
}
