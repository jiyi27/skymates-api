package handler

import (
	"regexp"
)

func isValidEmail(email string) bool {
	// ^[a-zA-Z0-9] 要求以字母或数字开头
	// [a-zA-Z0-9._%+-]* 允许中间部分有字母、数字、符号等
	pattern := `^[a-zA-Z0-9][a-zA-Z0-9._%+-]*@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(pattern).MatchString(email)
}
