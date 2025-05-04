package errors

import (
	"errors"
	"net/http"
)

// HTTPStatus 根据 ServerError.Kind 返回对应的 HTTP Status Code
// 如果不是 *ServerError，就返回 500
func HTTPStatus(err error) int {
	var se *ServerError
	if errors.As(err, &se) {
		switch se.Kind {
		case KindNotFound:
			return http.StatusNotFound // 404
		case KindAlreadyExists, KindConflict:
			return http.StatusConflict // 409
		case KindValidation:
			return http.StatusBadRequest // 400
		case KindUnauthorized:
			return http.StatusUnauthorized // 401
		case KindForbidden:
			return http.StatusForbidden // 403
		case KindInternal:
			fallthrough
		default:
			return http.StatusInternalServerError // 500
		}
	}
	return http.StatusInternalServerError
}
