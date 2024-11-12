package confluence

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/djk-lgtm/bongkoes/pkg/httpreq"
	"net/http"
	"net/url"
)

const (
	PagesAPI = "/wiki/api/v2/pages"

	VersionAPI = "/rest/api/2/version"
	ProjectAPI = "/rest/api/2/project"
)

type API interface {
	GetPageByID(context.Context, string) (*Page, error)
	CreatePage(context.Context, *CreatePageRequest) (*Page, error)

	GetProjectDetail(context.Context, string) (*ProjectDetailResponse, error)

	CreateVersion(context.Context, *CreateVersionRequest) (*CreateVersionResponse, error)
	GetLatestVersion(context.Context, *FetchLatestVersionRequest) (*Version, error)
}

type confluenceAPI struct {
	httpClient *httpreq.HTTPClient
}

type Opts struct {
	ConfluenceHost string
	Email          string
	Token          string
}

func NewConfluenceAPI(o *Opts) API {
	return &confluenceAPI{
		httpClient: httpreq.NewHTTPClient(&httpreq.Opts{
			Endpoint: o.ConfluenceHost,
			Username: o.Email,
			Password: o.Token,
		}),
	}
}

func (c *confluenceAPI) GetPageByID(ctx context.Context, pageID string) (*Page, error) {
	path := fmt.Sprintf("%s/%s?body-format=storage", PagesAPI, pageID)
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodGet, path, make(map[string]string), nil)
	if err != nil {
		return nil, err
	}

	responseBody := new(Page)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *confluenceAPI) CreatePage(ctx context.Context, request *CreatePageRequest) (*Page, error) {
	requestBytes, _ := json.Marshal(request)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodPost, PagesAPI, headers, requestBytes)
	if err != nil {
		return nil, err
	}

	responseBody := new(Page)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *confluenceAPI) GetProjectDetail(ctx context.Context, projectIDOrKey string) (*ProjectDetailResponse, error) {
	path := fmt.Sprintf("%s/%s", ProjectAPI, projectIDOrKey)
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodGet, path, make(map[string]string), nil)
	if err != nil {
		return nil, err
	}

	responseBody := new(ProjectDetailResponse)
	err = json.Unmarshal(response, &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *confluenceAPI) CreateVersion(ctx context.Context, request *CreateVersionRequest) (*CreateVersionResponse, error) {
	requestBytes, _ := json.Marshal(request)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodPost, VersionAPI, headers, requestBytes)
	if err != nil {
		return nil, err
	}

	responseBody := new(CreateVersionResponse)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *confluenceAPI) GetLatestVersion(ctx context.Context, request *FetchLatestVersionRequest) (*Version, error) {
	queryParams := url.Values{}
	queryParams.Set("orderBy", "-releaseDate")
	queryParams.Set("query", request.Query)
	queryParams.Set("maxResults", "1")
	queryParams.Set("status", request.Status.String())
	path := fmt.Sprintf("%s/%s/version?%s", ProjectAPI, request.ProjectKey, queryParams.Encode())
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodGet, path, make(map[string]string), nil)
	if err != nil {
		return nil, err
	}

	responseBody := new(VersionResponse)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	if len(responseBody.Values) == 0 {
		return nil, nil
	}

	return &responseBody.Values[0], nil
}
