package bsonutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)

	L.SetField(mod, "_DEBUG", lua.LBool(false))
	L.SetField(mod, "_VERSION", lua.LString("0.0.0"))

	RegisterType(L)

	// consts, after type registered
	L.SetField(mod, "Null", LNull(L))

	return 1
}

func Preload(L *lua.LState) {
	L.PreloadModule("bson", Loader)
}

func TestUnmarshalBSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	val, err := UnmarshalBSON("")
	require.Equal(ErrInvalidBSON, err)
	assert.Nil(val)

	val, err = UnmarshalBSON(`{"a": 1, "b": 2}`)
	require.NoError(err)
	assert.Equal(bson.D{bson.E{Key: "a", Value: int32(1)}, bson.E{Key: "b", Value: int32(2)}}, val)

	val, err = UnmarshalBSON(`["a", 1]`)
	require.NoError(err)
	assert.Equal(bson.A{"a", int32(1)}, val)
}

func TestCastToBSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	Preload(L)

	script := `
		return '{"a": 1, "b": 2}', {a = 1, b = 2}, {}, nil
	`
	require.NoError(L.DoString(script))
	require.Equal(4, L.GetTop())
	assert.Equal(primitive.D(primitive.D{primitive.E{Key: "a", Value: int32(1)}, primitive.E{Key: "b", Value: int32(2)}}), CastBSON(L, 1))
	assert.Equal(map[string]interface{}{"a": 1, "b": 2}, CastBSON(L, 2))
	assert.Equal(map[string]interface{}{}, CastBSON(L, 3))
	assert.Equal(map[string]interface{}{"a": 1, "b": 2}, ToBSON(L, 2))
	assert.Equal(map[string]interface{}{}, ToBSON(L, 3))
	assert.Equal(nil, ToBSON(L, 4))
	assert.Panics(func() {
		// ArgError here
		assert.Equal(nil, CastBSON(L, 4))
	})
}
