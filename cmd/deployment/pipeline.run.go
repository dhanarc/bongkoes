package deployment

import (
	"context"
	"fmt"
	"github.com/djk-lgtm/bongkoes/cmd/shared"
	"github.com/djk-lgtm/bongkoes/internal/page"
	"github.com/spf13/cobra"
	"time"
)

var PipelineRunCommand = &cobra.Command{
	Use:   "deployment:pipeline",
	Short: "Run Pipeline",
	Run:   runPipeline,
}

func init() {
	PipelineRunCommand.Flags().StringVarP(&pipelineAlias, "pipeline", "-p", "", "pipeline alias")
}

func runPipeline(_ *cobra.Command, _ []string) {
	cfg := shared.InitConfig()

	dbConnection := shared.InitDatabase(cfg)
	deploymentPlan := page.NewPlan(&page.Opts{
		Config: cfg,
		DBConn: dbConnection,
	})
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	projectCfg := shared.GetProjectConfig()
	alias := projectCfg.PipelineAlias[pipelineAlias]
	fmt.Println(fmt.Sprintf("[bongkoes] Running Pipeline %s - branch %s", alias.Pipeline, alias.Branch))
	pipelineLink, err := deploymentPlan.RunPipelineBranch(ctx, projectCfg.RepositoryName, alias.Branch, alias.Pipeline)
	goPanic(err, "[deployment:issue-diff] failed to running pipeline")

	fmt.Println(fmt.Sprintf("[bongkoes] Pipeline Link:%s", *pipelineLink))
}
