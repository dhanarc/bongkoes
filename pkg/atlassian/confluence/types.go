package confluence

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
