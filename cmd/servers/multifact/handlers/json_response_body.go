package handlers

import "encoding/json"

type ErrorResponseBody struct {
	Error string `json:"error"`
}

func (r ErrorResponseBody) ToBytes() []byte {
	json, _ := json.Marshal(r)
	return json
}

func AsError(err error) ErrorResponseBody {
	var errString string
	if err != nil {
		errString = err.Error()
	}
	return ErrorResponseBody{
		Error: errString,
	}
}
