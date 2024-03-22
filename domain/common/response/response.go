package response

import (
	"encoding/json"
	"net/http"
)

const (
	STATUS_SUCCESS = "success"
	STATUS_FAIL    = "fail"
	STATUS_ERROR   = "error"

	ERROR_WALLET_DISABLED       = "wallet disabled"
	ERROR_WALLET_NOT_FOUND      = "wallet not found"
	ERROR_INSSUFICIENT_FUND     = "insufficient fund"
	ERROR_REFERENCE_ID_CONFLICT = "reference id already used"
	ERROR_BAD_REQUEST           = "bad request: invalid value provided"
	ERROR_UNAUTHORIZED          = "unauthorized"
)

var (
	userErrors = map[string]struct{}{
		ERROR_WALLET_DISABLED:       {},
		ERROR_WALLET_NOT_FOUND:      {},
		ERROR_INSSUFICIENT_FUND:     {},
		ERROR_REFERENCE_ID_CONFLICT: {},
		ERROR_BAD_REQUEST:           {},
	}
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

func (payload *Response[T]) Error(msg string) {
	payload.Status = STATUS_ERROR
	payload.StatusCode = http.StatusInternalServerError

	if _, isUserError := userErrors[msg]; isUserError {
		payload.Status = STATUS_FAIL
		payload.StatusCode = http.StatusBadRequest
	}
}

func (res *Response[T]) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res)
}
