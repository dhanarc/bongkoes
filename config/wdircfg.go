package config

import (
	"github.com/spf13/viper"
	"log"
	"regexp"
	"strings"
)

type ProjectConfig struct {
	TribeName              string `mapstructure:"TRIBE_NAME" json:"tribe_name"`
	TeamName               string `mapstructure:"TEAM_NAME" json:"team_name"`
	ServiceCode            string `mapstructure:"SERVICE_CODE" json:"service_code"`
	ServiceName            string `mapstructure:"SERVICE_NAME" json:"service_name"`
	JiraProjectKey         string `mapstructure:"JIRA_PROJECT_KEY" json:"jira_project_key"`
	JiraProjectID          string `mapstructure:"JIRA_PROJECT_ID" json:"jira_project_id"`
	DeploymentTemplateID   uint64 `mapstructure:"DEPLOYMENT_TEMPLATE_ID" json:"deployment_template_id"`
	DeploymentPlanFolderID uint64 `mapstructure:"DEPLOYMENT_PLAN_FOLDER_ID" json:"deployment_parent_id"`
	MainBranch             string `mapstructure:"MAIN_BRANCH" json:"main_branch"`
	PipelineAlias          string `mapstructure:"PIPELINE_ALIAS" json:"-"`

	PipelineMap map[string]PipelineAlias `mapstructure:"-" json:"pipeline_alias"`
}

type PipelineAlias struct {
	Branch   string `json:"branch"`
	Pipeline string `json:"pipeline"`
}

func GetProjectConfig() *ProjectConfig {
	viper.SetConfigFile(".bongkoes")
	viper.SetConfigType("dotenv")

	// Read in the .env file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var projectConfig ProjectConfig
	if err := viper.Unmarshal(&projectConfig); err != nil {
		log.Fatalf("failed to get project config")
	}

	// deployStagingAli:staging[master];deployProduction:production[master]
	regexAlias := regexp.MustCompile(`(\w+):(\w+)\[(\w+)\]`)

	alias := strings.Split(projectConfig.PipelineAlias, ";")
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

	projectConfig.PipelineMap = pipelines
	return &projectConfig
}
