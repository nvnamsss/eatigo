package dtos

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Cursor  string `json:"cursor"`
}
