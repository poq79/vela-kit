package webdav

import (
	"github.com/vela-ssoc/vela-kit/grep"
	"path/filepath"
	"strings"
)

func Prefix(sub string) func(string) bool {
	return func(s string) bool {
		return strings.HasPrefix(s, sub)
	}
}

func Suffix(sub string) func(string) bool {
	return func(s string) bool {
		return strings.HasSuffix(s, sub)
	}
}

func Ext(ext ...string) func(string) bool {
	return func(s string) bool {
		e := filepath.Ext(s)
		for _, item := range ext {
			if e == item {
				return true
			}
		}
		return false
	}
}

func Grep(pat string) func(string) bool {
	return grep.New(pat)
}
