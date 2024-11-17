package confluence

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/djk-lgtm/bongkoes/pkg/atlassian/confluence/types"
	"github.com/samber/lo"
	"net/http"
)

func (c *confluenceAPI) GetCurrentUserMention(ctx context.Context) (*string, error) {
	response, err := c.httpClient.ExecuteBasicAuth(ctx, http.MethodGet, CurrentUserAPI, make(map[string]string), nil)
	if err != nil {
		return nil, err
	}

	responseBody := new(types.User)
	err = json.Unmarshal(response, responseBody)
	if err != nil {
		return nil, err
	}

	mentionElement := fmt.Sprintf("<ac:link><ri:user ri:account-id=\"%s\"/></ac:link>", responseBody.AccountID)
	return lo.ToPtr(mentionElement), nil
}
