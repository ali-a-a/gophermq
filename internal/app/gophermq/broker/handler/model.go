package handler

// PublishReq represents publish request body.
type PublishReq struct {
	Subject string `json:"subject"`
	Data    string `json:"data"`
}

// SubscribeReq represents subscribe request body.
type SubscribeReq struct {
	Subject string `json:"subject"`
}

// SubscribeRes represents subscribe response body.
type SubscribeRes struct {
	Subject string `json:"subject"`
	ID      string `json:"id"`
}

// FetchReq represents fetch request body.
type FetchReq struct {
	Subject string `json:"subject"`
	ID      string `json:"id"`
}

// FetchRes represents fetch response body.
type FetchRes struct {
	Subject string   `json:"subject"`
	ID      string   `json:"id"`
	Data    []string `json:"data"`
}
