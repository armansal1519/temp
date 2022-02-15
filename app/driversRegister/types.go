package driversRegister

import "github.com/arangodb/go-driver"

type carCategory struct {
	CarTypes   string  `json:"carTypes" validate:"required,min=2,max=32"`
	CarTonnage float64 `json:"carTonnage" validate:"required"`
	CarVolume  float64 `json:"carVolume" validate:"required"`
	CreatedAt  int64   `json:"createdAt"`
}

type getCategory struct {
	carCategory
	driver.DocumentMeta
}

// type Date struct {
// 	Day   int16 `json:"day" validate:"required, number"`
// 	Month int16 `json:"month" validate:"required, number"`
// 	Year  int16 `json:"year" validate:"required, number"`
// }

type Description struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Role string `json:"role"`
}

type driverInfo struct {
	FirstName            string        `json:"firstName" validate:"required"`
	LastName             string        `json:"lastName" validate:"required,min=2,max=32"`
	Gender               string        `json:"gender" validate:"required,min=2,max=32"`
	NationalNo           int64         `json:"nationalNo" validate:"required"`
	BirthdayDate         string        `json:"birthdayDate" validate:"required"`
	City                 string        `json:"city" validate:"required"`
	Province             string        `json:"province" validate:"required,min=2,max=32"`
	Address              string        `json:"address" validate:"required"`
	Latitude             float64       `json:"latitude" validate:"required"`
	Longitude            float64       `json:"longitude" validate:"required"`
	PostCode             int64         `json:"postCode" validate:"required"`
	PhoneNo              int64         `json:"phoneNo" validate:"required"`
	CardNo               int64         `json:"cardNo" validate:"required"`
	AccountNo            int64         `json:"accountNo" validate:"required"`
	CarType              string        `json:"carType" validate:"required"`
	PlateNo              string        `json:"plateNo" validate:"required"`
	AppearanceStatus     int32         `json:"appearanceStatus" validate:"required"`
	CarTechnical         int32         `json:"carTechnical" validate:"required"`
	InsuranceStatus      int64         `json:"insuranceStatus" validate:"required"`
	BarbandStatus        string        `json:"BarbandStatus" validate:"required"`
	CarColor             string        `json:"carColor" validate:"required"`
	CarTools             int32         `json:"carTools" validate:"required"`
	DriverEthics         int32         `json:"driverEthics" validate:"required"`
	Sabeghe              int32         `json:"sabeghe" validate:"required"`
	PhysicalCondition    int32         `json:"physicalCondition" validate:"required"`
	Panctuality          int32         `json:"panctuality" validate:"required"`
	IdCardImage          string        `json:"idCardImage" validate:"required"`
	IdBookPageOneImage   string        `json:"idBookPageOneImage" validate:"required"`
	IdBookPageTwoImage   string        `json:"idBookPageTwoImage" validate:"required"`
	LicenseImage         string        `json:"licenseImage" validate:"required"`
	FaceImage            string        `json:"faceImage" validate:"required"`
	InsurancePolicyImage string        `json:"insurancePolicyImage" validate:"required"`
	CreatedAt            int64         `json:"createdAt"`
	Description          []Description `json:"description"`
}

type getDriverInfo struct {
	driverInfo
	driver.DocumentMeta
}

type editDriverInfo struct {
	FirstName            string        `json:"firstName" validate:"required"`
	LastName             string        `json:"lastName" validate:"required,min=2,max=32"`
	Gender               string        `json:"gender" validate:"required,min=2,max=32"`
	NationalNo           int64         `json:"nationalNo" validate:"required"`
	BirthdayDate         string        `json:"birthdayDate" validate:"required"`
	City                 string        `json:"city" validate:"required"`
	Province             string        `json:"province" validate:"required,min=2,max=32"`
	Address              string        `json:"address" validate:"required"`
	Latitude             float64       `json:"latitude" validate:"required"`
	Longitude            float64       `json:"longitude" validate:"required"`
	PostCode             int64         `json:"postCode" validate:"required"`
	PhoneNo              int64         `json:"phoneNo" validate:"required"`
	CardNo               int64         `json:"cardNo" validate:"required"`
	AccountNo            int64         `json:"accountNo" validate:"required"`
	CarType              string        `json:"carType" validate:"required"`
	PlateNo              string        `json:"plateNo" validate:"required"`
	AppearanceStatus     int32         `json:"appearanceStatus" validate:"required"`
	CarTechnical         int32         `json:"carTechnical" validate:"required"`
	InsuranceStatus      int64         `json:"insuranceStatus" validate:"required"`
	BarbandStatus        string        `json:"BarbandStatus" validate:"required"`
	CarColor             string        `json:"carColor" validate:"required"`
	CarTools             int32         `json:"carTools" validate:"required"`
	DriverEthics         int32         `json:"driverEthics" validate:"required"`
	Sabeghe              int32         `json:"sabeghe" validate:"required"`
	PhysicalCondition    int32         `json:"physicalCondition" validate:"required"`
	Panctuality          int32         `json:"panctuality" validate:"required"`
	IdCardImage          string        `json:"idCardImage" validate:"required"`
	IdBookPageOneImage   string        `json:"idBookPageOneImage" validate:"required"`
	IdBookPageTwoImage   string        `json:"idBookPageTwoImage" validate:"required"`
	LicenseImage         string        `json:"licenseImage" validate:"required"`
	FaceImage            string        `json:"faceImage" validate:"required"`
	InsurancePolicyImage string        `json:"insurancePolicyImage" validate:"required"`
	CreatedAt            int64         `json:"createdAt"`
	Description          []Description `json:"description"`
}

type input struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type carType struct {
	Type string `json:"type" validate:"required"`
}
