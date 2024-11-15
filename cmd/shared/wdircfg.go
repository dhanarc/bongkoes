package shared

import (
	"github.com/spf13/viper"
	"log"
	"regexp"
	"strings"
)

type ProjectConfig struct {
	RepositoryName string
	PipelineAlias  map[string]PipelineAlias
}

type PipelineAlias struct {
	Branch   string
	Pipeline string
}

func GetProjectConfig() *ProjectConfig {
	viper.SetConfigFile(".bongkoes")
	viper.SetConfigType("dotenv")

	// Read in the .env file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Access the environment variables using Viper
	repositoryName := viper.GetString("SERVICE_CODE")

	// deployStagingAli:staging[master];deployProduction:production[master]
	regexAlias := regexp.MustCompile(`(\w+):(\w+)\[(\w+)\]`)
	pipelineAlias := viper.GetString("PIPELINE_ALIAS")

	alias := strings.Split(pipelineAlias, ";")
	pipelines := make(map[string]PipelineAlias)
	if len(alias) > 0 {
		for i := range alias {
			match := regexAlias.FindStringSubmatch(alias[i])
			if len(match) > 1 {
				pipeline := match[1]
				pAlias := match[2]
				branch := match[3]
				p := PipelineAlias{
					Branch:   branch,
					Pipeline: pipeline,
				}
				pipelines[pAlias] = p
			}
		}
	}
	return &ProjectConfig{
		RepositoryName: repositoryName,
		PipelineAlias:  pipelines,
	}
}
