package bsonutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestObjectID(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	Preload(L)

	script := `
		local bson = require 'bson'
		local oid1 = bson.ObjectID()
		local oid2 = bson.ObjectID('6092e50e4ed1be4939967323')
		local oid3 = bson.ObjectID('6092e50e4ed1be4939967323')
		print(oid1)
		return oid1 == oid2, oid1 == 0, oid1 ~= oid2, oid2 == oid3
	`
	require.NoError(L.DoString(script))
	require.Equal(4, L.GetTop())
	assert.Equal(lua.LFalse, L.Get(1))
	assert.Equal(lua.LFalse, L.Get(2))
	assert.Equal(lua.LTrue, L.Get(3))
	assert.Equal(lua.LTrue, L.Get(4))

	script = `
		local bson = require 'bson'
		local oid0 = bson.ObjectID('invalid')
		return oid0
	`
	require.Error(L.DoString(script))
}
