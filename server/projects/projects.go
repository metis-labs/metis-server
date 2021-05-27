package projects

import (
	"context"

	"github.com/yorkie-team/yorkie/client"
	"github.com/yorkie-team/yorkie/pkg/document"

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
) (*types.Project, error) {
	project, err := db.CreateProject(ctx, projectName)
	if err != nil {
		return nil, err
	}

	// TODO(youngteac.hong): Extract yorkie packages such as database.
	cli, err := client.Dial(yorkieConf.Addr, client.Option{
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

	doc := document.New(yorkieConf.Collection, project.ID.String())
	if err := cli.Attach(ctx, doc); err != nil {
		return nil, err
	}

	// TODO(youngteac.hong): Updates the basic data of the project content or template data.
	// using `doc.Update()`.

	if err := cli.Detach(ctx, doc); err != nil {
		return nil, err
	}

	return project, nil
}
