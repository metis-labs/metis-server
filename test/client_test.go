package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"oss.navercorp.com/metis/metis-server/client"
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
