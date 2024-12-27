package errors

import (
	"errors"
	"fmt"
)

type ErrorKind int

const (
	KindDatabase ErrorKind = iota
	KindNotFound
	//KindValidation
	//KindUnauthorized
)

type ServerError struct {
	Kind    ErrorKind
	Message string
	Err     error
}

// *ServerError implements error interface
func (e *ServerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Is 自定义 通过 errors.Is(err, target error) 比较错误时的比较逻辑
func (e *ServerError) Is(target error) bool {
	// 将 target (interface) 转换为 *ServerError 类型, 并将结果存储在 t 中
	var t *ServerError
	ok := errors.As(target, &t)

	if !ok {
		return false
	}
	return e.Kind == t.Kind
}

func NewDatabaseError(msg string, err error) *ServerError {
	return &ServerError{
		Kind:    KindDatabase,
		Message: msg,
		Err:     err,
	}
}

func NewNotFoundError(msg string, err error) *ServerError {
	return &ServerError{
		Kind:    KindNotFound,
		Message: msg,
		Err:     err,
	}
}
