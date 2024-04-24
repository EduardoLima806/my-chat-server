package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
	ErrNotFound            = errors.New("NOT_FOUND")
	ErrBadRequest          = errors.New("BAD_REQUEST")
	ErrConflict            = errors.New("CONFLICT")
	ErrInsufficientFund    = errors.New("INSUFFICIENT_FUND")
	ErrUnauthorized        = errors.New("UNAUTHORIZED")
	errStruct              ErrorCodesStruct
)

type ErrorCodesStruct struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func CreateError(code string, message string) error {
	err := fmt.Sprintf(`{"error_code": "%s", "error_message": "%s"}`, code, message)
	return errors.New(err)
}

func GetHttpStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	errCode := extractErroCode(err)

	switch errCode {
	case ErrInternalServerError.Error():
		return http.StatusInternalServerError
	case ErrNotFound.Error():
		return http.StatusNotFound
	case ErrConflict.Error():
		return http.StatusConflict
	case ErrInsufficientFund.Error():
		return http.StatusBadRequest
	case ErrUnauthorized.Error():
		return http.StatusUnauthorized
	case ErrBadRequest.Error():
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func extractErroCode(err error) string {
	s := err.Error()

	_ = json.Unmarshal([]byte(s), &errStruct)

	return strings.ToUpper(errStruct.ErrorCode)
}

func ErrorCodeResponse(err error) ErrorCodesStruct {
	s := err.Error()

	_ = json.Unmarshal([]byte(s), &errStruct)

	return errStruct
}
