package bsonutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestDateTime(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	Preload(L)

	script := `
		local bson = require 'bson'
		local dt1 = bson.DateTime()
		local dt2 = bson.DateTime(1620277291038)
		local dt3 = bson.DateTime(1620277291038)
		print(dt2)
		return dt1 == dt2, dt1 == 0, dt1 ~= dt2, dt2 == dt3
	`
	require.NoError(L.DoString(script))
	require.Equal(4, L.GetTop())
	assert.Equal(lua.LFalse, L.Get(1))
	assert.Equal(lua.LFalse, L.Get(2))
	assert.Equal(lua.LTrue, L.Get(3))
	assert.Equal(lua.LTrue, L.Get(4))

	script = `
		local bson = require 'bson'
		local dt0 = bson.DateTime('invalid')
		return dt0
	`
	require.Error(L.DoString(script))
}

func TestTimestamp(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	Preload(L)

	script := `
		local bson = require 'bson'
		local ts1 = bson.Timestamp()
		local ts2 = bson.Timestamp(1620277291, 1)
		local ts3 = bson.Timestamp(1620277291, 1)
		print(ts2)
		return ts1 == ts2, ts1 == 0, ts1 ~= ts2, ts2 == ts3
	`
	require.NoError(L.DoString(script))
	require.Equal(4, L.GetTop())
	assert.Equal(lua.LFalse, L.Get(1))
	assert.Equal(lua.LFalse, L.Get(2))
	assert.Equal(lua.LTrue, L.Get(3))
	assert.Equal(lua.LTrue, L.Get(4))

	script = `
		local bson = require 'bson'
		local ts0 = bson.Timestamp(0)
		return ts0
	`
	err := L.DoString(script)
	require.Error(err)
	require.Contains(err.Error(), "Timestamp needs 0 or 2 arguments")
}

func TestNull(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	Preload(L)

	script := `
		local bson = require 'bson'
		local dt1 = bson.Null
		local dt2 = bson.Null
		print(dt1)
		return dt1 == dt2, dt1 == 0, dt1 ~= dt2
	`
	require.NoError(L.DoString(script))
	require.Equal(3, L.GetTop())
	assert.Equal(lua.LTrue, L.Get(1))
	assert.Equal(lua.LFalse, L.Get(2))
	assert.Equal(lua.LFalse, L.Get(3))
}
