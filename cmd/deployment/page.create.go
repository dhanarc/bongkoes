package deployment

import (
	"context"
	"fmt"
	"github.com/djk-lgtm/bongkoes/cmd/shared"
	"github.com/djk-lgtm/bongkoes/internal/page"
	"github.com/spf13/cobra"
	"strings"
)

var CreateCommand = &cobra.Command{
	Use:   "deployment:create",
	Short: "Create Deployment Plan",
	Run:   runCreate,
}

func runCreate(_ *cobra.Command, _ []string) {
	cfg := shared.InitConfig()

	deploymentPlan := page.NewPlan(&page.Opts{
		Config: cfg,
	})
	serviceCode, err := readStdIn("Input Service Code: ")
	goPanic(err, "failed to get service code")

	tag, err := readStdIn("Input Tag: ")
	goPanic(err, "failed to get tribe name")

	rollbackTag, err := readStdIn("Input Rollback Tag: ")
	goPanic(err, "failed to get rollback tag")

	deploymentTime, err := readStdIn("Deployment Time (E.g. 12:00 AM): ")
	goPanic(err, "failed to get deployment time")

	downTime, err := readStdIn("Input Down Time Estimation (E.g. 120 Minutes; Optional): ")
	goPanic(err, "failed to get down time estimation")

	published, err := readStdIn("Published?(Y/N)")
	goPanic(err, "failed to get published")

	publishFlag := false
	if strings.ToLower(*published) == "y" {
		publishFlag = true
	}

	link, err := deploymentPlan.InitDocument(context.Background(), page.CreateDeploymentArgs{
		ServiceCode:    *serviceCode,
		Tag:            *tag,
		DeploymentTime: *deploymentTime,
		DownTimeEst:    *downTime,
		RollbackTag:    *rollbackTag,
		Published:      publishFlag,
	})
	goPanic(err, "[deployment:create] failed to execute init document")

	fmt.Println(fmt.Sprintf("[deployment:create] generated deployment document: %s", *link))
}
