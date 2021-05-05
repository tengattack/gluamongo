package bsonutil

import (
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// bson types
const (
	OBJECTID_TYPENAME = "bson{objectid}"
)

// ObjectID mongo
type ObjectID struct {
	OID primitive.ObjectID
}

var objectIDMethods = map[string]lua.LGFunction{}

// NewObjectID new ObjectID for glua
func NewObjectID(L *lua.LState) int {
	str := L.OptString(1, "")

	var oid primitive.ObjectID

	if str != "" {
		if !primitive.IsValidObjectID(str) {
			L.ArgError(1, "invalid format")
			return 0
		}
		var err error
		oid, err = primitive.ObjectIDFromHex(str)
		if err != nil {
			L.ArgError(1, err.Error())
			return 0
		}
	} else {
		oid = primitive.NewObjectID()
	}

	ud := L.NewUserData()
	ud.Value = &ObjectID{OID: oid}
	L.SetMetatable(ud, L.GetTypeMetatable(OBJECTID_TYPENAME))
	L.Push(ud)
	return 1
}

// LObjectID creates ObjectID value for glua
func LObjectID(L *lua.LState, oid primitive.ObjectID) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &ObjectID{OID: oid}
	L.SetMetatable(ud, L.GetTypeMetatable(OBJECTID_TYPENAME))
	return ud
}

func checkObjectID(L *lua.LState, idx int) *ObjectID {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*ObjectID); ok {
		return v
	}
	L.ArgError(1, "bson objectid expected")
	return nil
}

func objectIDToStringMethod(L *lua.LState) int {
	objectID := checkObjectID(L, 1)

	L.Push(lua.LString(objectID.OID.String()))
	return 1
}

func objectIDEqMethod(L *lua.LState) int {
	objectID1 := checkObjectID(L, 1)
	objectID2 := checkObjectID(L, 2)

	L.Push(lua.LBool(objectID1.OID == objectID2.OID))
	return 1
}
