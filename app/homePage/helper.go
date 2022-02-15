package homepage

import "fmt"

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//let a1=(return {type:"banner",numberOfBanners:4,data:[{imageUrl:"",link:""},{imageUrl:"",link:""},{imageUrl:"",link:""},{imageUrl:"",link:""},]}) return [a1[0],a1[0]]

func createBannerQuery(b Banners) string {
	dataString := ""
	for i, v := range b.Data {
		dataString += fmt.Sprintf(" {imageUrl:\"%v\", link: \"%v\" } ", v.ImageURL, v.Link)
		if i < len(b.Data)-1 {
			dataString += " , "
		}
	}

	q := fmt.Sprintf("return {numberOfBanners:%v,data:[%v]}", b.NumberOfBanners, dataString)

	return q

}

func createProductQuery(p ProductSlider) string {
	return fmt.Sprintf("for i in %v sort i.%v limit 15 return i", p.CategoryName, p.Sort)
}

func CreateCategoryQuery(c CategorySlider) string {
	return fmt.Sprintf("for i in categories filter i.url==\"%v\" for v in 1..2 outbound i graph \"categoryGraph\" return v", c.CategoryName)
}

func createCarousel(c []Carousel) string {
	q := ""
	for i, v := range c {
		q += fmt.Sprintf(" {imageUrl:\"%v\", link: \"%v\" } ", v.ImageURL, v.Link)
		if i < len(c)-1 {
			q += " , "
		}
	}
	query := fmt.Sprintf("return [%v]", q)

	return query
}




func insert(a []saveQuery, index int, value saveQuery) []saveQuery {
    if len(a) == index { // nil or empty slice or after last element
        return append(a, value)
    }
    a = append(a[:index+1], a[index:]...) // index < len(a)
    a[index] = value
    return a
}