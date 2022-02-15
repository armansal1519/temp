package products

import (
	"bamachoub-backend-go-v1/app/products/productStructure"
	"bamachoub-backend-go-v1/config/database"
	"context"
	"fmt"
)

func validateProductBody(pb productBodyDto) error {
	//categoryKey := pb.CategoryKey
	//ps := getps(categoryKey)
	//fmt.Println(ps)

	//err := validateTitleMaker(pb.TitleMaker, ps)
	//if err != nil {
	//	return err
	//}
	////err = validateMainSpecification(pb.MainSpecification, ps)
	////if err != nil {
	////	return err
	////}
	////err = validateCompleteSpecification(pb.CompleteSpecification, ps)
	////if err != nil {
	////	return err
	////}

	return nil
}

func getps(key string) productStructure.CreateProductStructureDto {
	query := fmt.Sprintf("for i in productStructures filter i.categoryKey==\"%v\" limit 1 return i", key)
	db := database.GetDB()
	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		panic(fmt.Sprintf("error while running query:%v", query))
	}
	defer cursor.Close()

	var doc productStructure.CreateProductStructureDto
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		panic("error in cursor -in gps")
	}

	return doc
}

func validateTitleMaker(titleMaker string, ps productStructure.CreateProductStructureDto) error {
	r := []rune(titleMaker)
	//fmt.Println(titleMaker)
	var start []int
	var end []int
	for i, v := range r {
		if v == '[' {
			start = append(start, i)
		}
		if v == ']' {
			end = append(end, i)
		}
	}
	fmt.Println(start, end)
	fieldList := make([]string, 0)
	for i := 0; i < len(start); i++ {
		a := r[start[i]+1 : end[i]]
		fieldList = append(fieldList, string(a))
	}
	for _, v := range fieldList {
		isInProductStructure := false
		for _, u := range ps.ProductFieldList {

			if v == u.Name {
				isInProductStructure = true
				break
			}
		}
		if !isInProductStructure {
			err := fmt.Errorf("field %v does not exist in Product structure titleMaker", v)
			return err
		}
	}
	return nil
}

//func validateMainSpecification(main []string, ps productStructure.CreateProductStructureDto) error {
//	for _, v := range main {
//		isInProductStructure := false
//		for _, u := range ps.ProductFieldList {
//			if v == u.Name {
//				isInProductStructure = true
//				break
//			}
//
//		}
//		if !isInProductStructure {
//			err := fmt.Errorf("field %v does not exist in Product structure MainSpecification", v)
//			return err
//
//		}
//	}
//	return nil
//}
//
//func validateCompleteSpecification(sp []sp, ps productStructure.CreateProductStructureDto) error {
//	for _, v := range sp {
//		isInProductStructure := false
//		for _, y := range v.SpecList {
//			for _, u := range ps.ProductFieldList {
//				fmt.Printf("%v %v %v\n", v, u, y)
//				if y == u.Name {
//					isInProductStructure = true
//					break
//				}
//
//			}
//
//		}
//		if !isInProductStructure {
//			err := fmt.Errorf("field %v does not exist in Product structure CompleteSpecification", v)
//			return err
//
//		}
//	}
//	return nil
//}
