package projects

import (
	"context"

	"github.com/yorkie-team/yorkie/client"
	"github.com/yorkie-team/yorkie/pkg/document"
	"github.com/yorkie-team/yorkie/pkg/document/proxy"

	"oss.navercorp.com/metis/metis-server/internal/log"
	"oss.navercorp.com/metis/metis-server/server/database"
	"oss.navercorp.com/metis/metis-server/server/types"
	"oss.navercorp.com/metis/metis-server/server/yorkie"
)

// Create creates a new project of the given name.
func Create(
	ctx context.Context,
	db database.Database,
	yorkieConf *yorkie.Config,
	projectName string,
) (*types.ProjectInfo, error) {
	projectInfo, err := db.CreateProject(ctx, projectName)
	if err != nil {
		return nil, err
	}

	// TODO(youngteac.hong): Extract yorkie packages such as database.
	cli, err := client.Dial(yorkieConf.RPCAddr, client.Option{
		Token: yorkieConf.WebhookToken,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cli.Close(); err != nil {
			log.Logger.Error(err)
		}
	}()

	if err := cli.Activate(ctx); err != nil {
		return nil, err
	}
	defer func() {
		if err := cli.Deactivate(ctx); err != nil {
			log.Logger.Error(err)
		}
	}()

	doc := document.New(yorkieConf.Collection, projectInfo.ID.String())
	if err := cli.Attach(ctx, doc); err != nil {
		return nil, err
	}

	project := types.NewProject(projectInfo.ID.String(), projectInfo.Name)
	if err := updateProject(doc, project); err != nil {
		return nil, err
	}

	if err := cli.Detach(ctx, doc); err != nil {
		return nil, err
	}

	return projectInfo, nil
}

func updateProject(doc *document.Document, p *types.Project) error {
	return doc.Update(func(root *proxy.ObjectProxy) error {
		// project
		project := root.SetNewObject("project")
		project.SetString("id", p.ID)
		project.SetString("name", p.Name)
		networks := project.SetNewObject("networks")

		// networks
		for nID, n := range p.Networks {
			network := networks.SetNewObject(nID)
			network.SetString("id", n.ID)
			network.SetString("name", n.Name)
			dependencies := network.SetNewObject("dependencies")
			thirdPartyDeps := dependencies.SetNewObject("thirdPartyDeps")
			blocks := network.SetNewObject("blocks")
			links := network.SetNewObject("links")

			// dependencies
			for dID, d := range n.Dependencies.ThirdPartyDeps {
				dependency := thirdPartyDeps.SetNewObject(dID)
				dependency.SetString("id", d.ID)
				dependency.SetString("name", d.Name)
				if d.Alias != "" {
					dependency.SetString("alias", d.Alias)
				}
				if d.Package != "" {
					dependency.SetString("package", d.Package)
				}
			}

			// blocks
			for bID, b := range n.Blocks {
				block := blocks.SetNewObject(bID)
				block.SetString("id", b.ID)
				block.SetString("name", b.Name)
				block.SetString("type", string(b.Type))
				position := block.SetNewObject("position")
				position.SetInteger("x", b.Position.X)
				position.SetInteger("y", b.Position.Y)

				if b.Type == types.InType {
					block.SetString("initVariables", b.InitVariables)
				} else if b.Type == types.NetworkType {
					block.SetString("refNetwork", b.RefNetwork)
					block.SetInteger("repeats", b.Repeats)
					updateParameters(block, b.Parameters)
				} else {
					block.SetInteger("repeats", b.Repeats)
					updateParameters(block, b.Parameters)
				}
			}

			// links
			for nID, n := range n.Links {
				link := links.SetNewObject(nID)
				link.SetString("id", n.ID)
				link.SetString("from", n.From)
				link.SetString("to", n.To)
			}
		}

		return nil
	})
}

func updateParameters(block *proxy.ObjectProxy, params types.Parameters) {
	parameters := block.SetNewObject("parameters")
	for pID, p := range params {
		switch v := p.(type) {
		case string:
			parameters.SetString(pID, v)
		case int:
			parameters.SetInteger(pID, v)
		case bool:
			parameters.SetBool(pID, v)
		default:
			log.Logger.Fatal("unsupported type")
		}
	}
}
