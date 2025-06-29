package http

import (
	"errors"
	"net/http"
)

func httpStatusCodeByErr(err error) int {
	if errors.Is(err, ErrJsonMarshalResponseData) {
		return http.StatusInternalServerError
	}
	if errors.Is(err, ErrWriteResponseBytes) {
		return http.StatusInternalServerError
	}

	return http.StatusBadRequest
}
