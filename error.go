package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorRes struct {
	Message string `json:"error"`
}

func marshalError(err error, message string, status int, w http.ResponseWriter) {
	w.WriteHeader(status)

	response := errorRes{}
	response.Message = fmt.Sprintf("%s: %v", message, err)

	bytes, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte("marshalling error"))
	} else {
		w.Write(bytes)
	}
}
