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

type Selector struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}
type TargetPipeline struct {
	Selector Selector `json:"selector"`
	Type     string   `json:"type"`
	RefType  string   `json:"ref_type"`
	RefName  string   `json:"ref_name"`
}

type RunPipelineResponse struct {
	UUID        string `json:"uuid"`
	BuildNumber uint64 `json:"build_number"`
	Links       struct {
		Type string `json:"type"`
	} `json:"links"`
}

type PaginationRequest struct {
	Page uint
	Size uint
}

type PullRequestResponse struct {
	Values  []PullRequest `json:"values"`
	PageLen uint          `json:"pagelen"`
	Size    uint64        `json:"size"`
	Page    uint64        `json:"page"`
}

type PullRequest struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Status      string `json:"status"`
	MergeCommit struct {
		Hash string `json:"hash"`
		Type string `json:"type"`
	} `json:"merge_commit"`
	Reason    string `json:"reason"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}
