package lua

import (
	"os"
)

var CompatVarArg = true

// var FieldsPerFlush = 50
// var RegistrySize = 256 * 20
// var RegistryGrowStep = 32
// var CallStackSize = 256
// var MaxTableGetLoop = 100
// var MaxArrayIndex = 67108864

var FieldsPerFlush = 32
var RegistrySize = 128
var RegistryGrowStep = 32
var CallStackSize = 128
var MaxTableGetLoop = 100
var MaxArrayIndex = 10001

type LNumber float64
type LInt int
type LUint uint
type LInt64 int64
type LUint64 uint64

const LNumberBit = 64
const LNumberScanFormat = "%f"
const LuaVersion = "Lua 5.1"

const LSNull = LString("")

var LuaPath = "LUA_PATH"
var LuaLDir string
var LuaPathDefault string
var LuaOS string
var LuaDirSep string
var LuaPathSep = ";"
var LuaPathMark = "?"
var LuaExecDir = "!"
var LuaIgMark = "-"

func init() {
	if os.PathSeparator == '/' { // unix-like
		LuaOS = "unix"
		LuaLDir = "/usr/local/share/lua/5.1"
		LuaDirSep = "/"
		LuaPathDefault = "./?.lua;" + LuaLDir + "/?.lua;" + LuaLDir + "/?/init.lua"
	} else { // windows
		LuaOS = "windows"
		LuaLDir = "!\\lua"
		LuaDirSep = "\\"
		LuaPathDefault = ".\\?.lua;" + LuaLDir + "\\?.lua;" + LuaLDir + "\\?\\init.lua"
	}
}
