package bsonutil

import (
	lua "github.com/yuin/gopher-lua"
)

var exports = map[string]lua.LGFunction{
	"ObjectID": NewObjectID,
}

// RegisterType registers bson types
func RegisterType(L *lua.LState) {
	mtObjectID := L.NewTypeMetatable(OBJECTID_TYPENAME)
	L.SetField(mtObjectID, "__index", L.SetFuncs(L.NewTable(), objectIDMethods))
	L.SetField(mtObjectID, "__eq", L.NewFunction(objectIDEqMethod))
	L.SetField(mtObjectID, "__tostring", L.NewFunction(objectIDToStringMethod))
}
