package transportation

import (
	"bamachoub-backend-go-v1/utils"
	"github.com/gofiber/fiber/v2"
)

//TODO fix transportation price

func GetTransportationPrice() int64 {
	return 1000
}

// getTransportationPrice get transportation price
// @Summary get transportation price
// @Description TransportationType must be bamachoub or user-address and SendingMethod must be fast or normal
// @Tags transportation
// @Accept json
// @Produce json
// @Param data body sendingInfo true "data"
// @Security ApiKeyAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} []string{}
// @Failure 500 {object} string{}
// @Failure 404 {object} string{}
// @Router /transportation [post]
func getTransportationPrice(c *fiber.Ctx) error {
	s := new(sendingInfo)
	//userKey := c.Locals("userKey").(string)

	if err := utils.ParseBodyAndValidate(c, s); err != nil {
		return c.JSON(err)
	}
	if s.TransportationType != "user-address" && s.TransportationType != "bamachoub" {
		return c.Status(400).SendString("TransportationType must be bamachoub or user-address but is : " + s.TransportationType)
	}
	if s.SendingMethod != "fast" && s.SendingMethod != "normal" {
		return c.Status(400).SendString("SendingMethod must be fast or normal but is : " + s.SendingMethod)
	}

	var p int64
	if s.TransportationType == "user-address" {
		//userAddressCol := database.GetCollection("userAddress")
		//var ao userAddress.AddressOut
		//_, err := userAddressCol.ReadDocument(context.Background(), s.AddressKey, &ao)
		//if err != nil {
		//	return c.Status(500).JSON(err)
		//}
		//if s.SendingMethod == "normal" {
		//	to, err := createTransportationObj(userKey, "normal", Location{
		//		Lat: ao.Latitude,
		//		Lan: ao.Longitude,
		//	})
		//	if err != nil {
		//		return c.Status(500).JSON(err)
		//	}
		//	return c.JSON(to)
		//}
		p = 1000
	} else {
		p = 0
	}

	//o, err := graphOrder.GetOrderAndPaymentsByUserKey(userKey)
	//if err != nil {
	//	return c.Status(500).JSON(err)
	//}
	//hasPrice := false
	//for _, out := range o.Payment {
	//	if out.Type == "price" {
	//		hasPrice = true
	//	}
	//}
	//if hasPrice {
	//	orderCol := database.GetCollection("gOrder")
	//	_, err = orderCol.UpdateDocument(context.Background(), o.Order.Key, updateTransportationPrice{TransportationPrice: p, TransportationPriceWithPrice: true})
	//	if err != nil {
	//		return c.Status(500).JSON(err)
	//	}
	//} else {
	//	orderCol := database.GetCollection("gOrder")
	//	_, err = orderCol.UpdateDocument(context.Background(), o.Order.Key, updateTransportationPrice{TransportationPrice: p, TransportationPriceWithPrice: false})
	//	if err != nil {
	//		return c.Status(500).JSON(err)
	//	}
	//}

	return c.JSON(fiber.Map{
		"transportationPrice": p,
	})
}
