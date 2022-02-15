package search

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"log"
)



// Search search products
// @Summary search products
// @Description search products
// @Tags search
// @Accept json
// @Produce json
// @Param data body search true "search term"
// @Success 200 {object} searchResponse{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /search [post]
func Search(c *fiber.Ctx) error {
	s := new(search)
	if err := utils.ParseBodyAndValidate(c, s); err != nil {
		return c.JSON(err)
	}
	query := "let a=(for i in productSearch Search LIKE(i.title, \"%"+s.SearchString+"%\") sort i.seen desc limit 5 return {id:i._id,title:i.title})" +
		"let b = (for j in categories filter LIKE(j.name, \"%"+s.SearchString+"%\") limit 5 return {id:j._id,name:j.name,url:j.url})" +
		"let ms=(for k in mostSearch sort k.searchCount limit 16 return k)\nreturn {products:a,categories:b,mostSearch:ms}"

	log.Println(query)
	res:=database.ExecuteGetQuery(query)
	return c.JSON(res[0])


}


// SetMostSearch create most search data
// @Summary create most search data
// @Description create most search data
// @Tags search
// @Accept json
// @Produce json
// @Param data body mostSearch true "data"
// @Success 200 {object} string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /search/ms [post]
func SetMostSearch(c *fiber.Ctx)error{
	ms := new(mostSearch)
	if err := utils.ParseBodyAndValidate(c, ms); err != nil {
		return c.JSON(err)
	}
	query := fmt.Sprintf("for i in mostSearch filter i.id==\"%v\" return i")
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()

	var doc mostSearchOut
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil && !errors.Is(err,driver.NoMoreDocumentsError{}) {
		return c.JSON(err)
	}

	if errors.Is(err,driver.NoMoreDocumentsError{}) {
		msCol:=database.GetCollection("mostSearch")
		ms.SearchCount=0
		_,err=msCol.CreateDocument(context.Background(),ms)
		return c.JSON(fiber.Map{"status":"ok"})

	}

	query = fmt.Sprintf("for i in mostSearch filter i.id==\"%v\" update i with { searchCount:i.searchCount + 1 } in mostSearch",ms.Id)
	database.ExecuteGetQuery(query)

	return c.JSON(fiber.Map{"status":"ok"})

}