# pipe
内置pipe操作接口

## pipe.check
- go 校验pipe逻辑
- []pipe = pipe.Check(L)
```golang
    func() {
        pv := pipe.Check(L)
        pipe.Do(pv , value, L , function(err) end)
    }
```

## pipe.LValue
- go 校验pipe逻辑
- pv = pipe.LValue(val)
```golang
   pv := pipe.LValue(L)
   pipe.Do([]pipe{pv} , value, L , function(err) end)
```
## pipe.do
- GO 调用pipe对象
- pipe.Do([]pipe , interface{} , state , func(error))
```golang
    pv := pipe.Check(L)
    vl := lua.LString("helo")
    pipe.Do(pv , vl , L , func(err error) 
	    xEnv.Errorf("%v" , err)
	end)
```