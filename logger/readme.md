# logger
默认日志配置库 控制框架中的日志输出

## rock.logger
启动logger的配置

- level 日志等级
- filename 日志输出文件
- console 是否打开console
- format  日志输出格式 如:text , json
- compress 是否开启压缩 zip,lz4,snappy

```lua
rock.logger{
    level = "DEBUG", 
    filename = "resource/logs/rock.error.log",
    console = true,
    format = "text"
}
```