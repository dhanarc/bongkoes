package page

import (
	"bytes"
	"context"
	"fmt"
	"github.com/djk-lgtm/bongkoes/internal/shared"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence"
	"github.com/samber/lo"
	"html/template"
	"strings"
	"time"
)

func (d *deploymentPlan) InitDocument(ctx context.Context, args CreateDeploymentArgs) (*string, error) {
	// get config
	var service Service
	err := d.db.Where("service_code = ?", args.ServiceCode).First(&service).Error
	if err != nil {
		return nil, fmt.Errorf("[InitDocument] failed to fetch service config, error=%+v", err)
	}

	// get current time
	currentTZ, _ := time.LoadLocation("Asia/Jakarta")
	currentTime := time.Now().In(currentTZ)

	// get template
	templatePage, err := d.confluenceAPI.GetPageByID(ctx, service.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("[InitDocument] failed to get template page, error=%+v", err)
	}

	// generate release
	createdVersion, err := d.GenerateVersionRelease(ctx, service, currentTime, args.Tag)
	if err != nil {
		return nil, fmt.Errorf("[InitDocument] failed to generate version release, error=%+v", err)
	}

	// bind issues
	err = d.CollectIssues(ctx, service, createdVersion.ID, args.Tag)
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
		JiraLink:          createdVersion.WebLink,
	})
	if err != nil {
		return nil, fmt.Errorf("[InitDocument] failed to rendering content page, error=%+v", err)
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
		return nil, fmt.Errorf("[InitDocument] failed to create confluence page, error=%+v", err)
	}

	// compose link page
	links := fmt.Sprintf("%s%s", createdPaged.Links.Base, createdPaged.Links.WebUI)
	return &links, nil
}

func (d *deploymentPlan) injectJiraLinkTemplate(template string) string {
	injectAppearanceCard := strings.ReplaceAll(template, shared.TemplateJiraLink, fmt.Sprintf("%s %s", shared.TemplateJiraLink, shared.JiraLinkProps))
	return strings.ReplaceAll(injectAppearanceCard, shared.TemplateJira, "{{ .JiraLink }}")
}

func (d *deploymentPlan) renderContent(args RenderArgs) (*string, error) {
	// Parse the template file
	contentTemplate := d.injectJiraLinkTemplate(args.PageTemplate)
	t, err := template.New(args.Service.ServiceCode.String()).Parse(contentTemplate)
	if err != nil {
		return nil, err
	}

	// Execute the template with the data
	stringWriter := bytes.NewBufferString("")
	err = t.Execute(stringWriter, DeploymentArgs{
		ServiceCode:    args.Service.ServiceCode.String(),
		ServiceName:    args.Service.ServiceName,
		TeamName:       args.Service.TeamName,
		TribeName:      args.Service.TribeName,
		Tag:            args.Tag,
		DeploymentTime: fmt.Sprintf("%s %s", args.CurrentTime.Format(shared.DefaultConfluenceTitleTimeLayout), args.DeploymentTime),
		DownTimeEst:    args.EstimatedDownTime,
		RollbackTag:    args.RollbackTag,
		JiraLink:       args.JiraLink,
	})
	if err != nil {
		return nil, err
	}
	return lo.ToPtr(stringWriter.String()), nil
}
