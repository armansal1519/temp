package similarityGraph

import (
	"bamachoub-backend-go-v1/app/admin"
	"bamachoub-backend-go-v1/app/products"
	"bamachoub-backend-go-v1/app/users"
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

// createSimilarityNode create similarity node
// @Summary create similarity node
// @Description create similarity node
// @Tags product similarity
// @Accept json
// @Produce json
// @Param data body SimilarityNodeRequest true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /p-similarity/node [post]
func createSimilarityNode(c *fiber.Ctx) error {
	fmt.Print("1111111111")
	snq := new(SimilarityNodeRequest)
	if err := utils.ParseBodyAndValidate(c, snq); err != nil {
		return c.JSON(err)
	}
	userKey := c.Locals("userKey").(string)
	u, err := users.GetUserByKey(userKey)
	if err != nil {

		return c.JSON(err)
	}

	sn := similarityNode{
		Title:            snq.Title,
		ProductsKeyArray: snq.ProductsKeyArray,
		ImageUrls:        []string{},
		Tags:             snq.Tags,
		UserKey:          userKey,
		Description:      snq.Description,
		Color:            snq.Color,
		Pattern:          snq.Pattern,
		UserMade:         false,
		Public:           snq.Public,
		IsCollection:     snq.IsCollection,
		Status:           "not",
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        0,
		CreatedBy:        fmt.Sprintf("%v %v", u.FirstName, u.LastName),
		UpdatedBy:        "",
	}
	if len(snq.ProductsKeyArray) > 0 {
		iu, err := getImageArrFromProductsKey(snq.ProductsKeyArray)
		if err != nil {
			return c.JSON(err)
		}
		sn.ImageUrls = iu

	}

	col := database.GetCollection("productSimilarityNode")
	meta, err := col.CreateDocument(context.Background(), sn)
	if err != nil {
		return c.JSON(err)
	}

	if len(snq.ProductsKeyArray) > 0 {
		err := createSimilarityInnerEdge(snq.ProductsKeyArray, meta.Key)
		if err != nil {
			return c.JSON(err)
		}
	}
	return c.JSON(meta)
}

// createSimilarityNodeByAdmin create similarity node by admin
// @Summary create similarity node by admin
// @Description create similarity node by admin
// @Tags product similarity
// @Accept json
// @Produce json
// @Param data body SimilarityNodeRequest true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /p-similarity/node/admin [post]
func createSimilarityNodeByAdmin(c *fiber.Ctx) error {
	snq := new(SimilarityNodeRequest)

	if err := utils.ParseBodyAndValidate(c, snq); err != nil {
		return c.JSON(err)
	}
	adminKey := c.Locals("adminKey").(string)
	adminCol := database.GetCollection("admin")
	var a admin.AdminOut
	_, err := adminCol.ReadDocument(context.Background(), adminKey, &a)
	if err != nil {
		return c.JSON(err)
	}

	sn := similarityNode{
		Title:            snq.Title,
		ProductsKeyArray: snq.ProductsKeyArray,
		ImageUrls:        []string{},
		Tags:             snq.Tags,
		UserKey:          adminKey,
		Description:      snq.Description,
		Color:            snq.Color,
		Pattern:          snq.Pattern,
		UserMade:         false,
		Public:           snq.Public,
		IsCollection:     snq.IsCollection,
		Status:           "not",
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        0,
		CreatedBy:        fmt.Sprintf("%v %v", a.FirstName, a.LastName),
		UpdatedBy:        "",
	}
	if len(snq.ProductsKeyArray) > 0 {
		iu, err := getImageArrFromProductsKey(snq.ProductsKeyArray)
		if err != nil {
			return c.JSON(err)
		}
		sn.ImageUrls = iu

	}

	col := database.GetCollection("productSimilarityNode")
	meta, err := col.CreateDocument(context.Background(), sn)
	if err != nil {
		return c.JSON(err)
	}

	if len(snq.ProductsKeyArray) > 0 {
		err := createSimilarityInnerEdge(snq.ProductsKeyArray, meta.Key)
		if err != nil {
			return c.JSON(err)
		}
	}
	return c.JSON(meta)
}

// createSimilarityEdge create similarity edge
// @Summary create similarity edge
// @Description create similarity edge
// @Tags product similarity
// @Accept json
// @Produce json
// @Param similarityEdge body similarityEdge true "data"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /p-similarity/edge [post]
func createSimilarityEdge(c *fiber.Ctx) error {
	sn := new(similarityEdge)
	if err := utils.ParseBodyAndValidate(c, sn); err != nil {
		return c.JSON(err)
	}
	if sn.Type != "sim" && sn.Type != "con" {
		return c.Status(400).SendString("type only can be sim or con but type : " + sn.Type)
	}
	edgeCol := database.GetCollection("productSimilarityEdge")
	sn.Key = strings.Split(sn.From, "/")[1] + strings.Split(sn.To, "/")[1]
	meta, err := edgeCol.CreateDocument(context.Background(), sn)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// AddOrRemoveProductToNode add or remove product to a node
// @Summary  add or remove product to a node
// @Description add or remove product to a node
// @Tags product similarity
// @Accept json
// @Produce json
// @Param op path string true "op"
// @Param productKey path string true "productKey"
// @Param nodeKey path string true "nodeKey"
// @Success 200 {object} similarityNode
// @Failure 404 {object} string{}
// @Router /p-similarity/{op}/{productKey}/{nodeKey} [post]
func AddOrRemoveProductToNode(c *fiber.Ctx) error {
	op := c.Params("op")
	if op != "add" && op != "remove" {
		return c.Status(fiber.StatusBadRequest).JSON("op only can be add or remove")
	}
	productKey := c.Params("productKey")
	nodeKey := c.Params("nodeKey")
	simCol := database.GetCollection("productSimilarityInnerEdge")
	if op == "remove" {
		_, err := simCol.RemoveDocument(context.Background(), productKey+nodeKey)
		if err != nil {
			return c.JSON(err)
		}
		q := fmt.Sprintf("for i in productSimilarityNode  \nfilter i._key==\"%v\" \nlet iu=( for s in sheet filter s._key==\"%v\" return s.imageArr[0])\nupdate i with {productsKeyArray: REMOVE_VALUE( i.productsKeyArray, \"%v\" ),imageUrl : REMOVE_VALUE(i.imageUrl,iu[0],1)} in productSimilarityNode\nreturn NEW", nodeKey, productKey, productKey)
		return c.JSON(database.ExecuteGetQuery(q))
	}
	var node similarityNode
	nodeCol := database.GetCollection("productSimilarityNode")
	_, err := nodeCol.ReadDocument(context.Background(), nodeKey, &node)
	if err != nil {
		return c.JSON(err)
	}
	for _, key := range node.ProductsKeyArray {
		if productKey == key {
			return c.Status(400).JSON("productKey is already in collection")
		}
	}
	add := addProductToNode{}
	if len(node.ImageUrls) < 4 {
		img, err := getImageArrFromProductsKey([]string{productKey})
		if err != nil {
			return c.JSON(err)
		}
		add.ImageUrls = append(node.ImageUrls, img[0])
		add.ProductsKeyArray = append(node.ProductsKeyArray, productKey)
	} else {
		add.ImageUrls = node.ImageUrls
		add.ProductsKeyArray = append(node.ProductsKeyArray, productKey)
	}
	err = createSimilarityInnerEdge([]string{productKey}, nodeKey)
	if err != nil {
		return c.JSON(err)
	}
	ctx := driver.WithReturnNew(context.Background(), &node)
	_, err = nodeCol.UpdateDocument(ctx, nodeKey, add)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(node)
}

// getAllSimilarityNodes get all similarity nodes
// @Summary return all similarity nodes
// @Description return all similarity nodes
// @Tags product similarity
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Success 200 {object} []similarityNode{}
// @Failure 404 {object} string{}
// @Router /p-similarity [get]
func getAllSimilarityNodes(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	query := fmt.Sprintf("for i in productSimilarityNode limit %v,%v return i", offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}

// getSimilarNodeToOneNodeByNodeKey get all similar nodes
// @Summary return all similar nodes
// @Description return all similar nodes
// @Tags product similarity
// @Accept json
// @Produce json
// @Param type query string false  "similarity type"
// @Param key path string true "key"
// @Success 200 {object} []similarityNodeOut{}
// @Failure 404 {object} string{}
// @Router /p-similarity/graph-sim/{key} [get]
func getSimilarNodeToOneNodeByNodeKey(c *fiber.Ctx) error {
	nodeKey := c.Params("key")
	edgeType := c.Query("type")

	typeString := ""
	if edgeType == "con" {
		typeString = "e.type==\"con\""
	} else if edgeType == "sim" {
		typeString = "e.type==\"sim\""
	} else {
		typeString = "e.type==\"con\" || e.type==\"sim\""
	}
	query := fmt.Sprintf("for n in productSimilarityNode filter n._key==\"%v\" for v,e in 0..1 any n graph \"sheetSimilarityGraph\" filter e.coreEdge==true\nfilter %v sort v.seen return v", nodeKey, typeString)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return c.JSON(err)
	}
	defer cursor.Close()
	var data []similarityNodeOut
	for {
		var doc similarityNodeOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return c.JSON(err)
		}
		data = append(data, doc)
	}
	return c.JSON(data)
}

// getSimilarProductsByProductKey get all similar products
// @Summary return all similar products
// @Description return all similar products
// @Tags product similarity
// @Accept json
// @Produce json
// @Param offset query int    false  "Offset"
// @Param limit  query int    false  "limit"
// @Param key path string true "key"
// @Success 200 {object} []products.Product{}
// @Failure 404 {object} string{}
// @Router /p-similarity/{key} [get]
func getSimilarProductsByProductKey(productKey string, offset string, limit string) (*[]products.Product, error) {
	limitString := ""
	if offset != "" && limit != "" {
		limitString = fmt.Sprintf(" limit %v,%v ", offset, limit)
	}
	query := fmt.Sprintf("for i in productSimilarityInnerEdge filter i._from==\"sheet/%v\"\nfor j in productSimilarityInnerEdge\nfilter i._to==j._to\nfor k in sheet\nfilter j._from==k._id %v return k", productKey, limitString)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()
	var data []products.Product
	for {
		var doc products.Product
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}
	return &data, nil

}

// getOneNodeByKey get one node with products
// @Summary return one node with products
// @Description return one node with products
// @Tags product similarity
// @Accept json
// @Produce json
// @Param offset query int    false  "Offset"
// @Param limit  query int    false  "limit"
// @Param key path string true "key"
// @Success 200 {object} []products.Product{}
// @Failure 404 {object} string{}
// @Router /p-similarity/node/{key} [get]
func getOneNodeByKey(c *fiber.Ctx) error {
	nodeKey := c.Params("nodeKey")
	offset := c.Query("offset")
	limit := c.Query("limit")
	if offset == "" || limit == "" {
		return c.Status(4040).JSON("offset or limit is empty")
	}
	query := fmt.Sprintf("let data=(for i in productSimilarityNode  filter i._key==\"%v\" return i)\nlet products =(for j in sheet filter j._key in data[0].productsKeyArray limit %v,%v return j)\nreturn {node:data,products:products}\n\n", nodeKey, offset, limit)
	return c.JSON(database.ExecuteGetQuery(query))
}

// updateSimilarityNode update similarity node
// @Summary update similarity node
// @Description update similarity node
// @Tags product similarity
// @Accept json
// @Produce json
// @Param data body similarityNode true "data"
// @Param key path string true "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /p-similarity/{key} [put]
func updateSimilarityNode(c *fiber.Ctx) error {
	snq := new(similarityNode)
	key := c.Params("key")
	if err := utils.ParseBodyAndValidate(c, snq); err != nil {
		return c.JSON(err)
	}

	isAdmin := c.Locals("isAdmin").(bool)
	if !isAdmin {
		userKey := c.Locals("userKey").(string)
		u, err := users.GetUserByKey(userKey)
		if err != nil {

			return c.JSON(err)
		}
		if userKey != u.Key {
			return c.Status(403).SendString("not user collection")
		}
		snq.UpdatedAt = time.Now().Unix()
		snq.UpdatedBy = u.FirstName + " " + u.LastName
		snq.Status = "not"

		col := database.GetCollection("productSimilarityNode")
		meta, err := col.UpdateDocument(context.Background(), key, snq)
		if err != nil {
			return c.JSON(err)
		}

		if len(snq.ProductsKeyArray) > 0 {
			err := createSimilarityInnerEdge(snq.ProductsKeyArray, meta.Key)
			if err != nil {
				return c.JSON(err)
			}
		}
		return c.JSON(meta)
	}

	adminKey := c.Locals("adminKey").(string)
	adminCol := database.GetCollection("admin")
	var a admin.AdminOut
	_, err := adminCol.ReadDocument(context.Background(), adminKey, &a)
	if err != nil {
		return c.JSON(err)
	}

	snq.UpdatedAt = time.Now().Unix()
	snq.UpdatedBy = a.FirstName + " " + a.LastName
	snq.Status = "not"

	col := database.GetCollection("productSimilarityNode")
	meta, err := col.UpdateDocument(context.Background(), key, snq)
	if err != nil {
		return c.JSON(err)
	}

	if len(snq.ProductsKeyArray) > 0 {
		err := createSimilarityInnerEdge(snq.ProductsKeyArray, meta.Key)
		if err != nil {
			return c.JSON(err)
		}
	}
	return c.JSON(meta)
}

// removeNode delete similarity node
// @Summary delete similarity node
// @Description delete similarity node
// @Tags product similarity
// @Accept json
// @Produce json
// @Param key path string true "key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /p-similarity/{key} [delete]
func removeNode(c *fiber.Ctx) error {
	nodeKey := c.Params("key")
	isAdmin := c.Locals("isAdmin").(bool)
	if isAdmin {
		query := fmt.Sprintf("for i in productSimilarityNode filter i._key==\"%v\"  remove i in productSimilarityNode\nfor j in productSimilarityInnerEdge filter j._to ==i._id remove j in productSimilarityInnerEdge", nodeKey)
		database.ExecuteGetQuery(query)
		return c.Status(204).JSON("node deleted")
	}

	userKey := c.Locals("userKey").(string)
	query := fmt.Sprintf("for i in productSimilarityNode filter i._key==\"%v\" && i.userKey==\"%v\" remove i in productSimilarityNode\nfor j in productSimilarityInnerEdge filter j._to ==i._id remove j in productSimilarityInnerEdge", nodeKey, userKey)
	database.ExecuteGetQuery(query)
	return c.Status(204).JSON("node deleted")
}
