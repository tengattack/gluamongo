package gluamongo

import (
	mongo "github.com/tengattack/gluamongo/mongo"
	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	mongo.RegisterType(L)
	L.PreloadModule("mongo", mongo.Loader)
}
