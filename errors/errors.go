package errors

import (
	"errors"
	"fmt"
)

// ErrorKind 用来区分不同类型的业务错误
type ErrorKind int

const (
	KindInternal      ErrorKind = iota // 系统内部错误
	KindNotFound                       // 资源未找到
	KindAlreadyExists                  // 资源已存在
	KindValidation                     // 参数校验失败
	KindUnauthorized                   // 需要认证
	KindForbidden                      // 权限不足
	KindConflict                       // 冲突，比如悲观锁、版本号不一致等
)

// ServerError 是所有可预知业务错误的统一类型
type ServerError struct {
	Kind    ErrorKind // 错误类型
	Message string    // 内部人可读的简短描述
	Err     error     // 底层原始错误（可选）
}

// Error 实现了 error 接口
func (e *ServerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 返回底层错误，方便 errors.Unwrap / errors.Is 链式比较
func (e *ServerError) Unwrap() error {
	return e.Err
}

// Is 支持 errors.Is(err, target) 直接比较 Kind
func (e *ServerError) Is(target error) bool {
	var t *ServerError
	if !errors.As(target, &t) {
		return false
	}
	return e.Kind == t.Kind
}

// --- 常见错误类型构造函数 ---

func NewInternalError(msg string, err error) *ServerError {
	return &ServerError{Kind: KindInternal, Message: msg, Err: err}
}

func NewNotFoundError(msg string, err error) *ServerError {
	return &ServerError{Kind: KindNotFound, Message: msg, Err: err}
}

func NewAlreadyExistsError(msg string, err error) *ServerError {
	return &ServerError{Kind: KindAlreadyExists, Message: msg, Err: err}
}

func NewValidationError(msg string, err error) *ServerError {
	return &ServerError{Kind: KindValidation, Message: msg, Err: err}
}

func NewUnauthorizedError(msg string, err error) *ServerError {
	return &ServerError{Kind: KindUnauthorized, Message: msg, Err: err}
}

func NewForbiddenError(msg string, err error) *ServerError {
	return &ServerError{Kind: KindForbidden, Message: msg, Err: err}
}

func NewConflictError(msg string, err error) *ServerError {
	return &ServerError{Kind: KindConflict, Message: msg, Err: err}
}
