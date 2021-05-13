package client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"oss.navercorp.com/metis/metis-server/client"
)

const (
	testUserA = "KR18401"
	testUserB = "KR18817"
)

func TestClient(t *testing.T) {
	t.Run("new/close test", func(t *testing.T) {
		cliA, err := client.Dial(testServer.RPCAddr(), client.Option{UserID: testUserA})
		assert.NoError(t, err)
		defer func() {
			err = cliA.Close()
			assert.NoError(t, err)
		}()
	})

	t.Run("invalid token test", func(t *testing.T) {
		cli, err := client.Dial(testServer.RPCAddr(), client.Option{UserID: ""})
		assert.NoError(t, err)

		_, err = cli.ListProjects(context.Background())
		assert.Equal(t, codes.Unauthenticated, status.Convert(err).Code())
	})
}

func TestProject(t *testing.T) {
	cliA, err := client.Dial(testServer.RPCAddr(), client.Option{UserID: testUserA})
	assert.NoError(t, err)
	defer func() {
		err = cliA.Close()
		assert.NoError(t, err)
	}()
	cliB, err := client.Dial(testServer.RPCAddr(), client.Option{UserID: testUserB})
	assert.NoError(t, err)
	defer func() {
		err = cliB.Close()
		assert.NoError(t, err)
	}()

	t.Run("list project test", func(t *testing.T) {
		ctxA := context.Background()
		ctxB := context.Background()

		pbProject, err := cliA.CreateProject(ctxA, t.Name())
		assert.NoError(t, err)
		assert.Equal(t, t.Name(), pbProject.Name)

		defer func() {
			err = cliA.DeleteProject(ctxA, pbProject.Id)
			assert.NoError(t, err)
		}()

		projects, err := cliB.ListProjects(ctxB)
		assert.NoError(t, err)
		assert.Empty(t, projects)

		projects, err = cliA.ListProjects(ctxA)
		assert.NoError(t, err)
		assert.NotEmpty(t, projects)
	})

	t.Run("update project test", func(t *testing.T) {
		ctxA := context.Background()

		pbProject, err := cliA.CreateProject(ctxA, t.Name())
		assert.NoError(t, err)
		assert.Equal(t, t.Name(), pbProject.Name)

		defer func() {
			err = cliA.DeleteProject(context.Background(), pbProject.Id)
			assert.NoError(t, err)
		}()

		err = cliA.UpdateProject(ctxA, pbProject.Id, "updated")
		assert.NoError(t, err)

		err = cliB.UpdateProject(ctxA, pbProject.Id, "updated")
		assert.Equal(t, codes.NotFound, status.Convert(err).Code())

		err = cliA.UpdateProject(ctxA, "invalid", "updated")
		assert.Equal(t, codes.InvalidArgument, status.Convert(err).Code())

		err = cliA.UpdateProject(ctxA, "000000000000000000000000", "updated")
		assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	})

	t.Run("delete project test", func(t *testing.T) {
		ctxA := context.Background()

		pbProject, err := cliA.CreateProject(ctxA, t.Name())
		assert.NoError(t, err)
		assert.Equal(t, t.Name(), pbProject.Name)

		err = cliA.DeleteProject(ctxA, pbProject.Id)
		assert.NoError(t, err)

		err = cliA.DeleteProject(ctxA, pbProject.Id)
		assert.NoError(t, err)
	})
}
