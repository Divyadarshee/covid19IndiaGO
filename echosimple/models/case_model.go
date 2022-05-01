package models

var GPS struct {
	Latitude  float64 `json:"latitude, omitempty" validate:"required"`
	Longitude float64 `json:"longitude, omitempty" validate:"required"`
}

type Cases struct {
	StateCode   string
	StateName   string
	Confirmed   int64  `json:"confirmed"`
	Deceased    int64  `json:"deceased"`
	Recovered   int64  `json:"recovered"`
	Tested      int64  `json:"tested"`
	Vaccinated1 int64  `json:"vaccinated1"`
	Vaccinated2 int64  `json:"vaccinated2"`
	LastUpdated string `json:"last_updated"`
}
