package gluamongo_mongo

import (
	"github.com/tengattack/gluamongo/bsonutil"
	lua "github.com/yuin/gopher-lua"
)

var exports = map[string]lua.LGFunction{
	"Client":    newClient,
	"ObjectID":  bsonutil.NewObjectID,
	"DateTime":  bsonutil.NewDateTime,
	"Timestamp": bsonutil.NewTimestamp,
}

// Loader mongo module loader
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)

	L.SetField(mod, "_DEBUG", lua.LBool(false))
	L.SetField(mod, "_VERSION", lua.LString("0.0.0"))

	registerType(L)

	// consts, after type registered
	L.SetField(mod, "Null", bsonutil.LNull(L))

	return 1
}

func registerType(L *lua.LState) {
	bsonutil.RegisterType(L)

	mtClient := L.NewTypeMetatable(CLIENT_TYPENAME)
	L.SetField(mtClient, "__index", L.SetFuncs(L.NewTable(), clientMethods))
	mtCollection := L.NewTypeMetatable(COLLECTION_TYPENAME)
	L.SetField(mtCollection, "__index", L.SetFuncs(L.NewTable(), collectionMethods))
	mtDatabase := L.NewTypeMetatable(DATABASE_TYPENAME)
	L.SetField(mtDatabase, "__index", L.SetFuncs(L.NewTable(), databaseMethods))
}
