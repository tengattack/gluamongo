package gluamongo_mongo

import (
	"context"
	"time"

	"github.com/tengattack/gluamongo/bsonutil"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CLIENT_TYPENAME = "mongo{client}"
)

// Client mongo
type Client struct {
	Client  *mongo.Client
	Timeout time.Duration
}

func (client *Client) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), client.Timeout)
}

func newClient(L *lua.LState) int {
	timeout := 10 * time.Second
	ud := L.NewUserData()
	ud.Value = &Client{
		Client:  nil,
		Timeout: timeout,
	}
	L.SetMetatable(ud, L.GetTypeMetatable(CLIENT_TYPENAME))
	L.Push(ud)
	return 1
}

var clientMethods = map[string]lua.LGFunction{
	"set_timeout": clientSetTimeoutMethod,
	"connect":     clientConnectMethod,
	"disconnect":  clientDisconnectMethod,

	"getCollection":    clientGetCollectionMethod,
	"getDatabase":      clientGetDatabaseMethod,
	"getDatabaseNames": clientGetDatabaseNamesMethod,
}

func checkClient(L *lua.LState) *Client {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Client); ok {
		return v
	}
	L.ArgError(1, "mongo client expected")
	return nil
}

func clientConnectMethod(L *lua.LState) int {
	client := checkClient(L)

	dsn := L.ToString(2)
	if dsn == "" {
		L.ArgError(2, "dsn required")
		return 0
	}

	if client.Client != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("mongo client already connected"))
		return 1
	}

	ctx, cancel := client.Context()
	defer cancel()
	opts := options.Client().ApplyURI(dsn)
	mongoClient, err := mongo.Connect(ctx, opts)

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		_ = mongoClient.Disconnect(ctx)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	client.Client = mongoClient

	L.Push(lua.LBool(true))
	return 1
}

func clientDisconnectMethod(L *lua.LState) int {
	client := checkClient(L)

	if client.Client == nil {
		L.Push(lua.LBool(true))
		return 1
	}

	ctx, cancel := client.Context()
	defer cancel()
	err := client.Client.Disconnect(ctx)
	// always clean
	client.Client = nil
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

func clientSetTimeoutMethod(L *lua.LState) int {
	client := checkClient(L)
	timeout := L.ToInt64(2) // timeout (in ms)

	client.Timeout = time.Millisecond * time.Duration(timeout)

	L.Push(lua.LBool(true))
	return 1
}

func clientGetCollectionMethod(L *lua.LState) int {
	client := checkClient(L)
	dbname := L.ToString(2)
	if dbname == "" {
		L.ArgError(2, "dbname required")
		return 0
	}
	collname := L.ToString(3)
	if collname == "" {
		L.ArgError(3, "collname required")
		return 0
	}

	mDb := client.Client.Database(dbname)
	mColl := mDb.Collection(collname)
	pushCollection(L, client, mColl)
	return 1
}

func clientGetDatabaseMethod(L *lua.LState) int {
	client := checkClient(L)
	dbname := L.ToString(2)
	if dbname == "" {
		L.ArgError(2, "dbname required")
		return 0
	}

	mDb := client.Client.Database(dbname)
	pushDatabase(L, client, mDb)
	return 1
}

func clientGetDatabaseNamesMethod(L *lua.LState) int {
	client := checkClient(L)

	ctx, cancel := client.Context()
	defer cancel()
	options := ToBSON(L, 2)
	if options == nil {
		options = bson.M{}
	}
	names, err := client.Client.ListDatabaseNames(ctx, options)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(bsonutil.ToLuaValue(L, names))
	return 1
}
