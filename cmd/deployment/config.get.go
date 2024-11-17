package deployment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/djk-lgtm/bongkoes/config"
	"github.com/spf13/cobra"
	"log"
)

var GetConfigCommand = &cobra.Command{
	Use:   "deployment:meta",
	Short: "Get Deployment Metadata",
	Run:   getGetConfig,
}

func init() {
	GetConfigCommand.Flags().StringVarP(&service, "service", "s", "", "service")
}

func getGetConfig(_ *cobra.Command, _ []string) {
	projectConfig := config.GetProjectConfig()
	rawJSON, _ := json.Marshal(projectConfig)

	// Beautify the parsed JSON
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, rawJSON, "", "  "); err != nil {
		log.Fatalf("Error indenting JSON: %v", err)
	}

	// Print the beautified JSON
	fmt.Println(prettyJSON.String())
}
