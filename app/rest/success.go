package rest

import (
	"T/app/core"
	"encoding/json"
	"fmt"
	"net/http"
)

func fmtSuccess(w http.ResponseWriter, m string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	_, _ = fmt.Fprintln(w, m)
}

func ServerSuccessOK(w http.ResponseWriter) {
	b, err := json.Marshal(core.StatusResponse{
		Status: core.Status{
			Code:    core.Success,
			Message: "Success!",
			Data:    nil,
		},
	})
	if err != nil {
		fmtError(w, critical)
		return
	}

	fmtSuccess(w, string(b))
}

func ServerSuccessStruct(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(core.StatusResponse{
		Status: core.Status{
			Code:    core.Success,
			Message: "Success!",
			Data:    data,
		},
	})
	if err != nil {
		fmtError(w, critical)
		return
	}

	fmtSuccess(w, string(b))
}
