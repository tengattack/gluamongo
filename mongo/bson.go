package gluamongo_mongo

import (
	"errors"
	"strings"

	"github.com/tengattack/gluamongo/bsonutil"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrInvalidBSON = errors.New("invalid BSON")

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
		val := bsonutil.GetValue(L, idx)
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

func ToBSON(L *lua.LState, idx int) interface{} {
	lv := L.Get(idx)
	if lv == lua.LNil {
		return nil
	}
	return CastBSON(L, idx)
}
