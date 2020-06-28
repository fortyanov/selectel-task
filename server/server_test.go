package server

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"selectel-task/vscale"
)

func init() {
	err := vscale.Init()
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateScalets(t *testing.T) {
	reqBody := []byte(`{"count": 5}`)
	req, err := http.NewRequest("POST", "/scalets", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createScalets)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestDeleteScalets(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/scalets", bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteScalets)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}
