package responses

import "go-swag-sample/echosimple/models"

type CaseResponse struct {
	Status       int          `json:"status"`
	Message      string       `json:"message"`
	CasesInState models.Cases `json:"cases_in_state"`
	CasesInIndia models.Cases `json:"cases_in_india"`
}
