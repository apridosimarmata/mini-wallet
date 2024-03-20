package response

import (
	"encoding/json"
	"net/http"
)

const (
	STATUS_SUCCESS               = "success"
	STATUS_FAIL                  = "fail"
	STATUS_NOT_FOUND             = "not found"
	STATUS_INTERNAL_SERVER_ERROR = "internal server error"
)

type Error struct {
	Error string `json:"error"`
}

type Response[T any] struct {
	Status     string `json:"status"`
	Data       T      `json:"data,omitempty"`
	StatusCode int    `json:"-"`
}

func (payload *Response[T]) Success(msg string, data T) Response[T] {
	return Response[T]{
		Status:     STATUS_SUCCESS,
		Data:       data,
		StatusCode: http.StatusOK,
	}
}

func (payload *Response[T]) BadRequest(msg string) Response[T] {
	return Response[T]{
		Status:     STATUS_FAIL,
		StatusCode: http.StatusBadRequest,
	}
}

func (payload *Response[T]) NotFound(msg string) Response[T] {
	return Response[T]{
		Status:     STATUS_NOT_FOUND,
		StatusCode: http.StatusBadRequest,
	}
}

func (payload *Response[T]) InternalServerError(msg string) Response[T] {
	return Response[T]{
		Status:     STATUS_INTERNAL_SERVER_ERROR,
		StatusCode: http.StatusBadRequest,
	}
}

func (res *Response[T]) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res)
}
