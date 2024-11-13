package page

import (
	"context"
	"fmt"
	"github.com/djk-lgtm/bongkoes/internal/shared"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence"
	"time"
)

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
