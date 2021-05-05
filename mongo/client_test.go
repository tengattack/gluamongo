package gluamongo_mongo_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gluamongo "github.com/tengattack/gluamongo"
	"github.com/tengattack/gluamongo/bsonutil"
	lua "github.com/yuin/gopher-lua"
)

func getLuaMongoConnection() string {
	return fmt.Sprintf(`
		local mongo = require 'mongo';
		local mongoClient = mongo.Client()
		local ok, err = mongoClient:connect('%s');
	`, "mongodb://localhost:27017/admin")
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluamongo.Preload(L)

	script := getLuaMongoConnection() + `
		local ok1 = false;
		local ok2 = false;
		local dbnames = nil;
		local err2 = nil;
		if mongoClient ~= nil then
		  ok1 = mongoClient:set_timeout(20000);
		  dbnames, err2 = mongoClient:getDatabaseNames();
		  ok2 = mongoClient:disconnect();
		end
		return err, ok, ok1, err2, dbnames, ok2;
	`

	require.NoError(L.DoString(script))
	require.Equal(6, L.GetTop())
	assert.Equal(lua.LNil, L.Get(1))
	assert.Equal(lua.LTrue, L.Get(2))
	assert.Equal(lua.LTrue, L.Get(3))
	assert.Equal(lua.LNil, L.Get(4))
	assert.NotEmpty(bsonutil.GetValue(L, 5))
	assert.Equal(lua.LTrue, L.Get(6))
}
