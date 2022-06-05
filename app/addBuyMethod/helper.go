package addBuyMethod

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"sort"
)

func getCategoriesName(isPrice bool) ([]string, error) {
	var buyMethod string
	if isPrice {
		buyMethod = "price"
	} else {
		buyMethod = "estelam"
	}
	q1 := fmt.Sprintf("for i in categories filter i.status==\"start\" return concat(\"supplier_\",i.url,\"_%v\")", buyMethod)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q1, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var data []string
	for {
		var doc string
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err

		}
		data = append(data, doc)
	}
	return data, nil
}

func checkAndUpdateLowestCheckPrice(productKey string, categoryUrl string, in PriceIn) error {
	m := make(map[int64]string, 0)
	if in.OneMonthPrice != 0 {
		m[in.OneMonthPrice] = "one"
	}
	if in.TwoMonthPrice != 0 {
		m[in.TwoMonthPrice] = "two"
	}
	if in.ThreeMonthPrice != 0 {
		m[in.ThreeMonthPrice] = "three"
	}
	fmt.Println("map", m)
	var p Product
	pCol := database.GetCollection(categoryUrl)
	_, err := pCol.ReadDocument(context.Background(), productKey, &p)
	if err != nil {
		return err
	}
	lp := p.LowestCheckPrice.Price
	if lp <= 0 {
		lp = 9999999999999
	}
	arr := DirRange{lp}

	if in.OneMonthPrice != 0 {
		arr = append(arr, in.OneMonthPrice)
	}
	if in.TwoMonthPrice != 0 {
		arr = append(arr, in.TwoMonthPrice)

	}
	if in.ThreeMonthPrice != 0 {
		arr = append(arr, in.ThreeMonthPrice)

	}
	sort.Sort(arr)
	if arr[0] == lp {
		return nil
	}

	fmt.Println("arr", arr, arr[0], m)
	u := updateCheckPrice{LowestCheckPrice: checkPrice{
		Type:  m[arr[0]],
		Price: arr[0],
	}}
	fmt.Println(u)
	_, err = pCol.UpdateDocument(context.Background(), productKey, u)
	if err != nil {
		return err
	}
	return nil
}

type DirRange []int64

func (a DirRange) Len() int           { return len(a) }
func (a DirRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DirRange) Less(i, j int) bool { return a[i] < a[j] }
