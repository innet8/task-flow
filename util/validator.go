package util

import (
	"regexp"
)

// IsDoubleStr 验证是否是浮点字符串
func IsDoubleStr(doubelstr string) (bool, error) {
	return regexp.MatchString(`^[-+]?[0-9]+(\.[0-9]+)?$`, doubelstr)
}
