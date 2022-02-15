package suppliers

import (
	"bamachoub-backend-go-v1/utils/middleware"
	"fmt"
	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/suppliers")
	r.Get("",
		middleware.Abac([]string{"manager", "admin", "super-admin"}),
		GetSuppliers)
	r.Post("", CreateSupplier)

	r.Get("/fav/:categoryUrl", middleware.GetSupplierByEmployee, getFavBySupplierKey)
	r.Get("/fav-product/:categoryUrl", middleware.GetSupplierByEmployee, getAllFavBySupplierKey)
	r.Put("/add-update-pool", middleware.SupplierEmployeeAuth([]string{}), addSupplierToUpdatePool)

	r.Post("/add-fav/:categoryUrl/:key", middleware.GetSupplierByEmployee, addFavorite)
	r.Post("/remove-fav/:key", middleware.GetSupplierByEmployee, deleteFavorite)

}

func WsRoutes(app fiber.Router) {
	r := app.Group("/suppliers-ws")

	ikisocket.On(ikisocket.EventDisconnect, func(ep *ikisocket.EventPayload) {
		// Remove the user from the local clients
		remove(ep.Kws.GetStringAttribute("supplier_key"))
		fmt.Println(fmt.Sprintf("Disconnection event - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})

	// On close event
	// This event is called when the server disconnects the user actively with .Close() method
	ikisocket.On(ikisocket.EventClose, func(ep *ikisocket.EventPayload) {
		// Remove the user from the local clients
		remove(ep.Kws.GetStringAttribute("supplier_key"))
		fmt.Println(fmt.Sprintf("Close event - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})

	// On error event
	ikisocket.On(ikisocket.EventError, func(ep *ikisocket.EventPayload) {
		fmt.Println(fmt.Sprintf("Error event - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})
	r.Get("/", middleware.GetSupplierByEmployee, ikisocket.New(newConn))

	//r.Post("/", func(c *fiber.Ctx) error {
	//	col:=database.GetCollection("onlineSuppliers")
	//	var u onlineUuid
	//	key:="9493985"
	//	_,err:=col.ReadDocument(context.Background(),key,&u)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	ikisocket.EmitTo(u.Uuid,[]byte("fuck you"))
	//	//NewEstelam([]string{"9493985"})
	//	return c.JSON("hi")
	//
	//})
}
