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

// Is 实现 Is() 方法是为了 自定义 通过 errors.Is(err, target error) 比较时的逻辑
// *ServerError 也实现了 error 接口
// 可是为什么这里的 Is() 方法的参数是 error 类型, 而不能是 *ServerError 类型呢?
func (e *ServerError) Is(target error) bool {
	// target 是 error,
	// 尝试将 target 转换为 *ServerError 类型, 并将结果存储在 t 中
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
