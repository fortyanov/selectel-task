package server

import (
	"fmt"
	"net/http"
)

func JsonError(writer http.ResponseWriter, err interface{}, code int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	fmt.Fprintf(writer, `{"error":"%v"}`, err)
}

func JsonCreated(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	fmt.Fprintf(writer, `{"status":"created"}`)
}

func JsonDeleted(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(writer, `{"status":"deleted"}`)
}
