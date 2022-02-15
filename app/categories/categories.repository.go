package categories

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"log"

	"github.com/arangodb/go-driver"
)

func createEdge(from string, to string) driver.DocumentMeta {

	edgeCollection := database.GetEdgeCollection("categoryGraph", "categoriesEdge")
	edge := createNewCategoryEdge{From: from, To: to}
	meta, err := edgeCollection.CreateDocument(context.TODO(), edge)
	if err != nil {
		log.Fatalf("Failed to create edge document: %v", err)
	}
	return meta
}

// func getPrevNodeCategory(key string, collectionName string) Category {
// 	ctx := context.Background()
// 	col := database.GetCollection(collectionName)
// 	var doc Category
// 	_, err := col.ReadDocument(ctx, key, &doc)
// 	if err != nil {
// 		panic(fmt.Sprintf("error in getPrevNodeCategory %v", err))
// 	}
// 	return doc

// }
