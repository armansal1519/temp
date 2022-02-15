package homepage

import (
	"bamachoub-backend-go-v1/config/database"
	"bamachoub-backend-go-v1/utils"
	"context"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func setBaseData(c *fiber.Ctx) error {
	m := new(homepageBase)
	if err := utils.ParseBodyAndValidate(c, m); err != nil {
		return c.JSON(err)
	}
	queryMap := make(map[int]saveQuery, 0)
	//banners validation
	for _, v := range m.Banners {
		_, ok := queryMap[v.Position]

		if ok {
			fmt.Println(queryMap)
			return c.Status(400).SendString(fmt.Sprintf("repetetive position : %v", v.Position))
		}

		if v.NumberOfBanners != len(v.Data) {
			return c.Status(400).SendString("number of banners does not match data lenght")
		}

		queryMap[v.Position] = saveQuery{
			Title: v.Title,
			Query: createBannerQuery(v),
			Type:  "banner",
		}
	}

	//product slider validation
	for _, v := range m.ProductSlider {
		_, ok := queryMap[v.Position]

		if ok {
			return c.Status(400).SendString(fmt.Sprintf("repetetive position : %v", v.Position))
		}

		if v.Sort != "seenNumber" && v.Sort != "buyNumber" && v.Sort != "createdAt" {
			return c.Status(400).SendString("sort in the product slider only can be seenNumber or buyNumber or createdAt")
		}
	}
	db := database.GetDB()

	for _, v := range m.ProductSlider {
		flag, err := db.CollectionExists(context.Background(), v.CategoryName)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		if !flag {
			return c.Status(400).SendString(fmt.Sprintf("category does not exist : %v", v.CategoryName))
		}

		queryMap[v.Position] = saveQuery{Query: createProductQuery(v), Title: v.Title, Type: "product"}
	}

	//category slider validation
	for _, v := range m.CategorySlider {
		_, ok := queryMap[v.Position]

		if ok {
			return c.Status(400).SendString(fmt.Sprintf("repetetive position : %v", v.Position))
		}
	}

	for _, v := range m.CategorySlider {
		flag, err := db.CollectionExists(context.Background(), v.CategoryName)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		if !flag {
			return c.Status(400).SendString(fmt.Sprintf("category does not exist : %v", v.CategoryName))
		}

		queryMap[v.Position] = saveQuery{
			Title: v.Title,
			Type:  "category",
			Query: CreateCategoryQuery(v),
		}

	}

	//brand slider validation
	for _, v := range m.BrandSlider {
		_, ok := queryMap[v.Position]

		if ok {
			return c.Status(400).SendString(fmt.Sprintf("repetetive position : %v", v.Position))
		}
	}

	for _, v := range m.BrandSlider {
		if v.CategoryName == "all" {
			continue
		}

		flag, err := db.CollectionExists(context.Background(), v.CategoryName)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		if !flag {
			return c.Status(400).SendString(fmt.Sprintf("category does not exist : %v", v.CategoryName))
		}
		q := ""
		if v.CategoryName != "all" {
			q = fmt.Sprintf("for i in categories filter i.url==\"%v\" for v in 1..1 inbound i graph \"brandCategory\" sort v.seen limit 8 return v", v.CategoryName)
		} else {
			q = "for i in brands sort i.seen limit 8 return i "
		}

		queryMap[v.Position] = saveQuery{
			Title: v.Title,
			Type:  "brand",
			Query: q,
		}

	}

	queryMap[0] = saveQuery{
		Query: createCarousel(m.Carousel),
		Type:  "carousel",
	}

	saveQueryArr := make([]saveQuery, 0)

	for i := 0; i < len(queryMap); i++ {
		saveQueryArr = append(saveQueryArr, queryMap[i])
	}

	//TODO proper haraji
	saveQueryArr = insert(saveQueryArr, 1, saveQuery{
		Title: "حراجی",
		Type:  "discount",
		Query: "for i in productSearch limit 15 return {discount:5,product:i}",
	})

	saveQueryArr = append(saveQueryArr, saveQuery{
		Title: "آخرین مطالب مجله",
		Type:  "blogContent",
		Query: fmt.Sprintf("return %v", m.BlogContent.Show),
	})

	saveQueryArr = append(saveQueryArr, saveQuery{
		Title: "از جدیدترین تخفیفات و جشنواره ها باخبر شوید",
		Type:  "email",
		Query: fmt.Sprintf("return %v", m.Email.Show),
	})

	f := finalHomePageData{
		Key:  "1",
		Data: saveQueryArr,
	}

	constCol := database.GetCollection("const")

	flag, err := constCol.DocumentExists(context.Background(), "1")
	if err != nil {
		return c.Status(500).JSON(err)
	}
	if flag {
		_, err := constCol.ReplaceDocument(context.Background(), "1", f)
		if err != nil {
			return c.Status(500).JSON(err)
		}
	} else {
		_, err = constCol.CreateDocument(context.Background(), f)
		if err != nil {
			return c.Status(500).JSON(err)
		}
	}

	return c.JSON(f)

}

// getHomePage get homepage data
// @Summary return homepage data
// @Description return homepage data
// @Tags homepage
// @Accept json
// @Produce json
// @Param offset query int    true  "Offset"
// @Param limit  query int    true  "limit"
// @Success 200 {object} string{}
// @Failure 404 {object} string{}
// @Router /homepage [get]
func getHomePage(c *fiber.Ctx) error {

	offset := c.Query("offset")
	limit := c.Query("limit")

	if offset == "" || limit == "" {
		return c.Status(400).SendString("Offset and Limit must have a value")
	}

	constCol := database.GetCollection("const")
	var d finalHomePageData
	_, err := constCol.ReadDocument(context.Background(), "1", &d)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	o, err := strconv.Atoi(offset)
	l, err := strconv.Atoi(limit)
	if err != nil {
		return c.Status(400).SendString("error parsing offset or limit")
	}

	query1 := ""
	query2 := ""

	if o >= len(d.Data) {
		return c.JSON(nil)
	}
	var sq []saveQuery

	if o+l > len(d.Data)-1 {
		sq = d.Data[o:]
	} else {
		sq = d.Data[o:l]
	}

	for i, v := range sq {
		query1 += fmt.Sprintf("let a%v = ( %v ) \n", i, v.Query)

		query2 += fmt.Sprintf("{type:\"%v\",title:\"%v\",data:a%v}", v.Type, v.Title, i)
		if i < len(sq)-1 {
			query2 += " , "
		}
	}

	finalQ := fmt.Sprintf("%v \n return [%v]", query1, query2)

	res := database.ExecuteGetQuery(finalQ)

	return c.JSON(res[0])

}
