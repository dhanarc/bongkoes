package page

import (
	"context"
	"strconv"
)

func (d *deploymentPlan) InitConfig(ctx context.Context, args CreateServiceArgs) error {
	projectDetail, err := d.confluenceAPI.GetProjectDetail(ctx, args.ProjectKey)
	if err != nil {
		return err
	}

	projectID, err := strconv.ParseUint(projectDetail.ID, 10, 64)
	if err != nil {
		return err
	}

	tx := d.db.WithContext(ctx).Create(&Service{
		TribeName:          args.TribeName,
		TeamName:           args.TeamName,
		ProjectKey:         args.ProjectKey,
		ProjectID:          projectID,
		ServiceCode:        ServiceCode(args.ServiceCode),
		ServiceName:        args.ServiceName,
		TemplateID:         args.TemplateID,
		DeploymentFolderID: args.DeploymentFolderID,
	})

	return tx.Error
}
