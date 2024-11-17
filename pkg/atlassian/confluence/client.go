package confluence

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence/types"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence/view"
	"github.com/djk-lgtm/bongkoes/pkg/httpreq"
	"net/http"
	"net/url"
)

const (
	PagesAPI       = "/wiki/api/v2/pages"
	CurrentUserAPI = "/wiki/rest/api/user/current"

	VersionAPI = "/rest/api/2/version"
	ProjectAPI = "/rest/api/2/project"
	IssueAPI   = "/rest/api/2/issue"
)

type API interface {
	GetPageByID(context.Context, string) (*types.Page, error)
	CreatePage(context.Context, *types.CreatePageRequest) (*types.Page, error)

	GetProjectDetail(context.Context, string) (*types.ProjectDetailResponse, error)

	CreateVersion(context.Context, *types.CreateVersionRequest) (*types.CreateVersionResponse, error)
	GetLatestVersion(context.Context, *types.FetchLatestVersionRequest) (*types.Version, error)

	AddIssueFixVersion(context.Context, string, string) error

	GetCurrentUserMention(context.Context) (*string, error)
	GenerateJiraLink(jql string, viewOptions []types.JiraLinkView) (*string, error)
}

type confluenceAPI struct {
	httpClient *httpreq.HTTPClient
	host       string
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
		host: o.ConfluenceHost,
	}
}

func (c *confluenceAPI) GetPageByID(ctx context.Context, pageID string) (*types.Page, error) {
	path := fmt.Sprintf("%s/%s?body-format=storage", PagesAPI, pageID)
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodGet, path, make(map[string]string), nil)
	if err != nil {
		return nil, err
	}

	responseBody := new(types.Page)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *confluenceAPI) CreatePage(ctx context.Context, request *types.CreatePageRequest) (*types.Page, error) {
	requestBytes, _ := json.Marshal(request)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodPost, PagesAPI, headers, requestBytes)
	if err != nil {
		return nil, err
	}

	responseBody := new(types.Page)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *confluenceAPI) GetProjectDetail(ctx context.Context, projectIDOrKey string) (*types.ProjectDetailResponse, error) {
	path := fmt.Sprintf("%s/%s", ProjectAPI, projectIDOrKey)
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodGet, path, make(map[string]string), nil)
	if err != nil {
		return nil, err
	}

	responseBody := new(types.ProjectDetailResponse)
	err = json.Unmarshal(response, &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *confluenceAPI) CreateVersion(ctx context.Context, request *types.CreateVersionRequest) (*types.CreateVersionResponse, error) {
	requestBytes, _ := json.Marshal(request)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodPost, VersionAPI, headers, requestBytes)
	if err != nil {
		return nil, err
	}

	responseBody := new(types.CreateVersionResponse)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *confluenceAPI) GetLatestVersion(ctx context.Context, request *types.FetchLatestVersionRequest) (*types.Version, error) {
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

	responseBody := new(types.VersionResponse)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	if len(responseBody.Values) == 0 {
		return nil, nil
	}

	return &responseBody.Values[0], nil
}

func (c *confluenceAPI) AddIssueFixVersion(ctx context.Context, issueKey string, versionID string) error {
	path := fmt.Sprintf("%s%s", IssueAPI, issueKey)

	updateOperation := new(types.FieldUpdateOperation)
	updateOperation.Add = map[string]interface{}{
		"id": versionID,
	}

	request := new(types.UpdateVersion)
	request.Update = types.FixVersionArgs{
		FixVersions: *updateOperation,
	}
	requestBytes, _ := json.Marshal(request)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodPut, path, headers, requestBytes)
	if err != nil {
		return err
	}

	responseBody := new(types.CreateVersionResponse)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *confluenceAPI) GenerateJiraLink(jql string, viewOptions []types.JiraLinkView) (*string, error) {
	return view.GetJiraIssuesLink(c.host, jql, viewOptions)
}
