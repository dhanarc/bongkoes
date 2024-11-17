package deployment

import (
	"context"
	"github.com/djk-lgtm/bongkoes/cmd/shared"
	"github.com/djk-lgtm/bongkoes/config"
	"github.com/djk-lgtm/bongkoes/internal/page"
	"github.com/spf13/cobra"
)

var DebugCommand = &cobra.Command{
	Use:   "deployment:debug",
	Short: "Deployment Debug",
	Run:   runDebug,
}

func runDebug(_ *cobra.Command, _ []string) {
	cfg := shared.InitConfig()
	projectCfg := config.GetProjectConfig()

	dbConnection := shared.InitDatabase(cfg)
	deploymentPlan := page.NewPlan(&page.Opts{
		Config:        cfg,
		ProjectConfig: projectCfg,
		DBConn:        dbConnection,
	})
	err := deploymentPlan.Debug(context.Background())
	goPanic(err, "[deployment:debug] failed to run debug")
}
