package lua

import "time"

// SuperVelaData 防止过多的方法定义
type SuperVelaData struct {
	Uptime time.Time
	Status VelaState
	TypeOf string
	code   string
}

func (sv *SuperVelaData) Init(typeof string) {
	sv.Status = VTInit
	sv.TypeOf = typeof
}

func (sv *SuperVelaData) vm(L *LState) {
	sv.code = L.CodeVM()
}

func (sv *SuperVelaData) CodeVM() string {
	return sv.code
}

func (sv *SuperVelaData) IsRun() bool {
	return sv.Status == VTRun
}

func (sv *SuperVelaData) IsPanic() bool {
	return sv.Status == VTPanic
}

func (sv *SuperVelaData) IsInit() bool {
	return sv.Status == VTInit
}

func (sv *SuperVelaData) IsClose() bool {
	return sv.Status == VTClose
}

func (sv *SuperVelaData) V(opts ...interface{}) {
	for _, item := range opts {

		switch v := item.(type) {

		//设置启动时间
		case time.Time:
			sv.Uptime = v
		//设置类型
		case string:
			sv.TypeOf = v
		//设置数据类型
		case VelaState:
			sv.Status = v

		default:

		}
	}
}

func (sv *SuperVelaData) Type() string { return sv.TypeOf }
func (sv *SuperVelaData) Name() string { return "" }

func (sv *SuperVelaData) NewMeta(*LState, LValue, LValue)  {}
func (sv *SuperVelaData) Meta(*LState, LValue) LValue      { return LNil }
func (sv *SuperVelaData) Index(*LState, string) LValue     { return LNil }
func (sv *SuperVelaData) NewIndex(*LState, string, LValue) {}

func (sv *SuperVelaData) Show(out Console) {
	out.Println("请定义对象的Show方法 ,如： func(a *A) Show( out lua.Console)")
}
func (sv *SuperVelaData) Help(out Console) {
	out.Println("请定义对象的Help方法 ,如： func(a *A) Help( out lua.Console)")
}

func (sv *SuperVelaData) State() VelaState { return sv.Status }
