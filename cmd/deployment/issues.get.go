package deployment

import (
	"context"
	"github.com/djk-lgtm/bongkoes/cmd/shared"
	"github.com/djk-lgtm/bongkoes/internal/page"
	"github.com/spf13/cobra"
	"time"
)

var GetLatestIssuesCommand = &cobra.Command{
	Use:   "deployment:issue-diff",
	Short: "Get Issue Diff",
	Run:   getIssueDiff,
}

func init() {
	GetLatestIssuesCommand.Flags().StringVarP(&service, "service", "s", "", "service")
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
	goPanic(err, "[deployment:issue-diff] failed to get issue list diff")
}
