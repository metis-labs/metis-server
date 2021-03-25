package client_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func TestClient(t *testing.T) {
	t.Run("new/close test", func(t *testing.T) {
		c, err := client.New()
		assert.NoError(t, err)

		err = c.Close()
		assert.NoError(t, err)
	})
}

func TestProject(t *testing.T) {
	cli, err := client.New()
	assert.NoError(t, err)
	defer func() {
		err = cli.Close()
		assert.NoError(t, err)
	}()

	t.Run("create/update/delete project test", func(t *testing.T) {
		pbProject, err := cli.CreateProject(context.Background(), t.Name())
		assert.NoError(t, err)
		assert.Equal(t, t.Name(), pbProject.Name)

		ctx := context.Background()

		projects, err := cli.ListProjects(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, projects)

		err = cli.UpdateProject(ctx, pbProject.Id, "updated")
		assert.NoError(t, err)

		err = cli.UpdateProject(ctx, "invalid", "updated")
		assert.Equal(t, codes.InvalidArgument, status.Convert(err).Code())

		err = cli.UpdateProject(ctx, "000000000000000000000000", "updated")
		assert.Equal(t, codes.NotFound, status.Convert(err).Code())

		err = cli.DeleteProject(context.Background(), pbProject.Id)
		assert.NoError(t, err)

		err = cli.DeleteProject(context.Background(), pbProject.Id)
		assert.NoError(t, err)
	})
}
