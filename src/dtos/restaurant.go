package dtos

type FindRestaurantsRequest struct {
	Place  string `form:"place" json:"place" binding:"required"`
	Cursor string `form:"cursor" json:"cursor"`
}

type FindRestaurantsResponse struct {
	Meta Meta                   `json:"meta"`
	Data []*FindRestaurantsData `json:"data"`
}

type FindRestaurantsData struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
