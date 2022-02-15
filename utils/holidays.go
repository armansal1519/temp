package utils

import (
	"encoding/json"
	"fmt"
	"github.com/jalaali/go-jalaali"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type holidayData struct {
	Events []struct {
		IsHoliday bool   `json:"isHoliday"`
		IsFriday  bool   `json:"isFriday"`
		Event     string `json:"event"`
	} `json:"events"`
	IsHoliday bool   `json:"isHoliday"`
	Date      string `json:"date"`
	Gdate     string `json:"gdate"`
}

func isHoliday(t time.Time) bool {
	y, m, d, _ := jalaali.ToJalaali(t.Year(), t.Month(), t.Day())
	months := []string{
		"فروردین", "اردیبهشت", "خرداد",
		"تیر", "مرداد", "شهریور",
		"مهر", "آبان", "آذر",
		"دی", "بهمن", "اسفند",
	}
	m_int := 0
	for i, month := range months {
		if m.String() == month {
			m_int = i + 1
		}
	}
	//Encode the data
	resp, err := http.Get(fmt.Sprintf("http://pholiday.herokuapp.com/date/%v-%v-%v", y, m_int, d))
	if err != nil {
		log.Println(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var h holidayData
	json.Unmarshal(body, &h)
	return h.IsHoliday
}
