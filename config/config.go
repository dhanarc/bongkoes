package config

type Config struct {
	Atlassianoto struct {
		DBLocation     string `mapstructure:"db_location"`
		AtlassianEmail string `mapstructure:"atlassian_email"`
		AtlassianToken string `mapstructure:"atlassian_token"`
		ConfluenceHost string `mapstructure:"confluence_host"`
	} `mapstructure:"atlassianoto"`
}
