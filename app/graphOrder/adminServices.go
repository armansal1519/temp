package graphOrder

import (
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"time"
)

func getOrderForAdmin(c *fiber.Ctx) error {
	//offset := c.Query("offset")
	//limit := c.Query("limit")
	//if offset == "" || limit == "" {
	//	return c.Status(400).SendString("Offset and Limit must have a value")
	//}

	//b := new(getOrderForAdminDto)
	//if err := utils.ParseBodyAndValidate(c, b); err != nil {
	//	return c.JSON(err)
	//}
	q := fmt.Sprintf("for u in users  \nlet order =(for v,e in 1..1 outbound u graph \"orderGraph\" return v)\nfor o in order \nlet payment =(for v,e in 1..1 outbound o graph \"orderGraph\"  return v)\nlet result=(for p in payment let oi=(for v,e in 1..1 outbound p graph \"orderGraph\" return v) return {payment:p,orderItems:oi})\nsort o.createdAt desc limit 0,10\nreturn {user:u,order:o,items:result} ")
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, q, nil)
	if err != nil {
		//fmt.Println(q)
		return c.Status(500).JSON(err)
	}
	defer cursor.Close()
	var data []GOrderResponseForAdminOut
	for {
		var doc GOrderResponseForAdminOut
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic("error in cursor -in GetAll")
		}
		data = append(data, doc)
	}

	//remove empty orders
	final := make([]GOrderResponseForAdminOut, 0)
	for _, datum := range data {
		flag := false
		for _, item := range datum.Items {
			if len(item.OrderItems) > 0 {
				flag = true
			}
		}
		if flag {
			final = append(final, datum)
		}
	}

	//calc status
	statusMap := make(map[int]string)
	statusMap[1] = "WaitingForPayment"
	statusMap[2] = "ApprovedBySupplier"
	statusMap[3] = "Processing"
	statusMap[4] = "Arrived"
	statusMap[5] = "Cancelled"
	statusMap[6] = "Referred"
	for ii, f := range final {
		statusArr := make([]int, 0)
		for _, i := range f.Items {
			for _, item := range i.OrderItems {
				if item.IsReferred {
					statusArr = append(statusArr, 6)
				} else if item.IsCancelled {
					statusArr = append(statusArr, 5)
				} else if item.IsArrived {
					statusArr = append(statusArr, 4)
				} else if item.IsProcessing {
					statusArr = append(statusArr, 3)
				} else if item.IsApprovedBySupplier {
					statusArr = append(statusArr, 2)
				} else if item.IsWaitingForPayment {
					statusArr = append(statusArr, 1)
				}
			}
		}

		statusScore := 6
		for _, i := range statusArr {
			if i < statusScore {
				statusScore = i
			}
		}
		final[ii].Order.Status = statusMap[statusScore]

	}

	//calc price
	for i, datum := range final {
		var totalPrice int64
		for _, j := range datum.Items {
			if datum.Order.Status == "wait-payment" {
				totalPrice += j.Payment.RemainingPrice
			} else {
				totalPrice += j.Payment.TotalPrice
			}

		}

		data[i].Order.TotalAmount = totalPrice

	}

	//add reserved
	for i, out := range final {
		for _, out2 := range out.Items {
			if out2.Payment.IsRejected {
				if time.Now().Unix() < out2.Payment.RejectionTime {
					final[i].Reserved = reservedInfo{
						IsReserved: true,
						TimeToEnd:  out2.Payment.RejectionTime - time.Now().Unix(),
					}
				}
			}
		}
	}
	return c.JSON(final)
}
