package bsonutil

import (
	"errors"
	"strings"

	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrInvalidBSON = errors.New("invalid BSON")

var exports = map[string]lua.LGFunction{
	"ObjectID":  NewObjectID,
	"DateTime":  NewDateTime,
	"Timestamp": NewTimestamp,
}

// RegisterType registers bson types
func RegisterType(L *lua.LState) {
	mtObjectID := L.NewTypeMetatable(OBJECTID_TYPENAME)
	L.SetField(mtObjectID, "__index", L.SetFuncs(L.NewTable(), objectIDMethods))
	L.SetField(mtObjectID, "__eq", L.NewFunction(objectIDEqMethod))
	L.SetField(mtObjectID, "__tostring", L.NewFunction(objectIDToStringMethod))

	mtDateTime := L.NewTypeMetatable(DATETIME_TYPENAME)
	L.SetField(mtDateTime, "__index", L.SetFuncs(L.NewTable(), dateTimeMethods))
	L.SetField(mtDateTime, "__eq", L.NewFunction(dateTimeEqMethod))
	L.SetField(mtDateTime, "__tostring", L.NewFunction(dateTimeToStringMethod))

	mtTimestamp := L.NewTypeMetatable(TIMESTAMP_TYPENAME)
	L.SetField(mtTimestamp, "__index", L.SetFuncs(L.NewTable(), timestampMethods))
	L.SetField(mtTimestamp, "__eq", L.NewFunction(timestampEqMethod))
	L.SetField(mtTimestamp, "__tostring", L.NewFunction(timestampToStringMethod))

	mtNull := L.NewTypeMetatable(NULL_TYPENAME)
	L.SetField(mtNull, "__index", L.SetFuncs(L.NewTable(), nullMethods))
	L.SetField(mtNull, "__eq", L.NewFunction(nullEqMethod))
	L.SetField(mtNull, "__tostring", L.NewFunction(nullToStringMethod))
}

// UnmarshalBSON unmarshals extended json to bson
func UnmarshalBSON(str string) (interface{}, error) {
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "{") {
		// document
		var val bson.D
		err := bson.UnmarshalExtJSON([]byte(str), false, &val)
		return val, err
	} else if strings.HasPrefix(str, "[") {
		// array
		var val bson.A
		err := bson.UnmarshalExtJSON([]byte(str), false, &val)
		return val, err
	}
	return nil, ErrInvalidBSON
}

// CastBSON casts glua value to bson, nil if not a valid bson value
func CastBSON(L *lua.LState, idx int) interface{} {
	lv := L.Get(idx)
	switch lv.Type() {
	case lua.LTString:
		val, err := UnmarshalBSON(lua.LVAsString(lv))
		if err != nil {
			L.ArgError(idx, err.Error())
		}
		return val
	case lua.LTTable:
		val := GetValue(L, idx)
		if arr, ok := val.([]interface{}); ok {
			if len(arr) == 0 {
				// empty doc treats as {} instead of []
				return map[string]interface{}{}
			}
		}
		return val
	default:
		L.ArgError(idx, "string or table expected")
		return nil
	}
}

// ToBSON converts glua value to bson, allow nil
func ToBSON(L *lua.LState, idx int) interface{} {
	lv := L.Get(idx)
	if lv == lua.LNil {
		return nil
	}
	return CastBSON(L, idx)
}
