package vela

import "fmt"

func WithEnv(env Environment) {
	if _G == nil {
		_G = env
		fmt.Println("env constructor over..")
		return
	}
	fmt.Println("env already running")
}

func GxEnv() Environment {
	return _G
}
