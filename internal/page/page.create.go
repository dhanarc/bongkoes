package page

import (
	"bytes"
	"context"
	"fmt"
	"github.com/djk-lgtm/bongkoes/internal/shared"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence"
	"github.com/samber/lo"
	"html/template"
	"os"
	"regexp"
	"strings"
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

	// generate release
	createdVersion, err := d.GenerateVersionRelease(ctx, service, currentTime, args.Tag)
	if err != nil {
		return nil, fmt.Errorf("failed generate version release, error=%v", err)
	}

	// bind issues
	err = d.CollectIssues(ctx, service, createdVersion.ID, args.Tag)
	if err != nil {
		return nil, fmt.Errorf("failed collecting issues, error=%v", err)
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
		return nil, fmt.Errorf("failed rendering content, error=%v", err)
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
		return nil, fmt.Errorf("failed creating deployment plan, error=%v", err)
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

func (d *deploymentPlan) CollectIssues(ctx context.Context, service Service, versionID, tag string) error {
	issueList, err := d.fetchShippedIssue(ctx, service, tag)
	if err != nil {
		return err
	}
	if len(issueList) == 0 {
		return nil
	}

	return d.bindIssueVersion(ctx, issueList, versionID)
}

func (d *deploymentPlan) bindIssueVersion(ctx context.Context, issues []string, jiraID string) error {
	for i := range issues {
		err := d.confluenceAPI.AddIssueFixVersion(ctx, issues[i], jiraID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *deploymentPlan) fetchShippedIssue(ctx context.Context, service Service, newTag string) ([]string, error) {
	tagsListResponse, err := d.bitbucketAPI.GetTagsByDateDesc(ctx, service.ServiceCode.String())
	if err != nil {
		return nil, err
	}
	latestTag := tagsListResponse.Values[0].Name

	destinationPath := "./.shipped_issues"
	d.git.CreateLocalTag(newTag)
	err = d.git.GenerateCommitDiff(latestTag, newTag, destinationPath)
	if err != nil {
		return nil, err
	}

	// load text
	issuesBytes, err := os.ReadFile(destinationPath)
	if err != nil {
		return nil, err
	}

	issuesRawList := string(issuesBytes)
	issueRegex := fmt.Sprintf("%s-\\d+", service.ProjectKey)

	cIssueRegex := regexp.MustCompile(issueRegex)
	issueMatches := cIssueRegex.FindAllString(issuesRawList, -1)

	return lo.Uniq(issueMatches), nil
}

func (d *deploymentPlan) GenerateVersionRelease(ctx context.Context, service Service, releaseTime time.Time, tags string) (*confluence.CreateVersionResponse, error) {
	// get latest version
	latestVersion, err := d.confluenceAPI.GetLatestVersion(ctx, &confluence.FetchLatestVersionRequest{
		Query:      service.ServiceName,
		Status:     confluence.VersionReleased,
		ProjectKey: service.ProjectKey,
	})
	if err != nil {
		return nil, err
	}

	latestVersionReleased, err := time.Parse(shared.DefaultDateYYYYMMDD, latestVersion.ReleaseDate)
	if err != nil {
		return nil, err
	}
	newReleaseStart := latestVersionReleased.AddDate(0, 0, 1)

	// create version
	createdRelease, err := d.confluenceAPI.CreateVersion(ctx, &confluence.CreateVersionRequest{
		Archived:    false,
		Description: "", //TODO::TBD
		Name:        fmt.Sprintf("%s - %s", service.ServiceName, tags),
		ProjectID:   service.ProjectID,
		ReleaseDate: releaseTime.Format(shared.DefaultDateYYYYMMDD),
		StartDate:   newReleaseStart.Format(shared.DefaultDateYYYYMMDD),
		Released:    false,
	})
	if err != nil {
		return nil, err
	}

	jiraLink := fmt.Sprintf("%s/projects/%s/versions/%s/tab/release-report-all-issues", d.cfg.Bongkoes.ConfluenceHost, service.ProjectKey, createdRelease.ID)
	createdRelease.WebLink = jiraLink
	return createdRelease, nil
}
