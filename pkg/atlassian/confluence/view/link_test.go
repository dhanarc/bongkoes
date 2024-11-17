package view_test

import (
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence/types"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence/view"
	"testing"
)

func Setup() confluence.API {
	confluenceAPI := confluence.NewConfluenceAPI(&confluence.Opts{
		ConfluenceHost: "https://ur.atlassian.net",
		Email:          "your_name@mail.com",
		Token:          "<TOKEN>",
	})
	return confluenceAPI
}

func TestGetJiraIssueLink(t *testing.T) {
	JQL := "project = KRED AND fixversion = \"Mekari Payment - v1.193.0\" ORDER BY fixVersion, Rank ASC"
	viewOptions := []types.JiraLinkView{
		{
			Type: "table",
			Properties: types.JiraViewProperty{
				Columns: []types.JiraPropertyKey{
					{
						Key:   "key",
						Width: 91,
					},
					{
						Key:       "summary",
						IsWrapped: true,
					},
					{
						Key:       "assignee",
						Width:     148,
						IsWrapped: true,
					},
					{
						Key:       "customfield_11620",
						Width:     304,
						IsWrapped: true,
					},
					{
						Key: "status",
					},
				},
			},
		},
	}

	jiraLink, err := view.GetJiraIssuesLink("https://yourproject.atlassian.net", JQL, viewOptions)
	if err != nil {
		t.Errorf("[GetJiraIssuesLink] err should be empty")
	}
	if jiraLink == nil {
		t.Errorf("[GetJiraIssuesLink] jiraLink should not be nil")
	}
}
