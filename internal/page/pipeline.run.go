package page

import "context"

func (d *deploymentPlan) RunPipelineBranch(ctx context.Context, serviceCode, branch, pipeline string) (*string, error) {
	return d.bitbucketAPI.RunPipelineBranch(ctx, serviceCode, branch, pipeline)
}
