package deployment

import (
	"context"
	"fmt"
	"github.com/djk-lgtm/bongkoes/cmd/shared"
	"github.com/djk-lgtm/bongkoes/config"
	"github.com/djk-lgtm/bongkoes/internal/page"
	"github.com/spf13/cobra"
)

var PipelineRunCommand = &cobra.Command{
	Use:   "deployment:pipeline",
	Short: "Run Pipeline",
	Run:   runPipeline,
}

func init() {
	PipelineRunCommand.Flags().StringVarP(&pipelineAlias, "pipeline", "p", "", "pipeline alias")
}

func runPipeline(_ *cobra.Command, _ []string) {
	cfg := shared.InitConfig()

	dbConnection := shared.InitDatabase(cfg)
	deploymentPlan := page.NewPlan(&page.Opts{
		Config: cfg,
		DBConn: dbConnection,
	})

	ctx := context.Background()

	projectCfg := config.GetProjectConfig()
	alias := projectCfg.PipelineMap[pipelineAlias]
	fmt.Println(fmt.Sprintf("[bongkoes] Running Pipeline %s - branch %s", alias.Pipeline, alias.Branch))
	pipelineLink, err := deploymentPlan.RunPipelineBranch(ctx, projectCfg.ServiceCode, alias.Branch, alias.Pipeline)
	goPanic(err, "[deployment:pipeline] failed to running pipeline")

	fmt.Println(fmt.Sprintf("[bongkoes] Pipeline Link:%s", *pipelineLink))
}
