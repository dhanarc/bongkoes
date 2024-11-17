package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (b *bitbucketAPI) GetLatestPullRequests(ctx context.Context, repository string, pageRequest PaginationRequest) (*PullRequestResponse, error) {
	const fetchedFields = "values.type,values.id,values.title,values.state,values.merge_commit.hash,values.merge_commit.type,values.created_on,values.updated_on"
	const sortField = "-updated_on"
	query := fmt.Sprintf("destination.branch.name=\"%s\"+AND+state=\"MERGED\"", b.mainBranch)
	queryParams := url.Values{}
	queryParams.Set("fields", fetchedFields)
	queryParams.Set("sort", sortField)
	path := fmt.Sprintf("/%s/repositories/%s/%s/pullrequests?%s&q=%s", Version2, b.workspace, repository, queryParams.Encode(), query)

	response, err := b.httpClient.ExecuteBasicAuth(ctx, http.MethodGet, path, make(map[string]string), nil)
	if err != nil {
		return nil, err
	}

	responseBody := new(PullRequestResponse)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
