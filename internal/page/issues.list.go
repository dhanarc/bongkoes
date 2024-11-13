package page

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"os"
	"regexp"
)

func (d *deploymentPlan) GetIssueListDiff(ctx context.Context, serviceCode, previousTag, tag string) error {
	// get config
	var service Service
	err := d.db.Where("service_code = ?", serviceCode).First(&service).Error
	if err != nil {
		return err
	}

	issues, err := d.getIssueDiff(ctx, service, previousTag, tag)
	if err != nil {
		return err
	}

	fmt.Println(issues)

	return nil
}

func (d *deploymentPlan) getIssueDiff(ctx context.Context, service Service, previousTag, newTag string) ([]string, error) {
	destinationPath := fmt.Sprintf("./.%s", newTag)
	d.git.CreateLocalTag(newTag)
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
	issueRegex := fmt.Sprintf("%s-\\d+", service.ProjectKey)

	cIssueRegex := regexp.MustCompile(issueRegex)
	issueMatches := cIssueRegex.FindAllString(issuesRawList, -1)

	return lo.Uniq(issueMatches), nil
}
