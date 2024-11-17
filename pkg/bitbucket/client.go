package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/djk-lgtm/bongkoes/pkg/httpreq"
	"github.com/samber/lo"
	"net/http"
	"net/url"
)

const (
	BaseURL  = "https://api.bitbucket.org"
	Version2 = "2.0"
)

type API interface {
	GetTagsByDateDesc(context.Context, string) (*RefsTagsResponse, error)
	RunPipelineBranch(context.Context, string, string, string) (*string, error)
	GetLatestPullRequests(context.Context, string, PaginationRequest) (*PullRequestResponse, error)
}

type bitbucketAPI struct {
	httpClient *httpreq.HTTPClient
	workspace  string
	mainBranch string
}

type Opts struct {
	BitbucketWorkspace   string
	BitbucketUsername    string
	BitbucketAppPassword string
	MainBranch           string
}

func NewBitbucketAPI(o *Opts) API {
	return &bitbucketAPI{
		httpClient: httpreq.NewHTTPClient(&httpreq.Opts{
			Endpoint: BaseURL,
			Username: o.BitbucketUsername,
			Password: o.BitbucketAppPassword,
		}),
		workspace:  o.BitbucketWorkspace,
		mainBranch: o.MainBranch,
	}
}

func (b *bitbucketAPI) RunPipelineBranch(ctx context.Context, repository, branch, pipeline string) (*string, error) {
	path := fmt.Sprintf("/%s/repositories/%s/%s/pipelines", Version2, b.workspace, repository)

	triggerRequest := &TriggerPipelineTargetRequest{
		Target: TargetPipeline{
			Type:    "pipeline_ref_target",
			RefType: PipelineBranch.String(),
			RefName: branch,
			Selector: Selector{
				Type:    "custom",
				Pattern: pipeline,
			},
		},
	}

	requestBytes, _ := json.Marshal(triggerRequest)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	response, err := b.httpClient.ExecuteBasicAuth(ctx, http.MethodPost, path, headers, requestBytes)
	if err != nil {
		return nil, err
	}

	responseBody := new(RunPipelineResponse)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	return lo.ToPtr(fmt.Sprintf("https://bitbucket.org/%s/%s/pipelines/results/%d", b.workspace, repository, responseBody.BuildNumber)), nil
}

func (b *bitbucketAPI) GetTagsByDateDesc(ctx context.Context, repository string) (*RefsTagsResponse, error) {
	const fetchedFields = "values.name,values.type,values.target.date"
	const sortField = "-target.date"
	queryParams := url.Values{}
	queryParams.Set("fields", fetchedFields)
	queryParams.Set("sort", sortField)
	path := fmt.Sprintf("/%s/repositories/%s/%s/refs/tags?%s", Version2, b.workspace, repository, queryParams.Encode())

	response, err := b.httpClient.ExecuteBasicAuth(ctx, http.MethodGet, path, make(map[string]string), nil)
	if err != nil {
		return nil, err
	}

	responseBody := new(RefsTagsResponse)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
