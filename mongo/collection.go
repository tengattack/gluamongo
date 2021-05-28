package gluamongo_mongo

import (
	"fmt"

	"github.com/tengattack/gluamongo/bsonutil"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	COLLECTION_TYPENAME = "mongo{collection}"
)

// Collection mongo
type Collection struct {
	Client     *Client
	Collection *mongo.Collection
}

var collectionMethods = map[string]lua.LGFunction{
	// "drop":                nil,
	"aggregate": collectionAggregateMethod,
	"count":     collectionCountMethod,
	"find":      collectionFindMethod,
	"findOne":   collectionFindOneMethod,
	"getName":   collectionGetNameMethod,
	"insert":    collectionInsertMethod,
	"remove":    collectionRemoveMethod,
	"update":    collectionUpdateMethod,
}

func pushCollection(L *lua.LState, client *Client, collection *mongo.Collection) {
	ud := L.NewUserData()
	ud.Value = &Collection{
		Client:     client,
		Collection: collection,
	}
	L.SetMetatable(ud, L.GetTypeMetatable(COLLECTION_TYPENAME))
	L.Push(ud)
}

func checkCollection(L *lua.LState) *Collection {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Collection); ok {
		return v
	}
	L.ArgError(1, "mongo collection expected")
	return nil
}

func toNumber(v interface{}) (int64, error) {
	switch i := v.(type) {
	case int64:
		return i, nil
	case int:
		return int64(i), nil
	case int32:
		return int64(i), nil
	case float64:
		return int64(i), nil
	}
	return 0, fmt.Errorf("unknown value: %v (%T)", v, v)
}

func collectionFindOptions(opts interface{}) (*options.FindOptions, error) {
	fo := &options.FindOptions{}
	if opts == nil {
		return fo, nil
	}
	if m, ok := opts.(map[string]interface{}); ok {
		if v, ok := m["projection"]; ok {
			fo.SetProjection(v)
		}
		if v, ok := m["sort"]; ok {
			fo.SetSort(v)
		}
		if v, ok := m["skip"]; ok {
			i, err := toNumber(v)
			if err != nil {
				return nil, err
			}
			fo.SetSkip(i)
		}
		if v, ok := m["limit"]; ok {
			i, err := toNumber(v)
			if err != nil {
				return nil, err
			}
			fo.SetLimit(i)
		}
	} else if m, ok := opts.(bson.D); ok {
		for _, v := range m {
			if v.Key == "projection" {
				fo.SetProjection(v.Value)
			} else if v.Key == "sort" {
				fo.SetSort(v.Value)
			} else if v.Key == "skip" {
				i, err := toNumber(v)
				if err != nil {
					return nil, err
				}
				fo.SetSkip(i)
			} else if v.Key == "limit" {
				i, err := toNumber(v)
				if err != nil {
					return nil, err
				}
				fo.SetLimit(i)
			}
		}
	}
	return fo, nil
}

func collectionAggregateMethod(L *lua.LState) int {
	coll := checkCollection(L)

	query := bsonutil.CastBSON(L, 2)

	ctx, cancel := coll.Client.Context()
	defer cancel()

	cur, err := coll.Collection.Aggregate(ctx, query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var results []bson.M
	err = cur.All(ctx, &results)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(bsonutil.ToLuaValue(L, results))
	return 1
}

func collectionCountMethod(L *lua.LState) int {
	coll := checkCollection(L)

	query := bsonutil.CastBSON(L, 2)
	opts, err := collectionFindOptions(bsonutil.ToBSON(L, 3))
	if err != nil {
		L.ArgError(3, err.Error())
		return 0
	}
	countOptions := &options.CountOptions{}
	if opts != nil {
		if opts.Limit != nil {
			countOptions.SetLimit(*opts.Limit)
		}
		if opts.Skip != nil {
			countOptions.SetSkip(*opts.Skip)
		}
	}

	ctx, cancel := coll.Client.Context()
	defer cancel()

	count, err := coll.Collection.CountDocuments(ctx, query, countOptions)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(count))
	return 1
}

func collectionFindMethod(L *lua.LState) int {
	coll := checkCollection(L)

	query := bsonutil.CastBSON(L, 2)
	opts, err := collectionFindOptions(bsonutil.ToBSON(L, 3))
	if err != nil {
		L.ArgError(3, err.Error())
		return 0
	}

	ctx, cancel := coll.Client.Context()
	defer cancel()

	cur, err := coll.Collection.Find(ctx, query, opts)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var results []bson.M
	err = cur.All(ctx, &results)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(bsonutil.ToLuaValue(L, results))
	return 1
}

func collectionFindOneMethod(L *lua.LState) int {
	coll := checkCollection(L)

	query := bsonutil.CastBSON(L, 2)
	opts, err := collectionFindOptions(bsonutil.ToBSON(L, 3))
	if err != nil {
		L.ArgError(3, err.Error())
		return 0
	}
	foOptions := &options.FindOneOptions{}
	if opts != nil {
		if opts.Projection != nil {
			foOptions.SetProjection(opts.Projection)
		}
		if opts.Sort != nil {
			foOptions.SetSort(opts.Sort)
		}
		if opts.Skip != nil {
			foOptions.SetSkip(*opts.Skip)
		}
	}

	ctx, cancel := coll.Client.Context()
	defer cancel()

	res := coll.Collection.FindOne(ctx, query, foOptions)
	err = res.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var result bson.M
	err = res.Decode(&result)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(bsonutil.ToLuaValue(L, res))
	return 1
}

func collectionGetNameMethod(L *lua.LState) int {
	coll := checkCollection(L)

	name := coll.Collection.Name()
	L.Push(lua.LString(name))
	return 1
}

func newInsertResult(nInserted int) map[string]int {
	return map[string]int{"nInserted": nInserted}
}

func collectionInsertMethod(L *lua.LState) int {
	coll := checkCollection(L)

	doc := bsonutil.CastBSON(L, 2)

	ctx, cancel := coll.Client.Context()
	defer cancel()

	if arr, ok := doc.([]interface{}); ok {
		res, err := coll.Collection.InsertMany(ctx, arr)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(bsonutil.ToLuaValue(L, newInsertResult(len(res.InsertedIDs))))
		return 1
	} else {
		_, err := coll.Collection.InsertOne(ctx, doc)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(bsonutil.ToLuaValue(L, newInsertResult(1)))
		return 1
	}
}

func newRemoveResult(nRemoved int) map[string]int {
	return map[string]int{"nRemoved": nRemoved}
}

func collectionRemoveMethod(L *lua.LState) int {
	coll := checkCollection(L)

	query := bsonutil.CastBSON(L, 2)
	var justOne bool
	lv := L.Get(3)
	if lv.Type() == lua.LTBool {
		justOne = lua.LVAsBool(lv)
	} else {
		options := bsonutil.ToBSON(L, 3)
		if options != nil {
			// TODO: bson.D
			if v, ok := options.(map[string]interface{}); ok {
				if v2, ok2 := v["justOne"]; ok2 {
					if justOneVal, ok3 := v2.(bool); ok3 {
						justOne = justOneVal
					} else {
						L.ArgError(3, "invalid justOne option")
						return 0
					}
				}
			}
		}
	}

	ctx, cancel := coll.Client.Context()
	defer cancel()

	var res *mongo.DeleteResult
	var err error
	if justOne {
		res, err = coll.Collection.DeleteOne(ctx, query)
	} else {
		res, err = coll.Collection.DeleteMany(ctx, query)
	}
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(bsonutil.ToLuaValue(L, newRemoveResult(int(res.DeletedCount))))
	return 1
}

func newUpdateResult(res *mongo.UpdateResult) map[string]int {
	return map[string]int{
		"nMatched":  int(res.MatchedCount),
		"nUpserted": int(res.UpsertedCount),
		"nModified": int(res.ModifiedCount),
	}
}

func collectionUpdateMethod(L *lua.LState) int {
	coll := checkCollection(L)

	query := bsonutil.CastBSON(L, 2)
	document := bsonutil.CastBSON(L, 3)
	opts := &options.UpdateOptions{}

	var multi bool
	options := bsonutil.ToBSON(L, 3)
	if options != nil {
		if v, ok := options.(map[string]interface{}); ok {
			if v2, ok2 := v["multi"]; ok2 {
				if multiVal, ok3 := v2.(bool); ok3 {
					multi = multiVal
				} else {
					L.ArgError(3, "invalid multi option")
					return 0
				}
			}
			if v2, ok2 := v["upsert"]; ok2 {
				if upsertVal, ok3 := v2.(bool); ok3 {
					opts.SetUpsert(upsertVal)
				} else {
					L.ArgError(3, "invalid upsert option")
					return 0
				}
			}
		}
	}

	ctx, cancel := coll.Client.Context()
	defer cancel()

	var res *mongo.UpdateResult
	var err error
	if multi {
		res, err = coll.Collection.UpdateMany(ctx, query, document, opts)
	} else {
		res, err = coll.Collection.UpdateOne(ctx, query, document, opts)
	}
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(bsonutil.ToLuaValue(L, newUpdateResult(res)))
	return 1
}
