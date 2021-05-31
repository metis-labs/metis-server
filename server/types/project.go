package types

import "github.com/rs/xid"

// Link represents the connection of blocks in the diagram.
type Link struct {
	ID   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
}

// Dependency represents package dependency used in import statements.
type Dependency struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Alias   string `json:"alias"`
	Package string `json:"package"`
}

// NewDefaultDependencies creates a new instance of dependency map with default value.
func NewDefaultDependencies() map[string]*Dependency {
	dependencies := make(map[string]*Dependency)

	torch := &Dependency{
		ID:   xid.New().String(),
		Name: "torch",
	}
	dependencies[torch.ID] = torch

	torchNN := &Dependency{
		ID:    xid.New().String(),
		Name:  "torch.nn",
		Alias: "nn",
	}
	dependencies[torchNN.ID] = torchNN

	return dependencies
}

// Network is data that represents a machine learning network.
type Network struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Dependencies map[string]*Dependency `json:"dependencies"`
	Blocks       map[string]*Block      `json:"blocks"`
	Links        map[string]*Link       `json:"links"`
}

// NewDefaultNetwork creates a new instance of Network with the default value.
func NewDefaultNetwork() *Network {
	blocks := make(map[string]*Block)

	inBlock := &Block{
		ID:            xid.New().String(),
		Name:          "in",
		Type:          InType,
		Position:      &Position{X: 100, Y: 100},
		InitVariables: "",
	}
	blocks[inBlock.ID] = inBlock

	outBlock := &Block{
		ID:            xid.New().String(),
		Name:          "out",
		Type:          OutType,
		Position:      &Position{X: 100, Y: 200},
		InitVariables: "",
	}
	blocks[outBlock.ID] = outBlock

	return &Network{
		ID:           xid.New().String(),
		Name:         "Main",
		Dependencies: NewDefaultDependencies(),
		Blocks:       blocks,
		Links:        make(map[string]*Link),
	}
}

// Project is the project contents of Metis.
type Project struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	Networks map[string]*Network `json:"networks"`
}

// NewProject creates a new instance of Project with the default value.
func NewProject(id string, name string) *Project {
	network := NewDefaultNetwork()
	networks := make(map[string]*Network)
	networks[network.ID] = network

	return &Project{
		ID:       id,
		Name:     name,
		Networks: networks,
	}
}
