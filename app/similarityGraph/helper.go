package similarityGraph

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"log"
)

func getImageArrFromProductsKey(pk []string) ([]string, error) {
	keyStr := "["
	for i, key := range pk {
		keyStr += fmt.Sprintf("\"%v\"", key)
		if i < len(pk)-1 {
			keyStr += " , "
		}
	}
	keyStr += "] "

	query := fmt.Sprintf("for i in sheet filter i._key in %v return i.imageArr[0]", keyStr)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []string
	for {
		var doc string
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return data, nil
}

func createSimilarityInnerEdge(productsKeyArray []string, similarityNodeKey string) error {
	edgeArr := make([]productEdge, 0)
	for _, key := range productsKeyArray {
		temp := productEdge{
			Key:   key + similarityNodeKey,
			From:  fmt.Sprintf("sheet/%v", key),
			To:    fmt.Sprintf("productSimilarityNode/%v", similarityNodeKey),
			Score: 1,
		}
		edgeArr = append(edgeArr, temp)
	}

	innerEdgeCol := database.GetCollection("productSimilarityInnerEdge")
	_, errorArr, err := innerEdgeCol.CreateDocuments(context.Background(), edgeArr)
	if err != nil {
		log.Println(1)
		return fmt.Errorf("%v", errorArr)
	}
	return nil

}
