package productStructure

type CreateProductStructureDto struct {
	CategoryKey      string         `json:"categoryKey" validate:"required"`
	ProductFieldList []productField `json:"productFieldList" validate:"required"`
}

type productField struct {
	Name   string `json:"name" validate:"required"`
	IsList bool   `json:"isList" validate:"required"`
}

type addProductsStructureToCategory struct {
	ProductStructureKey string `json:"productStructureKey"`
}

type Category struct {
	Name      string `json:"name"`
	GraphPath string `json:"graphPath"`
	Status    string `json:"status"`
}

type updateIn struct {
	ProductStructureKey string `json:"productStructureKey"`
	//Operation    string `json:"operation"`
	Name string `json:"name"`
	IsList    bool `json:"isList"`
}