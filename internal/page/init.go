package page

import (
	"context"
	"github.com/djk-lgtm/atlassianoto/config"
	"github.com/djk-lgtm/atlassianoto/pkg/atlassian/confluence"
	"gorm.io/gorm"
)

type Plan interface {
	InitConfig(context.Context, CreateServiceArgs) error
	GetConfig(context.Context, string) error
	InitDocument(context.Context, CreateDeploymentArgs) (*string, error)
}

type deploymentPlan struct {
	confluenceAPI confluence.API
	cfg           *config.Config
	db            *gorm.DB
}

type Opts struct {
	Config *config.Config
	DBConn *gorm.DB
}

func NewPlan(o *Opts) Plan {
	confluenceAPI := confluence.NewConfluenceAPI(&confluence.Opts{
		ConfluenceHost: o.Config.Atlassianoto.ConfluenceHost,
		Email:          o.Config.Atlassianoto.AtlassianEmail,
		Token:          o.Config.Atlassianoto.AtlassianToken,
	})
	return &deploymentPlan{
		cfg:           o.Config,
		db:            o.DBConn,
		confluenceAPI: confluenceAPI,
	}
}
