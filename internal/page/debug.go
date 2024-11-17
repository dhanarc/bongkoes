package page

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/djk-lgtm/bongkoes/pkg/bitbucket"
)

func (d *deploymentPlan) Debug(ctx context.Context) error {
	res, err := d.bitbucketAPI.GetLatestPullRequests(ctx, "mekari-payment", bitbucket.PaginationRequest{
		Page: 1,
		Size: 10,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	b, _ := json.Marshal(res)
	fmt.Println(string(b))
	return nil

}
