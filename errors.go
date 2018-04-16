package king

import (
	"net/http"
)

const (
	ErrorTypeNotFound            = "NOT_FOUND"
	ErrorTypeUnauthorized        = "UNAUTHORIZED"
	ErrorTypeNotImplemented      = "NOT_IMPLEMENTED"
	ErrorTypeBadRequest          = "BAD_REQUEST"
	ErrorTypeForbidden           = "STATUS_FORBIDDEN"
	ErrorTypeInternalServerError = "STATUS_INTERNAL_SERVER_ERROR"
)

var kingErrStatus = map[string]int{
	ErrorTypeNotFound:            http.StatusNotFound,
	ErrorTypeUnauthorized:        http.StatusUnauthorized,
	ErrorTypeNotImplemented:      http.StatusNotImplemented,
	ErrorTypeBadRequest:          http.StatusBadRequest,
	ErrorTypeForbidden:           http.StatusForbidden,
	ErrorTypeInternalServerError: http.StatusInternalServerError,
}

var statusKingErr = map[int]string{
	http.StatusNotFound:            ErrorTypeNotFound,
	http.StatusUnauthorized:        ErrorTypeUnauthorized,
	http.StatusNotImplemented:      ErrorTypeNotImplemented,
	http.StatusBadRequest:          ErrorTypeBadRequest,
	http.StatusForbidden:           ErrorTypeForbidden,
	http.StatusInternalServerError: ErrorTypeInternalServerError,
}

type KingError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
