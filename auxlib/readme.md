# auxlib

内置调用API

## cast

- 常见格式化字符
- [auxlib.ToBool(v)]()
- [auxlib.ToTime(v)]()
- [auxlib.ToDurtion(v)]()
- [auxlib.ToFloat64(v)]()
- [auxlib.ToFloat32(v)]()
- [auxlib.ToInt64(v)]()
- [auxlib.ToInt32(v)]()
- [auxlib.ToInt16(v)]()
- [auxlib.ToInt8(v)]()
- [auxlib.ToUint64(v)]()
- [auxlib.ToUint32(v)]()
- [auxlib.ToUint16(v)]()
- [auxlib.ToUint8(v)]()
- [auxlib.ToString(v)]()

```golang
    v := auxlib.ToBool("true")
    v := auxlib.ToInt("10")
    
```

- 下面方法会报错
- [auxlib.ToBoolE(v)]()
- [auxlib.ToTimeE(v)]()
- [auxlib.ToDurtionE(v)]()
- [auxlib.ToFloat64E(v)]()
- [auxlib.ToFloat32E(v)]()
- [auxlib.ToInt64E(v)]()
- [auxlib.ToInt32E(v)]()
- [auxlib.ToInt16E(v)]()
- [auxlib.ToInt8E(v)]()
- [auxlib.ToUint64E(v)]()
- [auxlib.ToUint32E(v)]()
- [auxlib.ToUint16E(v)]()
- [auxlib.ToUint8E(v)]()
- [auxlib.ToStringE(v)]()

```golang
    v , e := auxlib.ToBoolE("tre")
    v , e := auxlib.ToIntE("aa")
```


## sprintf
- v = auxlib.Format(L , seek)
- lua stack value sprintf 方法
- seek 从第几个参数开始

```golang
    -- 3: bb
    -- 2: cc
    -- 1: "%s %s" 

    v := auxlib.Format(L , 0)
    print(v) -- cc bb

```

## 字符串处理
- []byte = auxlib.S2B(string)
- string = auxlib.B2S([]byte)
```golang
    data   = auxlib.S2B("str")              -- []byte{'s' , 't' , 'r'}
    string = auxlib.S2B([]byte{'a' , 'b'})  -- ab

```

## 字符判断

- auxlib.IsInt(v)
- auxlib.IsChar(v)
```golang
    print(auxlib.IsInt('0')) -- true
    print(auxlib.IsInt('a')) -- false 

    print(auxlib.IsChar('a')) -- true
    print(auxlib.IsChar('0')) -- false 

```

## proc 名称判断 
- error = auxlib.Name(string)
- 判断进程名
```golang
    print(auxlib.Name("---aa")) -- invalid name
```

## 网络判断
- bool = auxlib.Ipv4(v)
- bool = auxlib.Ipv6(v)

## 目录判断
- string = auxlib.CheckDir(*LState , LValue)

## URL 判断
- URL = auxlib.CheckURL(LValue , *LState)
- URL , err = auxlib.NewURL(string)
#### 内置功能
- int    = URL.Int(string)
- bool   = URL.Bool(string)
- int    = URL.Port()
- []int  = URL.Ports()
- string = URL.Scheme()
- string = URL.Host()
- string = URL.Hostname()
- bool   = URL.IsNil()
- bool   = URL.V4()
- bool   = URL.V6()
- string = URL.Value(string)
```golang
    raw := "tcp://0.0.0.0:53?port=93, 100,200,667&exclude=99,94&aa=10&str=helo"
    URL , err = auxlib.NewURL(raw)

    URL.Port()       -- 53
    URL.Scheme()     -- tcp
    URL.V4()         -- true
    URL.V6()         -- false
    URL.Host()       -- 0.0.0.0:53
    URL.Hostname()   -- 0.0.0.0
    URL.Ports()      -- 93,95,..,98,100,200 , 667
    URL.Int("aa")    -- 10
    URL.Value("str") -- helo

```

## 文件hash
- string , error = auxlib.FileMd5(string)
- error  = auxlib.CheckSum(path , hash)
- writer = auxlib.CheckWriter(LValue , *LState) 