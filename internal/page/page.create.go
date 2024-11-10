package page

import (
	"bytes"
	"context"
	"fmt"
	"github.com/djk-lgtm/bongkoes/internal/shared"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence"
	"github.com/samber/lo"
	"html/template"
	"time"
)

func (d *deploymentPlan) InitDocument(ctx context.Context, args CreateDeploymentArgs) (*string, error) {
	// get config
	var service Service
	err := d.db.Where("service_code = ?", args.ServiceCode).First(&service).Error
	if err != nil {
		return nil, err
	}

	// get current time
	currentTZ, _ := time.LoadLocation("Asia/Jakarta")
	currentTime := time.Now().In(currentTZ)

	// get template
	templatePage, err := d.confluenceAPI.GetPageByID(ctx, service.TemplateID)
	if err != nil {
		return nil, err
	}

	// render template
	content, err := d.renderContent(RenderArgs{
		Service:           service,
		CurrentTime:       currentTime,
		Tag:               args.Tag,
		RollbackTag:       args.RollbackTag,
		DeploymentTime:    args.DeploymentTime,
		EstimatedDownTime: args.DownTimeEst,
		PageTemplate:      templatePage.Body.Storage.Value,
	})
	if err != nil {
		return nil, err
	}

	// create deployment page
	createPage := new(confluence.CreatePageRequest)
	createPage.SpaceID = templatePage.SpaceID
	createPage.ParentID = service.DeploymentFolderID
	createPage.Title = fmt.Sprintf("%s %s - %s", service.ServiceName, args.Tag, currentTime.Format(shared.DefaultConfluenceTitleTimeLayout))
	status := "draft"
	if args.Published {
		status = "current"
	}
	createPage.Status = status
	createPage.Body = confluence.BodyPage{
		Storage: confluence.BodyStorage{
			Representation: templatePage.Body.Storage.Representation,
			Value:          *content,
		},
	}
	createdPaged, err := d.confluenceAPI.CreatePage(ctx, createPage)
	if err != nil {
		return nil, err
	}

	// compose link page
	links := fmt.Sprintf("%s%s", createdPaged.Links.Base, createdPaged.Links.WebUI)
	return &links, nil
}

func (d *deploymentPlan) renderContent(args RenderArgs) (*string, error) {
	// Parse the template file
	t, err := template.New(args.Service.ServiceCode).Parse(args.PageTemplate)
	if err != nil {
		return nil, err
	}

	// Execute the template with the data
	stringWriter := bytes.NewBufferString("")
	err = t.Execute(stringWriter, DeploymentArgs{
		ServiceCode:    args.Service.ServiceCode,
		ServiceName:    args.Service.ServiceName,
		TeamName:       args.Service.TeamName,
		TribeName:      args.Service.TribeName,
		Tag:            args.Tag,
		DeploymentTime: fmt.Sprintf("%s %s", args.CurrentTime.Format(shared.DefaultConfluenceTitleTimeLayout), args.DeploymentTime),
		DownTimeEst:    args.EstimatedDownTime,
		RollbackTag:    args.RollbackTag,
		JiraLink:       "TBD",
	})
	if err != nil {
		return nil, err
	}
	return lo.ToPtr(stringWriter.String()), nil
}
