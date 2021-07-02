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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/metis-labs/metis-server/client"
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
