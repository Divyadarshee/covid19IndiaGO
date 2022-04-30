package responses

type CaseResponse struct {
	Status       int    `json:"status"`
	Message      string `json:"message"`
	CasesInState Cases  `json:"cases"`
	CasesInIndia Cases  `json:"cases_in_india"`
}

type Cases struct {
	Location    string
	Confirmed   int64  `json:"confirmed"`
	Deceased    int64  `json:"deceased"`
	Recovered   int64  `json:"recovered"`
	Tested      int64  `json:"tested"`
	Vaccinated1 int64  `json:"vaccinated1"`
	Vaccinated2 int64  `json:"vaccinated2"`
	LastUpdated string `json:"last_updated"`
}

//type UserResponse struct {
//	Status  int       `json:"status"`
//	Message string    `json:"message"`
//	Data    *echo.Map `json:"data"`
//}
