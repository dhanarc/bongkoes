package deployment

import (
	"context"
	"github.com/djk-lgtm/bongkoes/cmd/shared"
	"github.com/djk-lgtm/bongkoes/internal/page"
	"github.com/spf13/cobra"
	"time"
)

var GetConfigCommand = &cobra.Command{
	Use:   "deployment:get-config",
	Short: "Get Deployment Config",
	Run:   getGetConfig,
}

func init() {
	GetConfigCommand.Flags().StringVarP(&service, "service", "s", "", "service")
}

func getGetConfig(_ *cobra.Command, _ []string) {
	cfg := shared.InitConfig()

	dbConnection := shared.InitDatabase(cfg)
	deploymentPlan := page.NewPlan(&page.Opts{
		Config: cfg,
		DBConn: dbConnection,
	})
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := deploymentPlan.GetConfig(ctx, service)
	if err != nil {
		panic("failed to get config")
	}
}
