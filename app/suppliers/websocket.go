package suppliers

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/antoniodipinto/ikisocket"
	"github.com/arangodb/go-driver"
	"log"
	"strings"
)

func newConn(kws *ikisocket.Websocket) {

	col := database.GetCollection("onlineSuppliers")
	supplierId := kws.Locals("supplierId").(string)
	supplierKey := strings.Split(supplierId, "/")[1]
	con := onlineSuppliers{
		Key:  supplierKey,
		Uuid: kws.UUID,
	}
	meta, err := col.CreateDocument(context.Background(), con)
	if err != nil {
		log.Println(err)
	}
	log.Println(meta)

	// Retrieve the user id from endpoint
	//userId := kws.Params("id")

	// Add the connection to the list of the connected clients
	// The UUID is generated randomly and is the key that allow
	// ikisocket to manage Emit/EmitTo/Broadcast
	//clients[userId] = kws.UUID

	// Every websocket connection has an optional session key => value storage
	kws.SetAttribute("supplier_key", supplierKey)

	//Broadcast to all the connected users the newcomer
	kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", supplierKey, kws.UUID)), true)
	//Write welcome message
	kws.Emit([]byte(fmt.Sprintf("Hello user: %s with UUID: %s", supplierKey, kws.UUID)))
}

func remove(supplierKey string) error {
	col := database.GetCollection("onlineSuppliers")
	_, err := col.RemoveDocument(context.Background(), supplierKey)
	if err != nil {
		return err
	}
	return nil
}

func NewEstelam(keys []string) {
	log.Println(111, keys)
	keyStr := "["
	for i, key := range keys {
		keyStr += fmt.Sprintf("\"%v\"", key)
		if i < len(keys)-1 {
			keyStr += " , "
		}
	}
	keyStr += "] "

	query := fmt.Sprintf("for i in onlineSuppliers filter i._key in %v return i.uuid", keyStr)
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

	log.Println(222, data)
	if err != nil {
		log.Println(err)
	}
	ikisocket.EmitToList(data, []byte("new estelam"))
	//ikisocket.EmitTo(data[0],[]byte("new estelam"))
}

type onlineUuid struct {
	Uuid string `json:"uuid"`
}
