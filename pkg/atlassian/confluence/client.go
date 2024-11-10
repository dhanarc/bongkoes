package confluence

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/djk-lgtm/atlassianoto/pkg/httpreq"
	"net/http"
)

const (
	PagesAPI = "/wiki/api/v2/pages"
)

type API interface {
	GetPageByID(context.Context, string) (*Page, error)
	CreatePage(context.Context, *CreatePageRequest) (*Page, error)
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

	respEntity := Page{}
	err = json.Unmarshal(response, &respEntity)
	if err != nil {
		return nil, err
	}

	return &respEntity, nil
}

func (c *confluenceAPI) CreatePage(ctx context.Context, request *CreatePageRequest) (*Page, error) {
	requestBytes, _ := json.Marshal(request)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodPost, PagesAPI, headers, requestBytes)
	if err != nil {
		return nil, err
	}

	respEntity := Page{}
	err = json.Unmarshal(response, &respEntity)
	if err != nil {
		return nil, err
	}

	return &respEntity, nil
}
