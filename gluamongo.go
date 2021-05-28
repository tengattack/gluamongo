package gluamongo

import (
	mongo "github.com/tengattack/gluamongo/mongo"
	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule("mongo", mongo.Loader)
}
