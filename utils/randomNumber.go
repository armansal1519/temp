package utils

import (
	"math/rand"
	"time"
)

func GenRandomNUmber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(max-min+1) + min
}
