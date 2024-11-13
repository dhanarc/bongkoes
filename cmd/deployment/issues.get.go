package deployment

import (
	"context"
	"github.com/djk-lgtm/bongkoes/cmd/shared"
	"github.com/djk-lgtm/bongkoes/internal/page"
	"github.com/spf13/cobra"
	"time"
)

var GetLatestIssuesCommand = &cobra.Command{
	Use:   "deployment:latest-issues",
	Short: "Get Latest Issues",
	Run:   getIssueDiff,
}

func init() {
	GetConfigCommand.Flags().StringVarP(&service, "service", "s", "", "service")
	GetLatestIssuesCommand.Flags().StringVarP(&tag, "tag", "t", "", "tag")
	GetLatestIssuesCommand.Flags().StringVarP(&latestTag, "previous-tag", "p", "", "previous tag")
}

func getIssueDiff(_ *cobra.Command, _ []string) {
	cfg := shared.InitConfig()

	dbConnection := shared.InitDatabase(cfg)
	deploymentPlan := page.NewPlan(&page.Opts{
		Config: cfg,
		DBConn: dbConnection,
	})
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := deploymentPlan.GetIssueListDiff(ctx, service, latestTag, tag)
	if err != nil {
		panic("failed to load issues")
	}
}
