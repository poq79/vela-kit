package vela

type Closer interface {
	Name() string
	Close() error
}

type Start interface {
	From(string) Start     //注册code
	Err(func(error)) Start //注册错误处理
	Do()                   //处理业务
}
