package gluamongo_mongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gluamongo "github.com/tengattack/gluamongo"
	lua "github.com/yuin/gopher-lua"
)

func TestMongo(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluamongo.Preload(L)

	script := `
		return require('mongo')
	`
	assert.NoError(L.DoString(script))

	c := L.Get(1)
	assert.NotNil(c)
}
