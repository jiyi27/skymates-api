package errors

import (
	"errors"
	"fmt"
)

// ErrorKind 定义错误类型的枚举
type ErrorKind int

const (
	KindDatabase ErrorKind = iota
	KindNotFound
	KindValidation   // 取消注释，添加更多错误类型
	KindUnauthorized // 取消注释，添加更多错误类型
)

// ServerError 定义服务器错误结构
type ServerError struct {
	Kind    ErrorKind
	Message string
	Err     error
}

// Error 实现 error 接口
func (e *ServerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Is 自定义通过 errors.Is(err, target error) 比较错误时的比较逻辑
func (e *ServerError) Is(target error) bool {
	// 将 target (interface) 转换为 *ServerError 类型, 并将结果存储在 t 中
	var t *ServerError
	ok := errors.As(target, &t)

	if !ok {
		return false
	}
	return e.Kind == t.Kind
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(msg string, err error) *ServerError {
	return &ServerError{
		Kind:    KindDatabase,
		Message: msg,
		Err:     err,
	}
}

// NewNotFoundError 创建未找到资源错误
func NewNotFoundError(msg string, err error) *ServerError {
	return &ServerError{
		Kind:    KindNotFound,
		Message: msg,
		Err:     err,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(msg string, err error) *ServerError {
	return &ServerError{
		Kind:    KindValidation,
		Message: msg,
		Err:     err,
	}
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(msg string, err error) *ServerError {
	return &ServerError{
		Kind:    KindUnauthorized,
		Message: msg,
		Err:     err,
	}
}
