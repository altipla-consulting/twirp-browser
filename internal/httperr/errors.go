package httperr

import (
	"net/http"
)

const (
	ErrorTypeNotFound            = "NOT_FOUND"
	ErrorTypeUnauthorized        = "UNAUTHORIZED"
	ErrorTypeNotImplemented      = "NOT_IMPLEMENTED"
	ErrorTypeNotValid            = "NOT_VALID"
	ErrorTypeForbidden           = "STATUS_FORBIDDEN"
	ErrorTypeInternalServerError = "STATUS_INTERNAL_SERVER_ERROR"
)

var KingErrStatus = map[string]int{
	ErrorTypeNotFound:            http.StatusNotFound,
	ErrorTypeUnauthorized:        http.StatusUnauthorized,
	ErrorTypeNotImplemented:      http.StatusNotImplemented,
	ErrorTypeNotValid:            http.StatusBadRequest,
	ErrorTypeForbidden:           http.StatusForbidden,
	ErrorTypeInternalServerError: http.StatusInternalServerError,
}

var StatusKingErr = map[int]string{
	http.StatusNotFound:            ErrorTypeNotFound,
	http.StatusUnauthorized:        ErrorTypeUnauthorized,
	http.StatusNotImplemented:      ErrorTypeNotImplemented,
	http.StatusBadRequest:          ErrorTypeNotValid,
	http.StatusForbidden:           ErrorTypeForbidden,
	http.StatusInternalServerError: ErrorTypeInternalServerError,
}
