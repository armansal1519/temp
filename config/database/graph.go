package database

import (
	"github.com/arangodb/go-driver"
	"log"
)

type MyEdgeObject struct {
	From string `json:"_from"`
	To   string `json:"_to"`
}

func getGraph(name string) driver.Graph {
	db := GetDB()
	graph, err := db.Graph(nil, name)
	if err != nil {
		log.Fatalf("Failed to create graph: %v", err)
	}
	return graph
}

func GetEdgeCollection(graphName string, collectionName string) driver.Collection {
	graph := getGraph(graphName)

	edgeCollection, _, err := graph.EdgeCollection(nil, collectionName)
	if err != nil {
		log.Fatalf("Failed to select edge collection: %v", err)
	}
	return edgeCollection
}
