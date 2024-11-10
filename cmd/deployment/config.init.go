package deployment

import (
	"context"
	"github.com/djk-lgtm/bongkoes/cmd/shared"
	"github.com/djk-lgtm/bongkoes/internal/page"
	"github.com/spf13/cobra"
)

var InitCommand = &cobra.Command{
	Use:   "deployment:init",
	Short: "Deployment Plan Document Init Config",
	Run:   runInit,
}

func runInit(_ *cobra.Command, _ []string) {
	cfg := shared.InitConfig()

	dbConnection := shared.InitDatabase(cfg)

	deploymentPlan := page.NewPlan(&page.Opts{
		Config: cfg,
		DBConn: dbConnection,
	})
	teamName, err := readStdIn("Input Team Name: ")
	goPanic(err, "failed to get team name")

	tribeName, err := readStdIn("Input Tribe Name: ")
	goPanic(err, "failed to get tribe name")

	serviceCode, err := readStdIn("Input Service Code: ")
	goPanic(err, "failed to get service code")

	serviceName, err := readStdIn("Input Service Name: ")
	goPanic(err, "failed to get service name")

	templateID, err := readStdIn("Input Template ID(Confluence Page ID): ")
	goPanic(err, "failed to get template id")

	folderPageID, err := readStdIn("Input Folder Page ID(Confluence Page ID): ")
	goPanic(err, "failed to get folder page id")

	err = deploymentPlan.InitConfig(context.Background(), page.CreateServiceArgs{
		TeamName:           *teamName,
		TribeName:          *tribeName,
		ServiceCode:        *serviceCode,
		ServiceName:        *serviceName,
		TemplateID:         *templateID,
		DeploymentFolderID: *folderPageID,
	})
	goPanic(err, "failed to init deployment plan config")
}
