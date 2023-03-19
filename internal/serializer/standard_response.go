package serializer

import "encoding/json"

type Response struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Errors []ErrorInfo `json:"errors,omitempty"`
}

type ErrorInfo struct {
	Detail string `json:"detail,omitempty"`
}

func (resp *Response) ToJSON() []byte {
	jsonBytes, _ := json.Marshal(resp)
	return jsonBytes
}
