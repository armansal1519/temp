package database

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
)

func GetAll(collectionName string, offset int, limit int) []interface{} {
	if limit > 50 {
		panic("large limit in getAll")
	}
	db := GetDB()
	ctx := context.Background()
	query := fmt.Sprintf("FOR d IN %v \nLIMIT %v,%v RETURN d", collectionName, offset, limit)
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []interface{}
	for {
		var doc interface{}
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return data
}

func GetByKet(key string, collectionName string) interface{} {
	db := GetDB()
	ctx := context.Background()
	query := fmt.Sprintf("FOR d IN %v filter d._key==\"%v\" LIMIT 1 RETURN d", collectionName, key)
	fmt.Println(query)
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		// handle error
	}
	defer cursor.Close()
	var doc interface{}

	meta, err := cursor.ReadDocument(ctx, &doc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Got doc with key '%s' from query\n", meta.Key)

	return doc
}

func CreateDocument(data interface{}, collectionName string) interface{} {
	var doc interface{}
	ctx := context.Background()
	RNctx := driver.WithReturnNew(ctx, &doc)
	col := GetCollection(collectionName)
	_, err := col.CreateDocument(RNctx, data)
	if err != nil {
		panic(fmt.Sprintf("error creating document %v", err))
	}
	return doc

}

func UpdateDocument(key string, data interface{}, collectionName string) driver.DocumentMeta {
	ctx := context.Background()
	col := GetCollection(collectionName)
	meta, err := col.UpdateDocument(ctx, key, data)
	if err != nil {
		panic(fmt.Sprintf("error updating document %v with %v \n%v", key, data, err))
	}
	return meta
}
