package menu

import (
	"fmt"
	"log"
)

func getUnion(data []ReturnMenu ) []Item{
	orderSlice := make([]string,0)
	for _, d := range data {
		for _, item := range d.MenuItems {
			if !contains(orderSlice,item.Name) {
				orderSlice=append(orderSlice,item.Name)
			}
		}
	}
	fmt.Println(orderSlice)

	miArr:=make([][]Item,0)

	for _, u := range data {
		miArr=append(miArr,u.MenuItems)
	}

	newMap := make(map[string][]string)
	//log.Println(miArr)
	for _, m := range miArr {
		for _, m2 := range m {
			_,ok:=newMap[m2.Name]
			//log.Println(m2.Name,m2.Items,ok)
			if ok {
				a:=Union(m2.Items,newMap[m2.Name])
				newMap[m2.Name]=a

			}else {
				newMap[m2.Name]=m2.Items
			}
		}

	}
	log.Println(newMap)
	newArr:=make([]Item,0)
	for k, v := range newMap {
		temp:=Item{
			Name:  k,
			Items: v,
		}
		newArr=append(newArr,temp)
	}
	return newArr
}


func Union(a, b []string) []string {
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			a = append(a, item)
		}
	}
	return a
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
