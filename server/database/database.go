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

package database

import (
	"context"
	"errors"

	"github.com/metis-labs/metis-server/server/types"
)

var (
	// ErrInvalidID is returned when the given string is ID.
	ErrInvalidID = errors.New("invalid id")

	// ErrNotFound is returned when the requested resource cannot be found.
	ErrNotFound = errors.New("resource not found")
)

// Database represents database which reads or saves Metis data.
type Database interface {
	Dial(ctx context.Context) error
	Close(ctx context.Context) error

	CreateProject(ctx context.Context, name string) (*types.ProjectInfo, error)
	FindProject(ctx context.Context, id types.ID) (*types.ProjectInfo, error)
	ListProjects(ctx context.Context) ([]*types.ProjectInfo, error)
	UpdateProject(ctx context.Context, id types.ID, name string) error
	DeleteProject(ctx context.Context, id types.ID) error

	CreateTemplate(ctx context.Context, name, contents string) (*types.TemplateInfo, error)
	FindTemplate(ctx context.Context, id types.ID) (*types.TemplateInfo, error)
}
