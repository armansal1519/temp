package products

import (
	"bamachoub-backend-go-v1/app/categories"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"log"
	"strings"
	"time"
)

// createTheFuckingProduct  create products
// @Summary create products
// @Description create products
// @Tags products
// @Accept json
// @Produce json
// @Param data body productInfo true "data"
// @Success 200 {object} []string{}
// @Failure 404 {object} string{}
// @Router /products/ [post]
func createTheFuckingProduct(productInfo productInfo) (driver.DocumentMeta, error) {
	catCol := database.GetCollection("categories")
	var c categories.CategoryOut
	catMeta, err := catCol.ReadDocument(context.Background(), productInfo.CategoryKey, &c)
	if err != nil {
		return driver.DocumentMeta{}, err
	}
	baseCat, err := categories.GetBaseCategory(catMeta.Key)
	q := fmt.Sprintf("for p in productBody filter p.categoryKey==\"%v\" return p", productInfo.CategoryKey)
	db := database.GetDB()
	cursor, err := db.Query(context.Background(), q, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", q))
	}
	defer cursor.Close()

	var doc ProductBodyOut
	_, err = cursor.ReadDocument(context.Background(), &doc)

	ms, cs := fieldsToSpecs(productInfo.Fields, doc.MainSpecs, doc.CompleteSpec)
	filterString := createFilterString(productInfo.Fields)
	log.Println(filterString)

	pi := Product{
		CategoryKey:  productInfo.CategoryKey,
		CategoryPath: c.GraphPath,
		Title:        titleMakerToTitle(productInfo.Fields, doc.TitleMaker, productInfo.Brand),
		CompleteSpec: cs,
		MainSpecs:    ms,
		VariationObj: v{
			Title:      doc.VariationObj.Title,
			Variations: productInfo.VariationsObj,
		},
		ImageArr:               productInfo.ImageArr,
		Status:                 "ok",
		Description:            productInfo.Description,
		Brand:                  productInfo.Brand,
		Tags:                   productInfo.Tags,
		SpId:                   productInfo.SpId,
		CreatedAt:              time.Now().Unix(),
		CommissionPercent:      c.CommissionPercent,
		CheckCommissionPercent: c.CheckCommissionPercent,
		LowestPrice:            -1,
		FilterArr:              filterString,
	}

	if len(productInfo.ImageArr) <= 0 {
		return driver.DocumentMeta{}, fmt.Errorf("image arr can not be empty")
	}
	productCol := database.GetCollection(baseCat.Url)
	productMeta, err := productCol.CreateDocument(context.Background(), pi)

	if err != nil {
		log.Println(222222222, err)

		return driver.DocumentMeta{}, err
	}
	log.Println(1111111, productMeta)
	fmt.Println(111111111)
	edgeCol := database.GetCollection(fmt.Sprintf("categories-%v", baseCat.Url))
	e := database.MyEdgeObject{
		From: catMeta.ID.String(),
		To:   productMeta.ID.String(),
	}
	meta, err := edgeCol.CreateDocument(context.Background(), e)
	if err != nil {
		return driver.DocumentMeta{}, err
	}
	return meta, nil

}

func titleMakerToTitle(fields []f, titleMaker titleMaker, brand string) string {
	valueUsedInTitleMaker := make([]string, 0)
	for _, field := range fields {
		for _, t := range titleMaker.Data {
			if t == field.Name {
				valueUsedInTitleMaker = append(valueUsedInTitleMaker, field.Value)
			}
		}
	}
	s := titleMaker.Title
	for _, s2 := range valueUsedInTitleMaker {
		s = strings.Replace(s, "$", s2, 1)
	}

	s = strings.Replace(s, "!", brand, 1)
	return s
}

func fieldsToSpecs(fields []f, mainSpecs []f, completeSpec []csType) ([]f, []csType) {

	for i, ms := range mainSpecs {
		for _, field := range fields {
			if field.Name == ms.Name {
				mainSpecs[i].Value = field.Value
			}
		}
	}

	for i, subSpec := range completeSpec {
		for j, ss := range subSpec.Items {
			for _, field := range fields {
				if field.Name == ss.Name {
					completeSpec[i].Items[j].Value = field.Value
				}
			}
		}
	}

	return mainSpecs, completeSpec

}

func createFilterString(fields []f) []string {
	strArray := make([]string, 0)
	for _, ff := range fields {
		strArray = append(strArray, fmt.Sprintf("%v=%v", ff.Name, ff.Value))

	}
	return strArray
}
