package loki

type SendDataLogInput struct {
	Request SendDataLogRequest
}

type SendDataLogRequest struct {
	Streams []SendDataLogStramRequest `json:"streams"`
}

type SendDataLogStramRequest struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}
