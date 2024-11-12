package config

type Config struct {
	Bongkoes BongkoesConfig `mapstructure:"bongkoes"`
}

type BongkoesConfig struct {
	DBLocation           string `mapstructure:"db_location"`
	AtlassianEmail       string `mapstructure:"atlassian_email"`
	AtlassianToken       string `mapstructure:"atlassian_token"`
	ConfluenceHost       string `mapstructure:"confluence_host"`
	BitbucketUsername    string `mapstructure:"bitbucket_username"`
	BitbucketAppPassword string `mapstructure:"bitbucket_app_password"`
	BitbucketWorkspace   string `mapstructure:"bitbucket_workspace"`
}
