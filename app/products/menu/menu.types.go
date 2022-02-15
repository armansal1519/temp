package menu

import "github.com/arangodb/go-driver"

type CreateMenuDto struct {
	CategoryKey string `json:"categoryKey" validate:"required"`
	MenuItems   []Item `json:"menuItems" validate:"required"`
}

type Item struct {
	Name  string   `json:"name" validate:"required"`
	Items []string `json:"items" validate:"required"`
}

type addMenuKeyToCategory struct {
	MenuKey string `json:"menuKey"`
}

type ReturnMenu struct {
	driver.DocumentMeta
	CreateMenuDto
}
