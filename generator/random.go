package generator

import "math/rand"

func generateRandomNumber(min, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}

func generateRandomNumberInt64(min, max int) int64 {
	if min == max {
		return int64(min)
	}
	return rand.Int63n(int64(max -min)) + int64(min)
}
