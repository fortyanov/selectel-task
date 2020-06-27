package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"selectel-task/vscale"
)

type Scalets struct {
	Count int `json:"count"`
}

func scaletsCreate(w http.ResponseWriter, r *http.Request) {
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

	vscale.NewVscaleClient(time.Second * time.Duration(cfg.WriteTimeout)).CreateScalets(scalets.Count)
	JsonCreated(w)
}

func scaletsDelete(w http.ResponseWriter, r *http.Request) {
	vscale.NewVscaleClient(time.Second * time.Duration(cfg.WriteTimeout)).DeleteAllScalets()
	JsonDeleted(w)
}
