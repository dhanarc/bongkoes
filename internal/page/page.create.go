package page

import (
	"bytes"
	"context"
	"fmt"
	"github.com/djk-lgtm/bongkoes/internal/shared"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence/types"
	"github.com/samber/lo"
	"html/template"
	"regexp"
	"strings"
	"time"
)

func (d *deploymentPlan) InitDocument(ctx context.Context, tag, deploymentTime, downTimeEst string, published bool) (*string, error) {

	// get current time
	currentTZ, _ := time.LoadLocation("Asia/Jakarta")
	currentTime := time.Now().In(currentTZ)

	// get template
	templatePage, err := d.confluenceAPI.GetPageByID(ctx, fmt.Sprintf("%d", d.projectCfg.DeploymentTemplateID))
	if err != nil {
		return nil, fmt.Errorf("[InitDocument] failed to get template page, error=%+v", err)
	}

	// get latest version
	latestVersion, err := d.confluenceAPI.GetLatestVersion(ctx, &types.FetchLatestVersionRequest{
		Query:      d.projectCfg.ServiceName,
		Status:     types.VersionReleased,
		ProjectKey: d.projectCfg.JiraProjectKey,
	})
	if err != nil {
		return nil, err
	}

	// generate release
	createdVersion, err := d.GenerateVersionRelease(ctx, latestVersion, currentTime, tag)
	if err != nil {
		return nil, fmt.Errorf("[InitDocument] failed to generate version release, error=%+v", err)
	}

	// bind issues
	err = d.CollectIssues(ctx, createdVersion.ID, tag)
	if err != nil {
		return nil, err
	}

	currentUser, err := d.confluenceAPI.GetCurrentUserMention(ctx)
	if err != nil {
		return nil, err
	}

	issueJQL := fmt.Sprintf("project = %s AND fixversion = \"%s\" ORDER BY fixVersion, Rank ASC", d.projectCfg.JiraProjectKey, createdVersion.Name)

	// render template
	content, err := d.renderContent(RenderArgs{
		CurrentTime:       currentTime,
		Tag:               tag,
		RollbackTag:       d.getTagByVersionName(latestVersion.Name),
		DeploymentTime:    deploymentTime,
		EstimatedDownTime: downTimeEst,
		PageTemplate:      templatePage.Body.Storage.Value,
		JiraLink:          createdVersion.WebLink,
		CurrentUser:       *currentUser,
		IssueQuery:        issueJQL,
	})
	if err != nil {
		return nil, fmt.Errorf("[InitDocument] failed to rendering content page, error=%+v", err)
	}

	// create deployment page
	createPage := new(types.CreatePageRequest)
	createPage.SpaceID = templatePage.SpaceID
	createPage.ParentID = fmt.Sprintf("%d", d.projectCfg.DeploymentPlanFolderID)
	createPage.Title = fmt.Sprintf("%s %s - %s", d.projectCfg.ServiceName, tag, currentTime.Format(shared.DefaultConfluenceTitleTimeLayout))
	status := "draft"
	if published {
		status = "current"
	}
	createPage.Status = status
	createPage.Body = types.BodyPage{
		Storage: types.BodyStorage{
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

func (d *deploymentPlan) getTagByVersionName(versionName string) string {
	regexVersionName := regexp.MustCompile(`([\w\s]+) v(\d+.\d+.\d+)`)
	matches := regexVersionName.FindStringSubmatch(versionName)
	if len(matches) > 0 {
		version := matches[2]
		return version
	}
	return ""
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

	issueLink, err := d.confluenceAPI.GenerateJiraLink(args.IssueQuery, d.GetViewOptionsIssue())
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
		Me:             args.CurrentUser,
		IssueTable:     *issueLink,
	})
	if err != nil {
		return nil, err
	}
	return lo.ToPtr(stringWriter.String()), nil
}

func (d *deploymentPlan) GetViewOptionsIssue() []types.JiraLinkView {
	return []types.JiraLinkView{
		{
			Type: "table",
			Properties: types.JiraViewProperty{
				Columns: []types.JiraPropertyKey{
					{
						Key:   "key",
						Width: 91,
					},
					{
						Key:       "summary",
						IsWrapped: true,
					},
					{
						Key:       "assignee",
						Width:     148,
						IsWrapped: true,
					},
					{
						Key:       "customfield_11620",
						Width:     304,
						IsWrapped: true,
					},
					{
						Key: "status",
					},
				},
			},
		},
	}
}
