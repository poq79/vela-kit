package xreflect

import (
	"reflect"

	"github.com/vela-ssoc/vela-kit/lua"
)

// Config is used to define xreflect behaviour for a particular *lua.LState.
type Config struct {
	// The name generating function that defines under which names Go
	// struct fields will be accessed.
	//
	// If nil, the default behaviour is used:
	//   - if the "xreflect" tag of the field is "", the field name and its name
	//     with a lowercase first letter is returned
	//  - if the tag is "-", no name is returned (i.e. the field is not
	//    accessible)
	//  - for any other tag value, that value is returned
	FieldNames func(s reflect.Type, f reflect.StructField) []string

	// The name generating function that defines under which names Go
	// methods will be accessed.
	//
	// If nil, the default behaviour is used:
	//   - the method name and its name with a lowercase first letter
	MethodNames func(t reflect.Type, m reflect.Method) []string

	regular map[reflect.Type]*lua.LTable
	types   *lua.LTable
}

func newConfig() *Config {
	return &Config{
		regular: make(map[reflect.Type]*lua.LTable),
	}
}

// GetConfig returns the xreflect configuration options for the given *lua.LState.
func GetConfig(L *lua.LState) *Config {
	const registryKey = "github.com/layeh/gopher-xreflect"

	registry := L.Get(lua.RegistryIndex).(*lua.LTable)
	lConfig, ok := registry.RawGetString(registryKey).(*lua.LUserData)
	if !ok {
		lConfig = L.NewUserData()
		lConfig.Value = newConfig()
		registry.RawSetString(registryKey, lConfig)
	}
	return lConfig.Value.(*Config)
}

func defaultFieldNames(s reflect.Type, f reflect.StructField) []string {
	const tagName = "xreflect"

	tag := f.Tag.Get(tagName)
	if tag == "-" {
		return nil
	}
	if tag != "" {
		return []string{tag}
	}
	return []string{
		f.Name,
		getUnexportedName(f.Name),
	}
}

func defaultMethodNames(t reflect.Type, m reflect.Method) []string {
	return []string{
		m.Name,
		getUnexportedName(m.Name),
	}
}
