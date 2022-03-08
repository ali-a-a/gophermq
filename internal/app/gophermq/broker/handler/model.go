package handler

// PublishReq represents publish request body.
type PublishReq struct {
	Subject string `json:"subject"`
	Data    string `json:"data"`
}
