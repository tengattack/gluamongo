package gluamongo_mongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mongo "github.com/tengattack/gluamongo/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUnmarshalBSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	val, err := mongo.UnmarshalBSON("")
	require.Equal(mongo.ErrInvalidBSON, err)
	assert.Nil(val)

	val, err = mongo.UnmarshalBSON(`{"a": 1, "b": 2}`)
	require.NoError(err)
	assert.Equal(bson.D{bson.E{Key: "a", Value: int32(1)}, bson.E{Key: "b", Value: int32(2)}}, val)

	val, err = mongo.UnmarshalBSON(`["a", 1]`)
	require.NoError(err)
	assert.Equal(bson.A{"a", int32(1)}, val)
}
