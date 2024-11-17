package types

type VersionStatus string

func (v VersionStatus) String() string {
	return string(v)
}

const (
	VersionReleased   VersionStatus = "released"
	VersionUnReleased VersionStatus = "unreleased"
	VersionArchived   VersionStatus = "archived"
)

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
	Released    bool   `json:"released"`
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

	WebLink string `json:"-"`
}

type ProjectDetailResponse struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type FetchLatestVersionRequest struct {
	ProjectKey string
	Query      string
	Status     VersionStatus
}

type VersionResponse struct {
	Values []Version `json:"values"`
}

type Version struct {
	Self            string `json:"self"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	Archived        bool   `json:"archived"`
	Released        bool   `json:"released"`
	StartDate       string `json:"startDate"`
	ReleaseDate     string `json:"releaseDate"`
	UserStartDate   string `json:"userStartDate"`
	UserReleaseDate string `json:"userReleaseDate"`
	ProjectID       int64  `json:"projectId"`
}

type UpdateVersion struct {
	Update FixVersionArgs `json:"update"`
}

type FixVersionArgs struct {
	FixVersions FieldUpdateOperation `json:"fixVersions"`
}

type FieldUpdateOperation struct {
	Add    map[string]interface{} `json:"add,omitempty"`
	Copy   map[string]interface{} `json:"copy,omitempty"`
	Edit   map[string]interface{} `json:"edit,omitempty"`
	Remove map[string]interface{} `json:"remove,omitempty"`
	Set    map[string]interface{} `json:"set,omitempty"`
}

type User struct {
	Type        string `json:"type"`
	AccountID   string `json:"accountId"`
	AccountType string `json:"accountType"`
	Email       string `json:"email"`
	PublicName  string `json:"publicName"`
}
