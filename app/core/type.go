package core

const (
	Success         = "success"
	ServerError     = "server_error"
	ValidationError = "validation_error"
)

type Status struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type StatusResponse struct {
	Status Status `json:"status"`
}

type IdResponse struct {
	Status Status `json:"status"`
	Id     int32
}
