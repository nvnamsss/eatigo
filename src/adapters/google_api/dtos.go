package google_api

type Meta struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	NextPageToken string `json:"next_page_token"`
}

type FindRestaurantsRequest struct {
	Place         string `json:"place"`
	Radius        int    `json:"radius"`
	NextPageToken string `json:"next_page_token"`
}

type findRestaurantsDataAdapt struct {
	NextPageToken string                `json:"next_page_token"`
	Results       []FindRestaurantsData `json:"results"`
	Status        string                `json:"status"`
}

type FindRestaurantsResponse struct {
	Meta Meta                  `json:"meta"`
	Data []FindRestaurantsData `json:"data"`
}

type FindRestaurantsData struct {
	Name                string   `json:"name"`
	Address             string   `json:"formatted_address"`
	Icon                string   `json:"icon"`
	IconBackgroundColor string   `json:"icon_background_color"`
	PlaceID             string   `json:"place_id"`
	PriceLevel          int      `json:"price_level"`
	Rating              float64  `json:"rating"`
	UserRatingTotal     int      `json:"user_rating_total"`
	References          string   `json:"references"`
	PermanentClosed     bool     `json:"permanent_closed"`
	Types               []string `json:"types"`
	Photos              []Photo  `json:"photos"`
	PlusCode            PlusCode `json:"plus_code"`
}

type Geometry struct {
	Location `json:"location"`
	Viewport `json:"view_port"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Viewport struct {
	NorthEast Location `json:"northeast"`
	SouthEast Location `json:"southwest"`
}

type OpeningHours struct {
	OpenNow bool `json:"open_now"`
}

type Photo struct {
	Height           int      `json:"height"`
	Width            int      `json:"width"`
	HTMLAttributions []string `json:"html_attributions"`
	PhotoReferrence  string   `json:"photo_reference"`
}

type PlusCode struct {
	CompoundCode string `json:"compound_code"`
	GlobalCode   string `json:"global_code"`
}
