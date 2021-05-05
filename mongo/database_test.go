package gluamongo_mongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gluamongo "github.com/tengattack/gluamongo"
	"github.com/tengattack/gluamongo/bsonutil"
	lua "github.com/yuin/gopher-lua"
)

func TestGetDatabase(t *testing.T) {
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
		local mdb, err = mongoClient:getDatabase('admin');
		if err ~= nil then
		  error(err);
		end
		local name = mdb:getName();
		local names, err = mdb:getCollectionNames();
		mongoClient:disconnect();
		return name, err, names;
	`

	require.NoError(L.DoString(script))
	require.Equal(3, L.GetTop())
	assert.Equal("admin", L.ToString(1))
	assert.Equal(lua.LNil, L.Get(2))
	assert.NotEmpty(bsonutil.GetValue(L, 3))
}
