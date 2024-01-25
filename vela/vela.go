package vela

import (
	"github.com/vela-ssoc/vela-kit/stdutil"
)

func WithEnv(env Environment) {
	once.Do(func() {
		console := stdutil.New(stdutil.Console())
		defer func() {
			_ = console.Close()
		}()

		_G = env
	})
}

func GxEnv() Environment {
	return _G
}
