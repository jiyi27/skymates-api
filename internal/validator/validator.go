package validator

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	// validate 单例实例
	validate *validator.Validate
	once     sync.Once
)

func init() {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
	})
}

// ValidateRequest 校验传入的结构体 req
// 如果校验不通过，返回格式化后的错误字符串；
// 如果出现非校验类错误，则返回原始 error
// 校验通过时返回空字符串和 nil
func ValidateRequest(req interface{}) (string, error) {
	if err := validate.Struct(req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			return buildErrorString(errs), err
		}
		return "验证字段发生未知错误", err
	}
	return "", nil
}

// buildErrorString 将 validator.ValidationErrors 切片拼接为多行错误描述字符串
// 每行格式: 完整路径 (字段名): failed '校验标签' [ param=参数]
func buildErrorString(errs validator.ValidationErrors) string {
	var sb strings.Builder
	for _, e := range errs {
		// e.StructNamespace() 返回类似 Parent.Child.Field 的完整路径
		sb.WriteString(fmt.Sprintf("%s (%s): failed '%s'", e.StructNamespace(), e.Field(), e.Tag()))
		// 某些校验标签带参数，如 max=10
		if p := e.Param(); p != "" {
			sb.WriteString(" param=" + p)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
