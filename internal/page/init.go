package page

import (
	"context"
	"github.com/djk-lgtm/bongkoes/config"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence"
	"github.com/djk-lgtm/bongkoes/pkg/bitbucket"
	"github.com/djk-lgtm/bongkoes/pkg/git"
	"gorm.io/gorm"
)

type Plan interface {
	GetConfig(context.Context, string) error
	InitDocument(context.Context, CreateDeploymentArgs) (*string, error)
	GetIssueListDiff(context.Context, string, string, string) error
	RunPipelineBranch(ctx context.Context, serviceCode, branch, pipeline string) (*string, error)
	Debug(context.Context) error
}

type deploymentPlan struct {
	confluenceAPI confluence.API
	bitbucketAPI  bitbucket.API
	git           git.LocalGit
	cfg           *config.Config
	projectCfg    *config.ProjectConfig
	db            *gorm.DB
}

type Opts struct {
	Config        *config.Config
	ProjectConfig *config.ProjectConfig
}

func NewPlan(o *Opts) Plan {
	confluenceAPI := confluence.NewConfluenceAPI(&confluence.Opts{
		ConfluenceHost: o.Config.Bongkoes.ConfluenceHost,
		Email:          o.Config.Bongkoes.AtlassianEmail,
		Token:          o.Config.Bongkoes.AtlassianToken,
	})
	bitbucketAPI := bitbucket.NewBitbucketAPI(&bitbucket.Opts{
		BitbucketWorkspace:   o.Config.Bongkoes.BitbucketWorkspace,
		BitbucketUsername:    o.Config.Bongkoes.BitbucketUsername,
		BitbucketAppPassword: o.Config.Bongkoes.BitbucketAppPassword,
		MainBranch:           o.ProjectConfig.MainBranch,
	})
	return &deploymentPlan{
		cfg:           o.Config,
		confluenceAPI: confluenceAPI,
		bitbucketAPI:  bitbucketAPI,
		git:           git.NewGitLocal(),
	}
}
