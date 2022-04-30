package models

var GPS struct {
	Latitude  float64 `json:"latitude, omitempty" validate:"required"`
	Longitude float64 `json:"longitude, omitempty" validate:"required"`
}
