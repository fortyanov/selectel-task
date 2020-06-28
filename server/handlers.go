package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"selectel-task/vscale"
)

type Scalets struct {
	Count int `json:"count"`
}

func createScalets(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		scalets = new(Scalets)
	)

	if err = json.NewDecoder(r.Body).Decode(&scalets); err != nil {
		JsonError(w, err, http.StatusBadRequest)
		return
	}
	if scalets.Count < 1 {
		JsonError(w, errors.New(errScaletsCount), http.StatusBadRequest)
		return
	}

	if err = vscale.NewClient().CreateScalets(scalets.Count); err != nil {
		JsonError(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	JsonCreated(w)
}

func deleteScalets(w http.ResponseWriter, r *http.Request) {
	vscale.NewClient().DeleteAllScalets()
	JsonDeleted(w)
}
