package supplyWorkers

import (
	"github.com/arangodb/go-driver"
)

type SupplyWorker struct {
	FullName           string   `json:"fullName" validate:"required"`
	PhoneNumber        string   `json:"phoneNumber" validate:"required"`
	Email              string   `json:"email" validate:"required,email" `
	NationalCode       string   `json:"nationalCode" validate:"required"`
	Role               string   `json:"role" validate:"required"`
	Access             []string `json:"access"`
	HashPassword       string   `json:"hashPassword"`
	HashRefreshToken   string   `json:"hashRefreshToken"`
	SupplierKeyArray   []string `json:"supplierKeyArray" validate:"required"`
	CurrentSupplierKey string   `json:"currentSupplierKey"`
	FirstTimeLogin     bool     `json:"firstTimeLogin"`
}

type SupplyManagerCreateRequest struct {
	FullName     string `json:"fullName" validate:"required"`
	PhoneNumber  string `json:"phoneNumber" validate:"required"`
	Email        string `json:"email" validate:"required,email" `
	NationalCode string `json:"nationalCode" validate:"required"`
	SupplierKey  string `json:"supplierKey" validate:"required"`
}

type getSupplyWorkerByKeyType struct {
	driver.DocumentMeta
	SupplyWorker
}
