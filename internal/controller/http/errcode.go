package http

import (
	"errors"
)

type ErrCode interface {
	ErrCode() string
}

func errCode(err error) string {
	var errWithCode ErrCode
	if errors.As(err, &errWithCode) {
		return errWithCode.ErrCode()
	}

	switch {
	case errors.Is(err, ErrJsonMarshalResponseData):
		return ErrCodeUnmarshalJSONResponseData
	}

	return ErrCodeUnknown
}

const (
	ErrCodeUnknown                   = ""
	ErrCodeUnmarshalJSONResponseData = "UNMARSHAL_JSON_RESPONSE_DATA"
)
