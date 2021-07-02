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

package types

import "github.com/rs/xid"

// BlockType is a type of block.
type BlockType string

// Belows are the types of the block.
// TODO(youngteac.hong): We need to find a way to simplify types for Normal Block.
// It is burdensome to add a type every time we add a normal block.
const (
	InType          BlockType = "In"
	OutType         BlockType = "Out"
	NetworkType     BlockType = "Network"
	Conv2dType      BlockType = "Conv2d"
	BatchNorm2dType BlockType = "BachNorm2d"
	ReLUType        BlockType = "ReLU"
	MaxPool2dType   BlockType = "MaxPool2d"
)

// Position represents point on the canvas.
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ParameterValue is a value for the parameter.
type ParameterValue interface{}

// Parameters is the parameter value map of the block.
type Parameters map[string]ParameterValue

// Block is the block component of the network.
type Block struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Type     BlockType `json:"type"`
	Position *Position `json:"position"`

	// for IOBlock
	InitVariables string `json:"initVariables"`

	// for NetworkBlock
	RefNetwork string `json:"refNetwork"`

	// for NormalBlock + NetworkBlock
	Repeats    int        `json:"repeats"`
	Parameters Parameters `json:"parameters"`
}

// NewBlock creates a new instance of Block.
func NewBlock(blockType BlockType, name string) *Block {
	return &Block{
		ID:   xid.New().String(),
		Name: name,
		Type: blockType,
	}
}
