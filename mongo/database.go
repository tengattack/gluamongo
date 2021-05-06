package gluamongo_mongo

import (
	"github.com/tengattack/gluamongo/bsonutil"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE_TYPENAME = "mongo{database}"
)

// Database mongo
type Database struct {
	Client   *Client
	Database *mongo.Database
}

var databaseMethods = map[string]lua.LGFunction{
	// "drop":                nil,
	"getCollection":      databaseGetCollectionMethod,
	"getCollectionNames": databaseGetCollectionNamesMethod,
	"getName":            databaseGetNameMethod,
}

func pushDatabase(L *lua.LState, client *Client, database *mongo.Database) {
	ud := L.NewUserData()
	ud.Value = &Database{
		Client:   client,
		Database: database,
	}
	L.SetMetatable(ud, L.GetTypeMetatable(DATABASE_TYPENAME))
	L.Push(ud)
}

func checkDatabase(L *lua.LState) *Database {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Database); ok {
		return v
	}
	L.ArgError(1, "mongo database expected")
	return nil
}

func databaseGetCollectionMethod(L *lua.LState) int {
	db := checkDatabase(L)

	collname := L.ToString(2)
	if collname == "" {
		L.ArgError(2, "collname required")
		return 0
	}

	mColl := db.Database.Collection(collname)
	pushCollection(L, db.Client, mColl)
	return 1
}

func databaseGetCollectionNamesMethod(L *lua.LState) int {
	db := checkDatabase(L)

	ctx, cancel := db.Client.Context()
	defer cancel()

	options := bsonutil.ToBSON(L, 2)
	if options == nil {
		options = bson.M{}
	}
	names, err := db.Database.ListCollectionNames(ctx, options)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(bsonutil.ToLuaValue(L, names))
	return 1
}

func databaseGetNameMethod(L *lua.LState) int {
	db := checkDatabase(L)

	name := db.Database.Name()
	L.Push(lua.LString(name))
	return 1
}
