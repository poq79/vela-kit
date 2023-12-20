package lua

type OptionEx struct {
	State    *LState
	OnAfter  func() error
	OnBefore func() error
}

type OptionFunc func(*OptionEx)

func NewOption(v ...OptionFunc) *OptionEx {
	opt := &OptionEx{}

	for _, fn := range v {
		fn(opt)
	}
	return opt
}

func WithState(co *LState) OptionFunc {
	return func(opt *OptionEx) {
		opt.State = co
	}
}

func After(fn func() error) OptionFunc {
	return func(opt *OptionEx) {
		opt.OnAfter = fn
	}
}
func Before(fn func() error) OptionFunc {
	return func(opt *OptionEx) {
		opt.OnBefore = fn
	}
}
