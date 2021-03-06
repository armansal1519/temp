package suppliers

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
)

func GetSuppliers(c *fiber.Ctx) error {
	return c.JSON(
		database.GetAll("suppliers", 0, 32))
}

func CreateSupplier(c *fiber.Ctx) error {
	s := new(SupplierIn)

	if err := c.BodyParser(s); err != nil {
		return err
	}
	errors := utils.Validate(s)
	if errors != nil {
		c.JSON(errors)
		return nil
	}
	s.WalletAmount = 0
	suppliersCollection := database.GetCollection("suppliers")
	var newSupplier interface{}
	ctx := driver.WithReturnNew(context.Background(), &newSupplier)

	_, err := suppliersCollection.CreateDocument(ctx, s)
	if err != nil {
		panic(fmt.Sprintf("error creating Supplier :%v", err))
	}
	return c.JSON(newSupplier)

}

func GetSupplierByEmployeeKey(key string) (Supplier, error) {
	q := fmt.Sprintf("for se in supplierEmployee filter se._key==\"%v\" for s in suppliers filter s._key==se.supplierKey return s", key)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return Supplier{}, fmt.Errorf("error while excuting query: %v \n error:%v", q, err)
	}
	defer cursor.Close()

	var doc Supplier
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return Supplier{}, err
	}

	return doc, nil
}

func GetSupplierByKey(key string) (*Supplier, error) {
	var s Supplier
	sCol := database.GetCollection("suppliers")
	_, err := sCol.ReadDocument(context.Background(), key, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// getFavBySupplierKey get Favorite product to Supplier
// @Summary get Favorite product to Supplier
// @Description get Favorite product to Supplier
// @Tags  Supplier
// @Accept json
// @Produce json
// @Param   categoryUrl      path   string     true  "categoryUrl"
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []Fav{}
// @Failure 400 {object} string
// @Router /suppliers/fav/{categoryUrl} [get]
func getFavBySupplierKey(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")

	if limit == "" || offset == "" {
		return c.Status(400).SendString("limit and offset are required")
	}
	key := c.Locals("supplierId").(string)
	categoryUrl := c.Params("categoryUrl")

	q := fmt.Sprintf("for i in fav\nfilter i.supplierKey==\"%v\" && i.categoryUrl==\"%v\"\nlimit %v,%v\nreturn i", key, categoryUrl, offset, limit)
	log.Println(q)
	return c.JSON(database.ExecuteGetQuery(q))
}

// getAllFavBySupplierKey get Favorite product to Supplier
// @Summary get Favorite product to Supplier
// @Description get Favorite product to Supplier
// @Tags  Supplier
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param   categoryUrl      path   string     true  "categoryUrl"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} []Fav{}
// @Failure 400 {object} string
// @Router /suppliers/fav-product/{categoryUrl} [get]
func getAllFavBySupplierKey(c *fiber.Ctx) error {
	offset := c.Query("offset")
	limit := c.Query("limit")
	categoryUrl := c.Params("categoryUrl")

	if limit == "" || offset == "" {
		return c.Status(400).SendString("limit and offset are required")
	}
	key := c.Locals("supplierId").(string)
	log.Println(key)

	q := fmt.Sprintf("for f in fav filter f.supplierKey==\"%v\" for i in %v filter i._key==f.productKey limit %v,%v return i", key, categoryUrl, offset, limit)
	ql := fmt.Sprintf("let data=(for f in fav filter f.supplierKey==\"%v\" for i in %v filter i._key==f.productKey  return i) return length(data)", key, categoryUrl)
	res := database.ExecuteGetQuery(ql)
	log.Println(q)
	return c.JSON(fiber.Map{
		"data":   database.ExecuteGetQuery(q),
		"length": res[0],
	})
}

func getAllFavProductIdBysupplierKey(c *fiber.Ctx) error {
	key := c.Locals("supplierId").(string)
	q := fmt.Sprintf("for i in fav filter i.supplierKey==\"%v\"  return concat(i.categoryUrl,\"/\",i.productKey)", key)
	fmt.Println(q)
	res := database.ExecuteGetQuery(q)
	return c.JSON(res)
}

// addFavorite add Favorite product to Supplier
// @Summary add Favorite product to Supplier
// @Description add Favorite product to Supplier
// @Tags  Supplier
// @Accept json
// @Produce json
// @Param   key      path   string     true  "key of product you want to add"
// @Param   categoryUrl      path   string     true  "categoryUrl"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} Fav{}
// @Failure 400 {object} string
// @Router /suppliers/add-fav/{categoryUrl}/{key} [post]
func addFavorite(c *fiber.Ctx) error {
	supplierId := c.Locals("supplierId").(string)
	key := c.Params("key")
	cu := c.Params("categoryUrl")
	favCol := database.GetCollection("fav")
	s, err := GetSupplierByKey(supplierId)
	if err != nil {
		return c.JSON(err)
	}
	flag := false
	for _, fav := range s.CategoriesToSale {
		if fav == cu {
			flag = true
		}
	}
	if !flag {
		return c.Status(400).SendString("category url not found or not acceptable")
	}
	fav := Fav{
		Key:         s.Key + key,
		CategoryUrl: cu,
		ProductKey:  key,
		SupplierKey: s.Key,
		Status:      "ok",
	}
	_, err = favCol.CreateDocument(context.Background(), fav)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(fav)

}

// addSupplierToUpdatePool add Supplier to update pool
// @Summary add Supplier to update pool
// @Description add Supplier to update pool
// @Tags  Supplier
// @Accept json
// @Produce json
// @Param data body updateSupplier true "data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string{}
// @Failure 400 {object} string{}
// @Router /suppliers/add-update-pool [put]
func addSupplierToUpdatePool(c *fiber.Ctx) error {
	data := new(updateSupplier)
	if err := utils.ParseBodyAndValidate(c, data); err != nil {
		return c.JSON(err)
	}
	data.CreateAt = time.Now().Unix()
	data.SupplierKey = c.Locals("supplierKey").(string)
	updatePoolCol := database.GetCollection("supplierUpdatePool")
	meta, err := updatePoolCol.CreateDocument(context.Background(), data)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}

// deleteFavorite remove Favorite product to Supplier
// @Summary remove Favorite product to Supplier
// @Description remove Favorite product to Supplier
// @Tags  Supplier
// @Accept json
// @Produce json
// @Param   key      path   string     true  "fav key"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /suppliers/remove-fav/{key} [post]
func deleteFavorite(c *fiber.Ctx) error {
	key := c.Params("key")
	favCol := database.GetCollection("fav")
	meta, err := favCol.RemoveDocument(context.Background(), key)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(meta)
}
