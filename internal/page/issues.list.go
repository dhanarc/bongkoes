package page

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"os"
	"regexp"
)

func (d *deploymentPlan) CollectIssues(ctx context.Context, versionID, tag string) error {
	tagsListResponse, err := d.bitbucketAPI.GetTagsByDateDesc(ctx, d.projectCfg.ServiceCode)
	if err != nil {
		return fmt.Errorf("[CollectIssues] failed to execute GetTagsByDateDesc, error:%+v", err)
	}
	latestTag := tagsListResponse.Values[0].Name
	err = d.git.CreateLocalTag(tag)
	if err != nil {
		return err
	}
	issueList, err := d.getIssueDiff(latestTag, tag)
	if err != nil {
		return fmt.Errorf("[CollectIssues] failed to get issue diff, error:%+v", err)
	}
	if len(issueList) == 0 {
		return nil
	}

	return d.bindIssueVersion(ctx, issueList, versionID)
}

func (d *deploymentPlan) getIssueDiff(previousTag, newTag string) ([]string, error) {
	destinationPath := fmt.Sprintf("./.%s", newTag)
	err := d.git.GenerateCommitDiff(previousTag, newTag, destinationPath)
	if err != nil {
		return nil, err
	}

	// load text
	issuesBytes, err := os.ReadFile(destinationPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = os.Remove(destinationPath)
		if err != nil {
			fmt.Println("failed to remove")
		}
	}()

	issuesRawList := string(issuesBytes)
	issueRegex := fmt.Sprintf("%s-\\d+", d.projectCfg.JiraProjectKey)

	cIssueRegex := regexp.MustCompile(issueRegex)
	issueMatches := cIssueRegex.FindAllString(issuesRawList, -1)

	return lo.Uniq(issueMatches), nil
}

func (d *deploymentPlan) bindIssueVersion(ctx context.Context, issues []string, jiraID string) error {
	for i := range issues {
		err := d.confluenceAPI.AddIssueFixVersion(ctx, issues[i], jiraID)
		if err != nil {
			return fmt.Errorf("[BindIssueVersion] failed to add issue fix version, error=%+v", err)
		}
	}
	return nil
}
