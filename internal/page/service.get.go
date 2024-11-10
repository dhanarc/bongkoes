package page

import (
	"context"
	"github.com/djk-lgtm/atlassianoto/pkg/console"
)

func (d *deploymentPlan) GetConfig(ctx context.Context, serviceName string) error {
	var configs []Service
	tx := d.db.WithContext(ctx)
	if len(serviceName) > 0 {
		tx = tx.Where("service_name = ?", serviceName).First(&configs)
	} else {
		tx = tx.Find(&configs)
	}

	if tx.Error != nil {
		return tx.Error
	}

	console.PrintTable(configs)
	return nil
}
