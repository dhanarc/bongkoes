package bitbucket

type PipelineRefType string

func (p PipelineRefType) String() string {
	return string(p)
}

const PipelineBranch PipelineRefType = "branch"

type RefsTagsResponse struct {
	Values []RefsTag `json:"values"`
}

type RefsTag struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Target struct {
		Date string `json:"date"`
	} `json:"target"`
}

type Link struct {
	Base  string `json:"base"`
	WebUI string `json:"webui"`
}

type Page struct {
	ParentType string   `json:"parentType"`
	OwnerID    string   `json:"ownerId"`
	Title      string   `json:"title"`
	Body       BodyPage `json:"body"`
	ParentID   string   `json:"parentId"`
	SpaceID    string   `json:"spaceId"`
	Links      Link     `json:"_links"`
}

type CreatePageRequest struct {
	SpaceID  string   `json:"spaceId"`
	Status   string   `json:"status"`
	Title    string   `json:"title"`
	ParentID string   `json:"parentId"`
	Body     BodyPage `json:"body"`
}

type BodyPage struct {
	Storage BodyStorage `json:"storage"`
}

type BodyStorage struct {
	Representation string `json:"representation"`
	Value          string `json:"value"`
}

type CreateVersionRequest struct {
	Archived    bool   `json:"archived"`
	Description string `json:"description"`
	Name        string `json:"name"`
	ProjectID   uint64 `json:"projectId"`
	ReleaseDate string `json:"releaseDate"`
	StartDate   string `json:"startDate"`
	Released    string `json:"released"`
}

type CreateVersionResponse struct {
	Archived        bool   `json:"archived"`
	Description     string `json:"description"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	ProjectID       int64  `json:"projectId"`
	ReleaseDate     string `json:"releaseDate"`
	Released        bool   `json:"released"`
	Self            string `json:"self"`
	UserReleaseDate string `json:"userReleaseDate"`
}

type ProjectDetailResponse struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type TriggerPipelineTargetRequest struct {
	Target TargetPipeline `json:"target"`
}

type TargetPipeline struct {
	Type    string `json:"type"`
	RefType string `json:"ref_type"`
	RefName string `json:"ref_name"`
}

type RunPipelineResponse struct {
	UUID  string `json:"uuid"`
	Links struct {
		Type string `json:"type"`
	} `json:"links"`
}
