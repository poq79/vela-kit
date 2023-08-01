package lua

import "io"

type Console interface {
	Println(string)
	Printf(string, ...interface{})
	Invalid(string, ...interface{})
}

type LVFace interface {
	ToLValue() LValue
}

type Writer interface {
	VelaEntry
	io.Writer
}

type IO interface {
	VelaEntry
	io.Writer
	io.Reader
}

type Reader interface {
	VelaEntry
	io.Reader
}

type Closer interface {
	VelaEntry
	io.Closer
}

type ReaderCloser interface {
	VelaEntry
	io.Reader
	io.Closer
}

type WriterCloser interface {
	VelaEntry
	io.Writer
	io.Closer
}

type IndexEx interface {
	Index(*LState, string) LValue
}

type NewIndexEx interface {
	NewIndex(*LState, string, LValue)
}

// MetaTableEx 通过 string 获取 内容 a:key()
type MetaTableEx interface {
	MetaTable(*LState, string) LValue
}

// MetaEx  通过 LValue 获取 a[1]
type MetaEx interface {
	Meta(*LState, LValue) LValue
}

// NewMetaEx 通过LValue 设置 a[1] = 123
type NewMetaEx interface {
	NewMeta(*LState, LValue, LValue)
}
