package page

import "time"

type Service struct {
	ID                 uint      `gorm:"primarykey" header:"ID"`
	TribeName          string    `gorm:"column:tribe_name" header:"Tribe Name"`
	TeamName           string    `gorm:"column:team_name" header:"Team Name"`
	ServiceCode        string    `gorm:"column:service_code" header:"Service Code"`
	ServiceName        string    `gorm:"column:service_name" header:"Service Name"`
	TemplateID         string    `gorm:"column:template_id" header:"Template ID"`
	DeploymentFolderID string    `gorm:"column:deployment_folder_id" header:"Deployment Folder ID"`
	CreatedAt          time.Time `header:"Created At"`
	UpdatedAt          time.Time `header:"Updated At"`
}

type CreateServiceArgs struct {
	TeamName           string
	TribeName          string
	ServiceCode        string
	ServiceName        string
	TemplateID         string
	DeploymentFolderID string
}

type CreateDeploymentArgs struct {
	ServiceCode    string
	Tag            string
	DeploymentTime string
	DownTimeEst    string
	RollbackTag    string
	Published      bool
}

type DeploymentArgs struct {
	ServiceCode    string
	ServiceName    string
	TeamName       string
	TribeName      string
	Tag            string
	DeploymentTime string
	DownTimeEst    string
	RollbackTag    string
	JiraLink       string
}

type RenderArgs struct {
	Service                           Service
	CurrentTime                       time.Time
	Tag, RollbackTag                  string
	DeploymentTime, EstimatedDownTime string
	PageTemplate                      string
}
