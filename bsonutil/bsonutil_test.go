package bsonutil

import (
	lua "github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)

	L.SetField(mod, "_DEBUG", lua.LBool(false))
	L.SetField(mod, "_VERSION", lua.LString("0.0.0"))

	RegisterType(L)
	return 1
}

func Preload(L *lua.LState) {
	L.PreloadModule("bson", Loader)
}
