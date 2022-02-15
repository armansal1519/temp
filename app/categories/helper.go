package categories

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"log"

	driver "github.com/arangodb/go-driver"
)

func createGraphAndEdge(name string) error {
	db := database.GetDB()

	// define the edgeCollection to store the edges
	var edgeDefinition driver.EdgeDefinition
	edgeDefinition.Collection = fmt.Sprintf("categories_%v", name)
	// define a set of collections where an edge is going out...
	edgeDefinition.From = []string{"categories"}

	// repeat this for the collections where an edge is going into
	edgeDefinition.To = []string{name}

	var options driver.CreateGraphOptions
	options.OrphanVertexCollections = []string{"categories", name}
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}

	_, err := db.CreateGraph(context.TODO(), fmt.Sprintf("categories_%v_graph", name), &options)
	if err != nil {
		log.Fatalf("Failed to create graph: %v", err)
	}
	return nil

}

func  createSupplierToProductEdgeAndAddToGraph(name string) error {
	db := database.GetDB()
	// define the edgeCollection to store the edges

	var edgeDefinition1 driver.EdgeDefinition
	var edgeDefinition2 driver.EdgeDefinition
	edgeDefinition1.Collection = fmt.Sprintf("supplier_%v_price", name)
	edgeDefinition2.Collection = fmt.Sprintf("supplier_%v_estelam", name)
	// define a set of collections where an edge is going out...
	edgeDefinition1.From = []string{"suppliers"}
	edgeDefinition2.From = []string{"suppliers"}

	// repeat this for the collections where an edge is going into
	edgeDefinition1.To = []string{name}
	edgeDefinition2.To = []string{name}

	var options driver.CreateGraphOptions
	options.OrphanVertexCollections = []string{"suppliers", name}
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition1, edgeDefinition2}

	_, err := db.CreateGraph(context.TODO(), fmt.Sprintf("supplier_%v_graph", name), &options)
	if err != nil {
		log.Fatalf("Failed to create graph: %v", err)
	}
	return nil
}
