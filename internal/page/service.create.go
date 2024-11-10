package page

import (
	"context"
)

func (d *deploymentPlan) InitConfig(ctx context.Context, args CreateServiceArgs) error {
	tx := d.db.WithContext(ctx).Create(&Service{
		TribeName:          args.TribeName,
		TeamName:           args.TeamName,
		ServiceCode:        args.ServiceCode,
		ServiceName:        args.ServiceName,
		TemplateID:         args.TemplateID,
		DeploymentFolderID: args.DeploymentFolderID,
	})
	return tx.Error
}
