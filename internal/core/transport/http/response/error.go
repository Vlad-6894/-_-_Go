package core_http_response

type ErrorResponse struct {
	Error   string `json:"details"         example:"full error text"`
	Message string `json:"message"       example:"message about error"`
}
