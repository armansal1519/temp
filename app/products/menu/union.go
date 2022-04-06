package menu

import (
	"fmt"
	"log"
	"sort"
)

func getUnion(data []ReturnMenu) []Item {
	orderSlice := make([]string, 0)
	for _, d := range data {
		for _, item := range d.MenuItems {
			if !contains(orderSlice, item.Name) {
				orderSlice = append(orderSlice, item.Name)
			}
		}
	}
	fmt.Println(orderSlice)

	miArr := make([][]Item, 0)

	for _, u := range data {
		miArr = append(miArr, u.MenuItems)
	}

	newMap := make(map[string][]string)
	//log.Println(miArr)
	for _, m := range miArr {
		for _, m2 := range m {
			_, ok := newMap[m2.Name]
			//log.Println(m2.Name,m2.Items,ok)
			if ok {
				a := Union(m2.Items, newMap[m2.Name])
				newMap[m2.Name] = a

			} else {
				newMap[m2.Name] = m2.Items
			}
		}

	}
	log.Println(newMap)
	newArr := make([]Item, 0)
	for k, v := range newMap {
		temp := Item{
			Name:  k,
			Items: v,
		}
		newArr = append(newArr, temp)
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

func order(data []ReturnMenu, after []Item) []Item {

	nameArr := make([]orderType, 0)
	for _, datum := range data {
		for i, item := range datum.MenuItems {
			nameArr = append(nameArr, orderType{
				Name: item.Name,
				Rank: float64(i),
			})
		}
	}
	//for _, i2 := range nameArr {
	//	fmt.Println(i2)
	//
	//}

	merge := make(map[string]float64, 0)
	for _, n := range nameArr {
		v, ok := merge[n.Name]
		if ok {
			merge[n.Name] = (v + n.Rank) / 2
		} else {
			merge[n.Name] = n.Rank
		}
	}
	//for s, f := range merge {
	//	fmt.Println(s, " : ", f)
	//}
	rankedList := make([]orderType, 0)
	for k, v := range merge {
		rankedList = append(rankedList, orderType{
			Name: k,
			Rank: v,
		})
	}
	sort.SliceStable(rankedList, func(i, j int) bool {
		return rankedList[i].Rank < rankedList[j].Rank
	})
	for _, o := range rankedList {
		fmt.Println(o)
	}

	final := make([]Item, 0)
	for _, i := range rankedList {
		for _, item := range after {
			if item.Name == i.Name {
				final = append(final, item)
			}
		}
	}
	return final

}

//func insert(a []Item, index int, value Item) []Item {
//	if index >= len(a) {
//		return append(a, value)
//	}
//	if len(a) == index { // nil or empty slice or after last element
//		return append(a, value)
//	}
//	a = append(a[:index+1], a[index:]...) // index < len(a)
//	a[index] = value
//	return a
//}
