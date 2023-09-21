package rest

import (
	"T/app/core"
	"encoding/json"
	"fmt"
	"net/http"
)

const critical = `{"message": "Server error", "code": "server_error"}`

func fmtError(w http.ResponseWriter, m string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	_, _ = fmt.Fprintln(w, m)
}

func ServerError(w http.ResponseWriter, errMes error) {
	b, err := json.Marshal(core.StatusResponse{
		Status: core.Status{
			Code:    core.ServerError,
			Message: errMes.Error(),
		},
	})
	if err != nil {
		fmtError(w, critical)
		return
	}

	fmtError(w, string(b))
}

func ValidationError(w http.ResponseWriter, msg string) {
	b, err := json.Marshal(core.StatusResponse{
		Status: core.Status{
			Code:    core.ValidationError,
			Message: msg,
		},
	})
	if err != nil {
		fmtError(w, critical)
		return
	}

	fmtError(w, string(b))
}
