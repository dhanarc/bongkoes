package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/djk-lgtm/bongkoes/pkg/httpreq"
	"net/http"
	"net/url"
)

const (
	BaseURL  = "https://api.bitbucket.org"
	Version2 = "2.0"
)

type API interface {
	GetTagsByDateDesc(context.Context, string) (*RefsTagsResponse, error)
}

type bitbucketAPI struct {
	httpClient *httpreq.HTTPClient
	workspace  string
}

type Opts struct {
	BitbucketWorkspace   string
	BitbucketUsername    string
	BitbucketAppPassword string
}

func NewBitbucketAPI(o *Opts) API {
	return &bitbucketAPI{
		httpClient: httpreq.NewHTTPClient(&httpreq.Opts{
			Endpoint: BaseURL,
			Username: o.BitbucketUsername,
			Password: o.BitbucketAppPassword,
		}),
		workspace: o.BitbucketWorkspace,
	}
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
