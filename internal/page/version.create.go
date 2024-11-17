package page

import (
	"context"
	"fmt"
	"github.com/djk-lgtm/bongkoes/internal/shared"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence/types"
	"strconv"
	"time"
)

func (d *deploymentPlan) GenerateVersionRelease(ctx context.Context, previousVersion *types.Version, releaseTime time.Time, tags string) (*types.CreateVersionResponse, error) {
	latestVersionReleased, err := time.Parse(shared.DefaultDateYYYYMMDD, previousVersion.ReleaseDate)
	if err != nil {
		return nil, err
	}
	newReleaseStart := latestVersionReleased.AddDate(0, 0, 1)

	// create version
	jiraProjectID, err := strconv.ParseUint(d.projectCfg.JiraProjectID, 10, 64)
	if err != nil {
		return nil, err
	}
	createdRelease, err := d.confluenceAPI.CreateVersion(ctx, &types.CreateVersionRequest{
		Archived:    false,
		Description: "", //TODO::TBD
		Name:        fmt.Sprintf("%s - %s", d.projectCfg.ServiceName, tags),
		ProjectID:   jiraProjectID,
		ReleaseDate: releaseTime.Format(shared.DefaultDateYYYYMMDD),
		StartDate:   newReleaseStart.Format(shared.DefaultDateYYYYMMDD),
		Released:    false,
	})
	if err != nil {
		return nil, err
	}

	jiraLink := fmt.Sprintf("%s/projects/%s/versions/%s/tab/release-report-all-issues", d.cfg.Bongkoes.ConfluenceHost, d.projectCfg.JiraProjectKey, createdRelease.ID)
	createdRelease.WebLink = jiraLink
	return createdRelease, nil
}
