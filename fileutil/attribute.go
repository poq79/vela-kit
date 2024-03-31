package fileutil

import (
	"github.com/vela-ssoc/vela-kit/exception"
	"os"
	"path/filepath"
	"time"
)

type Attr struct {
	Filename string
	Dir      bool
	CTime    time.Time
	MTime    time.Time
	ATime    time.Time
	Size     int64
}

func Attribute(filename string) (Attr, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return Attr{Filename: filename}, err
	}

	at, mt, ct, size, e := StateByInfo(info)
	return Attr{
		Filename: filename,
		Dir:      info.IsDir(),
		CTime:    ct,
		MTime:    mt,
		ATime:    at,
		Size:     size,
	}, e
}

func Glob(pattern string, ignore func(Attr) bool, errFn func(error)) []Attr {
	var tab []Attr
	elem, err := filepath.Glob(pattern)
	if err != nil {
		errFn(err)
		return tab
	}

	n := len(elem)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		filename := elem[i]
		attr, e := Attribute(filename)

		if e != nil {
			errFn(e)
			continue
		}

		if ignore(attr) {
			continue
		}

		tab = append(tab, attr)
	}
	return tab
}

func Clean(attrs []Attr, last int, max int64, space int64) error {
	n := len(attrs)
	errs := exception.New()
	cleanFn := func(filename string) {
		err := os.Remove(filename)
		if err != nil {
			errs.Try(filename, err)
		}
	}

	var size int64
	for i := 0; i < n; i++ {
		size += attrs[i].Size
	}

	if space > 0 && space < size { //空间占用 统计 如果为0 忽略统计
		return nil
	}

	for i := 0; i < n; i++ {
		attr := attrs[i]
		if i < n-last { //
			cleanFn(attr.Filename)
		}

		if attr.Size > max { //删除过大文件
			cleanFn(attr.Filename)
		}
	}

	return errs.Wrap()
}
