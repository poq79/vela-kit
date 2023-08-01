package lua

import "errors"

var (
	InvalidVelaName  = errors.New("invalid name")
	AlreadyRun       = errors.New("already running")
	NotFoundCode     = errors.New("not found code")
	NotFoundVelaData = errors.New("not found vela data")
	InvalidVelaData  = errors.New("invalid vela data")
	InvalidTree      = errors.New("invalid tree")
	NotFoundTree     = errors.New("not found tree")
)

type VelaCode interface {
	Key() string
	NewVelaData(*LState, string, string) *VelaData
}

func velaCodeE(ls *LState) VelaCode {
	if ls.Exdata == nil {
		return nil
	}

	if vc, ok := ls.Exdata.(VelaCode); ok {
		return vc
	}
	return nil
}

func (ls *LState) CodeVM() string {
	vc := velaCodeE(ls)
	if vc == nil {
		return ""
	}
	return vc.Key()
}

func (ls *LState) CheckCodeVM(name string) bool {
	return ls.CodeVM() == name
}

func (ls *LState) NewVela(key string, typeof string) *VelaData {
	vc := velaCodeE(ls)
	if vc == nil {
		ls.RaiseError("new vela data are not allowed in vm without code")
		return nil
	}

	vla := vc.NewVelaData(ls, key, typeof)
	vla.code = vc.Key()
	vla.private = false
	return vla
}

func (ls *LState) NewVelaData(key string, typeof string) *VelaData {
	return ls.NewVela(key, typeof)
}

func (ls *LState) CheckVelaData(n int) *VelaData {
	v := ls.Get(n)
	if lv, ok := v.(*VelaData); ok {
		return lv
	}
	ls.TypeError(n, LTVelaData)
	return nil
}
