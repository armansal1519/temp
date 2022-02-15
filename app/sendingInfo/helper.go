package sendingInfo

import "fmt"

func addToInterval(array []interval, item interval) ([]interval, error) {
	if item.To <= item.From {
		return array, fmt.Errorf("to cannot be less than from")
	}
	if item.To > 24 {
		return array, fmt.Errorf("to cannot be bigger than 24")
	}
	if item.From < 1 {
		return array, fmt.Errorf("from can not be less than 1")
	}

	if item.To <= array[0].From {
		return insert(array, 0, item), nil
	}
	for i := 0; i < len(array)-1; i++ {
		if array[i].To <= item.From && array[i+1].From >= item.To {
			return insert(array, i+1, item), nil

		}
	}
	if item.From >= array[len(array)-1].To {
		return append(array, item), nil
	}
	return array, nil
}

func insert(a []interval, index int, value interval) []interval {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func remove(slice []interval, s int) []interval {
	return append(slice[:s], slice[s+1:]...)
}
