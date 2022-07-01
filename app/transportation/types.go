package transportation

type sendingInfo struct {
	AddressKey          string `json:"addressKey" validate:"required"`
	SendIntervalString  string `json:"sendIntervalString" validate:"required"`
	SendingMethod       string `json:"sendingMethod" validate:"required"`
	TransportationPrice int64  `json:"transportationPrice"`
	TransportationType  string `json:"transportationType"`
}
type sendingInfoOut struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	sendingInfo
}

type Location struct {
	Lat float64 `json:"lat"`
	Lan float64 `json:"lan"`
}

type NeshanResponse struct {
	Status string `json:"status"`
	Rows   []struct {
		Elements []struct {
			Status   string `json:"status"`
			Duration struct {
				Value int    `json:"value"`
				Text  string `json:"text"`
			} `json:"duration"`
			Distance struct {
				Value int    `json:"value"`
				Text  string `json:"text"`
			} `json:"distance"`
		} `json:"elements"`
	} `json:"rows"`
	OriginAddresses      []string `json:"origin_addresses"`
	DestinationAddresses []string `json:"destination_addresses"`
}

//type supplierOrderItem struct {
//	Supplier     suppliers.Supplier       `json:"supplier"`
//	SupplierKeys string                   `json:"supplierKeys"`
//	OrderItemArr graphOrder.GOrderItemOut `json:"orderItemArr"`
//}

type supplier struct {
	Key              string   `json:"_key"`
	Id               string   `json:"_id"`
	Rev              string   `json:"_rev"`
	Address          string   `json:"address"`
	State            string   `json:"state"`
	City             string   `json:"city"`
	Latitude         int      `json:"latitude"`
	Longitude        int      `json:"longitude"`
	Name             string   `json:"name"`
	Code             string   `json:"code"`
	Area             int      `json:"area"`
	AreaWithRoof     int      `json:"areaWithRoof"`
	PhoneNumber      string   `json:"phoneNumber"`
	CategoriesToSale []string `json:"categoriesToSale"`
	WalletAmount     int      `json:"walletAmount"`
	Status           string   `json:"status"`
	CreateAt         int      `json:"createAt"`
}

type supplierOrderItem struct {
	Supplier     supplier `json:"supplier"`
	SupplierKeys string   `json:"supplierKeys"`
	Data         []struct {
		OrderItem struct {
			Key                    string  `json:"_key"`
			Id                     string  `json:"_id"`
			Rev                    string  `json:"_rev"`
			PriceId                string  `json:"priceId"`
			SupplierKey            string  `json:"supplierKey"`
			ProductId              string  `json:"productId"`
			PricePerNumber         int     `json:"pricePerNumber"`
			Number                 int     `json:"number"`
			Variant                string  `json:"variant"`
			ProductTitle           string  `json:"productTitle"`
			ProductImageUrl        string  `json:"productImageUrl"`
			UserKey                string  `json:"userKey"`
			UserAuthType           string  `json:"userAuthType"`
			CommissionPercent      float64 `json:"commissionPercent"`
			CheckCommissionPercent int     `json:"checkCommissionPercent"`
			IsWaitingForPayment    bool    `json:"isWaitingForPayment"`
			IsApprovedBySupplier   bool    `json:"isApprovedBySupplier"`
			SupplierEmployeeId     string  `json:"supplierEmployeeId"`
			IsProcessing           bool    `json:"isProcessing"`
			IsArrived              bool    `json:"isArrived"`
			IsCancelled            bool    `json:"isCancelled"`
			CancelledById          string  `json:"cancelledById"`
			IsReferred             bool    `json:"isReferred"`
			ReferredReason         string  `json:"referredReason"`
			CreatedAt              int     `json:"createdAt"`
		} `json:"orderItem"`
		TrData struct {
			Key                     string  `json:"_key"`
			Id                      string  `json:"_id"`
			Rev                     string  `json:"_rev"`
			Id1                     string  `json:"id"`
			NumberInPallet          string  `json:"number_in_pallet"`
			Thickness               string  `json:"thickness"`
			Dimension               string  `json:"dimension"`
			DimensionTitle          string  `json:"dimension_title"`
			Volume                  float64 `json:"volume"`
			PercentageInNissan      float64 `json:"percentage_in_nissan"`
			MaxPalletInNissan       int     `json:"max_pallet_in_nissan"`
			PercentageInPaykanvanet float64 `json:"percentage_in_paykanvanet"`
			MaxPalletInPaykanvanet  int     `json:"max_pallet_in_paykanvanet"`
			PercentageInIsuzo       float64 `json:"percentage_in_isuzo"`
			MaxPalletInIsozo        int     `json:"max_pallet_in_isozo"`
			Percent                 int     `json:"percent"`
			Position                string  `json:"position"`
		} `json:"trData"`
	} `json:"data"`
}

type transportationObj struct {
	Id         string         `json:"id"`
	From       locationInfo   `json:"from"`
	To         locationInfo   `json:"to"`
	IsToUser   bool           `json:"IsToUser"`
	Distance   float64        `json:"distance"`
	Percentage float64        `json:"percentage"`
	Stops      []locationInfo `json:"stops"`
	Price      int64          `json:"price"`
}

type locationInfo struct {
	LocationId string `json:"locationId"`
	Location
}

type supplierVolume struct {
	SupplierKey supplier  `json:"supplier"`
	VolumeArr   []float64 `json:"volumeArr"`
}

type greedyItem struct {
	Supplier supplier `json:"supplier"`
	Volume   float64  `json:"volume"`
}

type updateTransportationPrice struct {
	TransportationPrice          int64 `json:"transportationPrice"`
	TransportationPriceWithPrice bool  `json:"transportationPriceWithPrice"`
}
