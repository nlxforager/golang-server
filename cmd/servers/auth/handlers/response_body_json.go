package handlers

import "encoding/json"

type ResponseBodyJSON struct {
	Error string `json:"error"`
}

func (r ResponseBodyJSON) ToBytes() []byte {
	json, _ := json.Marshal(r)
	return json
}

func AsError(err error) ResponseBodyJSON {
	var errString string
	if err != nil {
		errString = err.Error()
	}
	return ResponseBodyJSON{
		Error: errString,
	}
}
