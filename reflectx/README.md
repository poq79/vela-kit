# xreflect
这个是一个golang 到lua的反射库， 基于[gopher-luar](https://github.com/yuin/gopher-lua)

## xreflect.ToLValue(interface{} , *lua.LState) 
将常用的golang中的数据映射到lua变量中去,默认采用userdata的方法利用的是reflect反射到metatable中去，
绑定方法
```go
type User struct {
    Name  string
    token string
}

func (u *User) SetToken(t string) {
    u.token = t
}

func (u *User) Token() string {
    return u.token
}

func main() {
	 xEnv.Set("demo" , xreflect.ToLValue( u , nil))
}

/*
    lua代码演示:
    demo.name = "helo"
    demo.token = "x-xx-x-0"
    demo:SetToken("Helo") 注意此时要用冒号
*/
```

## xreflect.ToStruct(lua.Table , interface{})
主要是为了方便用户操作lua.table 映射到go中的struct的字段， 减少配置关系映射代码
- tag  分为以下几种关键字 lua , type
- lua  字段 格式默认用逗号分割 主要是 name,value 其中第二个数值是如果为默认值
- type 描述字段的类型 现在只支持 string,int,bool ,object ; 其中 object是对象 类似struct
_注意: 目前之支持struct映射 不支持二级嵌套_
```go
    type config struct {
	    Name string           `lua:"name"         type:"string"`
	    Age  int              `lua:"age,18"'      type:"int"`
	    Enable bool           `lua:"enable,false" type:"bool"`
	    Time  time.Duration   `lua:"time,5s"`
	    Sdk  lua.Writer        `lua:"sdk"         type:"object"`
    }
    
    func New(L *lua.LState) {
    	cfg := &config{}
    	tab := L.CheckTable(1)
    	
    	e := xreflect.ToStruct(tab , cfg)
    	if e != nil { return }
        	
    }
```