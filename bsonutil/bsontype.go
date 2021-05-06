package bsonutil

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// bson types
const (
	DATETIME_TYPENAME  = "bson{datetime}"
	TIMESTAMP_TYPENAME = "bson{timestamp}"
	NULL_TYPENAME      = "bson{null}"
)

// DateTime mongo
type DateTime struct {
	DT primitive.DateTime
}

// Timestamp mongo
type Timestamp struct {
	Ts primitive.Timestamp
}

// Null mongo
type Null struct {
}

var dateTimeMethods = map[string]lua.LGFunction{}
var timestampMethods = map[string]lua.LGFunction{}
var nullMethods = map[string]lua.LGFunction{}

// NewDateTime new DateTime for glua
func NewDateTime(L *lua.LState) int {
	dt := L.OptInt64(1, 0)

	ud := L.NewUserData()
	ud.Value = &DateTime{DT: primitive.DateTime(dt)}
	L.SetMetatable(ud, L.GetTypeMetatable(DATETIME_TYPENAME))
	L.Push(ud)
	return 1
}

// NewTimestamp new Timestamp for glua
func NewTimestamp(L *lua.LState) int {
	switch L.GetTop() {
	case 0, 2:
		// PASS
	default:
		L.ArgError(1, "Timestamp needs 0 or 2 arguments")
		return 0
	}
	t := L.OptInt(1, 0)
	i := L.OptInt(2, 0)

	if t == 0 && i == 0 {
		// mongodb will generates a new timestamp when saving documents
		// https://docs.mongodb.com/manual/reference/bson-types/#timestamps
		//
		// now := time.Now()
		// the most significant 32 bits are a time_t value (seconds since the Unix epoch)
		// t = int(now.Unix())
		// the least significant 32 bits are an incrementing ordinal for operations within a given second.
		// i = 1
	}

	ud := L.NewUserData()
	ud.Value = &Timestamp{Ts: primitive.Timestamp{T: uint32(t), I: uint32(i)}}
	L.SetMetatable(ud, L.GetTypeMetatable(TIMESTAMP_TYPENAME))
	L.Push(ud)
	return 1
}

// LDateTime creates DateTime value for glua
func LDateTime(L *lua.LState, dt primitive.DateTime) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &DateTime{DT: dt}
	L.SetMetatable(ud, L.GetTypeMetatable(DATETIME_TYPENAME))
	return ud
}

// LTimestamp creates Timestamp value for glua
func LTimestamp(L *lua.LState, ts primitive.Timestamp) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &Timestamp{Ts: ts}
	L.SetMetatable(ud, L.GetTypeMetatable(TIMESTAMP_TYPENAME))
	return ud
}

// LNull creates Null value for glua
func LNull(L *lua.LState) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &Null{}
	L.SetMetatable(ud, L.GetTypeMetatable(NULL_TYPENAME))
	return ud
}

func checkDateTime(L *lua.LState, idx int) *DateTime {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*DateTime); ok {
		return v
	}
	L.ArgError(1, "bson datetime expected")
	return nil
}

func checkTimestamp(L *lua.LState, idx int) *Timestamp {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*Timestamp); ok {
		return v
	}
	L.ArgError(1, "bson timestamp expected")
	return nil
}

func checkNull(L *lua.LState, idx int) *Null {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*Null); ok {
		return v
	}
	L.ArgError(1, "bson null expected")
	return nil
}

func dateTimeToStringMethod(L *lua.LState) int {
	dateTime := checkDateTime(L, 1)

	L.Push(lua.LString(fmt.Sprintf("DateTime(%d)", dateTime.DT)))
	return 1
}

func timestampToStringMethod(L *lua.LState) int {
	ts := checkTimestamp(L, 1)

	L.Push(lua.LString(fmt.Sprintf("Timestamp(%d, %d)", ts.Ts.T, ts.Ts.I)))
	return 1
}

func nullToStringMethod(L *lua.LState) int {
	_ = checkNull(L, 1)

	L.Push(lua.LString("null"))
	return 1
}

func dateTimeEqMethod(L *lua.LState) int {
	dateTime1 := checkDateTime(L, 1)
	dateTime2 := checkDateTime(L, 2) // REVIEW: ArgError required?

	L.Push(lua.LBool(dateTime1.DT == dateTime2.DT))
	return 1
}

func timestampEqMethod(L *lua.LState) int {
	ts1 := checkTimestamp(L, 1)
	ts2 := checkTimestamp(L, 2)

	L.Push(lua.LBool(ts1.Ts.Equal(ts2.Ts)))
	return 1
}

func nullEqMethod(L *lua.LState) int {
	_ = checkNull(L, 1)
	_ = checkNull(L, 2)

	// null always equal
	L.Push(lua.LTrue)
	return 1
}
