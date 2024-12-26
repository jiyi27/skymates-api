package types

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	// omitempty tag tells the JSON encoder Data field is optional
	Data interface{} `json:"data,omitempty"`
}
