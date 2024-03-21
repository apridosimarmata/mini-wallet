package response

import (
	"encoding/json"
	"net/http"
)

const (
	STATUS_SUCCESS = "success"
	STATUS_FAIL    = "fail"
	STATUS_ERROR   = "error"

	WALLET_DISABLED_ERROR = "Wallet disabled"
)

type Error struct {
	Error string `json:"error"`
}

type Response[T any] struct {
	Status     string `json:"status"`
	Data       *T     `json:"data,omitempty"`
	StatusCode int    `json:"-"`
}

func (payload *Response[T]) Success(msg string, data T) {
	payload.Status = STATUS_SUCCESS
	payload.StatusCode = http.StatusOK
}

func (payload *Response[T]) Fail() {
	payload.Status = STATUS_FAIL
	payload.StatusCode = http.StatusBadRequest
}

func (payload *Response[T]) Error() {
	payload.Status = STATUS_ERROR
	payload.StatusCode = http.StatusInternalServerError
}

func (res *Response[T]) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res)
}
