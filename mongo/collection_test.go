package gluamongo_mongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gluamongo "github.com/tengattack/gluamongo"
	"github.com/tengattack/gluamongo/bsonutil"
	lua "github.com/yuin/gopher-lua"
)

func TestGetCollection(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluamongo.Preload(L)

	script := getLuaMongoConnection() + `
		if err ~= nil then
		  error(err);
		end
		local mcoll, err = mongoClient:getCollection('test', 'test');
		if err ~= nil then
		  error(err);
		end
		local name = mcoll:getName();
		mongoClient:disconnect();
		return name, err;
	`

	require.NoError(L.DoString(script))
	require.Equal(2, L.GetTop())
	assert.Equal("test", L.ToString(1))
	assert.Equal(lua.LNil, L.Get(2))
}

func TestInsertUpdateRemove(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluamongo.Preload(L)

	script := getLuaMongoConnection() + `
		if err ~= nil then
		  error(err);
		end
		local mcoll, err = mongoClient:getCollection('test', 'test');
		if err ~= nil then
		  error(err);
		end
		mcoll:remove({}); -- remove all
		local res, err = mcoll:insert({a = 1});
		local res2, err2 = mcoll:update({a = 1}, '{"$set": {"a": 2}}', {multi = true});
		local res3, err3 = mcoll:remove({a = 2});
		mongoClient:disconnect();
		return res, err, res2, err2, res3, err3
	`

	require.NoError(L.DoString(script))
	require.Equal(6, L.GetTop())
	assert.Equal(map[string]interface{}{"nInserted": 1}, bsonutil.GetValue(L, 1))
	assert.Equal(lua.LNil, L.Get(2))
	assert.Equal(map[string]interface{}{"nMatched": 1, "nModified": 1, "nUpserted": 0}, bsonutil.GetValue(L, 3))
	assert.Equal(lua.LNil, L.Get(4))
	assert.Equal(map[string]interface{}{"nRemoved": 1}, bsonutil.GetValue(L, 5))
	assert.Equal(lua.LNil, L.Get(6))
}

func TestInsertFindRemove(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluamongo.Preload(L)

	script := getLuaMongoConnection() + `
		if err ~= nil then
		  error(err);
		end
		local mcoll, err = mongoClient:getCollection('test', 'test');
		if err ~= nil then
		  error(err);
		end
		mcoll:remove({}); -- remove all
		local res, err = mcoll:insert({a = 1});
		local res2, err2 = mcoll:find({a = 1});
		local res3, err3 = mcoll:remove({a = 1});
		mongoClient:disconnect();
		return res, err, res2, err2, res3, err3
	`

	require.NoError(L.DoString(script))
	require.Equal(6, L.GetTop())
	assert.Equal(map[string]interface{}{"nInserted": 1}, bsonutil.GetValue(L, 1))
	assert.Equal(lua.LNil, L.Get(2))
	assert.NotEmpty(bsonutil.GetValue(L, 3))
	assert.Equal(lua.LNil, L.Get(4))
	assert.Equal(map[string]interface{}{"nRemoved": 1}, bsonutil.GetValue(L, 5))
	assert.Equal(lua.LNil, L.Get(6))
}
