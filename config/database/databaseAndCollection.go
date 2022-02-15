package database

import (
	"bamachoub-backend-go-v1/config"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func GetDB() driver.Database {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{config.DB_HOST},
	})
	if err != nil {
		panic("error connecting arangodb")
	}
	//log.Println(conn)
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "rootpassword"),
	})
	if err != nil {
		panic("error while creating new arango client")
	}
	db, err := c.Database(nil, config.DB_NAME)
	if db == nil {
		panic("check your internet")
	}

	return db
}

func GetCollection(name string) driver.Collection {

	db := GetDB()
	ctx := context.Background()
	//found, err := db.CollectionExists(ctx, name)
	//if err != nil {
	//	panic(fmt.Sprintf("error while chacking collection: %v exist",name))
	//}
	col, err := db.Collection(ctx, name)
	if err != nil {
		panic(fmt.Sprintf("error collection: %v does not exist", name))
	}
	return col

	//options := &driver.CreateCollectionOptions{ /* ... */ }
	//col, err := db.CreateCollection(ctx, name, options)
	//if err != nil {
	//	panic(fmt.Sprintf("error while creating collection: %v exist",name))
	//}
	//fmt.Printf("collection: %v was created" ,name)
	//return col

}

func ExecuteGetQuery(query string) []interface{} {
	db := GetDB()
	ctx := context.Background()
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
