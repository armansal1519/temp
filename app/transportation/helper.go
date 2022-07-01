package transportation

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"encoding/json"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

func (l Location) isEmpty() bool {
	if l.Lat == -1 && l.Lan == -1 {
		return true
	}
	return false
}

func getDistance(from, to Location) float64 {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.neshan.org/v1/distance-matrix?type=car&origins=%v,%v&destinations=%v,%v", from.Lat, from.Lan, to.Lat, to.Lan), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(from, to)
	req.Header.Set("Api-Key", "service.HfmheE7dYqtDDt8cVebzun8vciaxOI2wYoNd3Vq1")
	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var n NeshanResponse
	err = json.Unmarshal(body, &n)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	fmt.Println(n.Rows[0].Elements[0].Distance.Value)
	fmt.Println(4)

	return float64(n.Rows[0].Elements[0].Distance.Value) / 1000
}

func createTransportationObj(userKey string, tType string, userLoc Location) (*[]transportationObj, error) {
	soi, err := getSupplierOrderItem(userKey)
	if err != nil {
		return nil, err
	}

	sv, err := createSupplierVolume(*soi)
	if tType == "normal" {
		fmt.Println(1)

		trObjArr := make([]transportationObj, 0)
		for _, volume := range sv {
			for _, f := range volume.VolumeArr {
				if f == 1.1 {
					supplierLoc := Location{
						Lat: float64(volume.SupplierKey.Latitude),
						Lan: float64(volume.SupplierKey.Longitude),
					}
					dd := getDistance(supplierLoc, userLoc)
					trObjArr = append(trObjArr, transportationObj{
						Id: uuid.New().String(),
						From: locationInfo{
							LocationId: volume.SupplierKey.Id,
							Location:   supplierLoc,
						},
						To: locationInfo{
							LocationId: "user/" + userKey,
							Location:   userLoc,
						},
						IsToUser:   true,
						Distance:   dd,
						Percentage: 1.1,
						Price:      int64((dd * 5600) + 160000),
					})
				}
			}
		}
		fmt.Println(2)

		for _, volume := range sv {
			fmt.Println(1, volume.VolumeArr)
		}
		svGreedyArr := make(greedyList, 0)
		for _, v := range sv {
			svGreedyArr = append(svGreedyArr, greedyItem{
				Supplier: v.SupplierKey,
				Volume:   v.VolumeArr[0],
			})
		}

		sort.Sort(svGreedyArr)
		usedIndex := make([]int, 0)
		sList := make([][]supplier, 0)

		for _, item := range svGreedyArr {
			fmt.Println(item.Volume)
		}
		for i, item := range svGreedyArr {
			sum := item.Volume
			temp := []supplier{item.Supplier}
			if contains(usedIndex, i) {
				continue
			}
			for i2, g := range svGreedyArr {
				if i == i2 {
					continue
				}
				if contains(usedIndex, i2) {
					continue
				}
				if sum+item.Volume < 1.1 {
					usedIndex = append(usedIndex, i2)
					temp = append(temp, g.Supplier)
				}

			}
			sList = append(sList, temp)

		}
		for _, suppliers := range sList {
			id := uuid.New().String()
			if len(suppliers) > 1 {
				for j := 0; j < len(suppliers)-2; j++ {
					supplier1Loc := Location{
						Lat: float64(suppliers[j].Latitude),
						Lan: float64(suppliers[j].Longitude),
					}
					supplier2Loc := Location{
						Lat: float64(suppliers[j+1].Latitude),
						Lan: float64(suppliers[j+1].Longitude),
					}
					dd := getDistance(supplier1Loc, supplier2Loc)
					trObjArr = append(trObjArr, transportationObj{
						Id: id,
						From: locationInfo{
							LocationId: suppliers[j].Id,
							Location:   supplier1Loc,
						},
						To: locationInfo{
							LocationId: suppliers[j+1].Id,
							Location:   supplier2Loc,
						},
						IsToUser:   false,
						Distance:   dd,
						Percentage: -1,
						Price:      int64(dd * 5600),
					})
				}
			} else {
				for j := 0; j < len(suppliers)-2; j++ {
					supplier1Loc := Location{
						Lat: float64(suppliers[j].Latitude),
						Lan: float64(suppliers[j].Longitude),
					}

					dd := getDistance(supplier1Loc, userLoc)
					trObjArr = append(trObjArr, transportationObj{
						Id: id,
						From: locationInfo{
							LocationId: suppliers[j].Id,
							Location:   supplier1Loc,
						},
						To: locationInfo{
							LocationId: "user/" + userKey,
							Location:   userLoc,
						},
						IsToUser:   true,
						Distance:   dd,
						Percentage: -1,
						Price:      int64((dd * 5600) + 160000),
					})
				}
			}
		}

		return &trObjArr, nil

	}

	return nil, nil

}

type greedyList []greedyItem

func (e greedyList) Len() int {
	return len(e)
}

func (e greedyList) Less(i, j int) bool {
	return e[i].Volume > e[j].Volume
}

func (e greedyList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func createSupplierVolume(soi []supplierOrderItem) ([]supplierVolume, error) {
	sv := make([]supplierVolume, 0)
	for _, s := range soi {
		var sum float64
		temp := supplierVolume{
			SupplierKey: s.Supplier,
			VolumeArr:   []float64{},
		}
		for _, s2 := range s.Data {
			v := s2.TrData.PercentageInNissan
			sum += v * float64(s2.OrderItem.Number)
		}
		var count float64
		for true {
			count = sum
			if count >= 1.1 {
				fmt.Println(count)
				temp.VolumeArr = append(temp.VolumeArr, 1.1)
				sum = sum - 1.1
			} else {
				temp.VolumeArr = append(temp.VolumeArr, sum)
				break
			}
			count = count - 1.1
		}
		sv = append(sv, temp)
	}
	for i, _ := range sv {
		sort.Float64s(sv[i].VolumeArr)
		//fmt.Println(sv[i].VolumeArr)
	}

	return sv, nil
}

func getSupplierOrderItem(userKey string) (*[]supplierOrderItem, error) {
	q := fmt.Sprintf("for i in gOrderItem  \nCOLLECT supplierKeys = i.supplierKey into orderItems\nlet data=(for j in orderItems[*].i for i in productTransportationData filter i.id==j.productId return {orderItem:j,trData:i})\nlet s=(for i in suppliers filter i._key==supplierKeys return i)\nreturn {supplier:s[0],supplierKeys,data}")
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var data []supplierOrderItem
	for {
		var doc supplierOrderItem
		_, err = cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Println("error in cursor -in GetAll")
			return nil, err
		}
		data = append(data, doc)
	}
	return &data, nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
