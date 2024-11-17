package view

import (
	"encoding/json"
	"fmt"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence/types"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"html"
	"net/url"
)

func GetLinkCard(link string) string {
	return fmt.Sprintf("<a href=\"%s\" data-card-appearance=\"inline\">%s<a>", link, link)
}

func GetJiraIssuesLink(host, jql string, viewOptions []types.JiraLinkView) (*string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("error parsing url:%+v", err)
	}

	// Set path
	u.Path = "issues"

	// Set query parameters
	q := u.Query()
	q.Set("jql", jql)
	u.RawQuery = q.Encode()
	jiraRawLink := u.String()

	// data source
	dataSource := constructJiraLinkDataSource(jql, viewOptions)
	ds, err := json.Marshal(dataSource)
	if err != nil {
		return nil, fmt.Errorf("error marshal datasource:%+v", err)
	}
	dataSourceString := string(ds)
	encodedDS := html.EscapeString(dataSourceString)
	link := fmt.Sprintf("<a href=\"%s\" data-layout=\"wide\" data-card-appearance=\"block\" data-datasource=\"%s\">%s</a>", jiraRawLink, encodedDS, jiraRawLink)

	return lo.ToPtr(link), nil
}

func constructJiraLinkDataSource(jql string, viewsOptions []types.JiraLinkView) *types.JiraLinkDataSource {
	ID := uuid.NewString()
	cloudID := uuid.NewString()
	return &types.JiraLinkDataSource{
		ID: ID,
		Parameters: types.JiraLinkParameters{
			JQL:     jql,
			CloudID: cloudID,
		},
		Views: viewsOptions,
	}
}
