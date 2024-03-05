package strutil

import "strings"

func Expr(str string) (string, bool) { //v(abc)
	if !strings.HasPrefix(str, "v(") {
		return "", false
	}

	n := len(str)
	if str[n-1] != ')' {
		return "", false
	}

	val := str[2 : n-1]
	if len(val) == 0 {
		return "", true
	}
	return val, true
}
