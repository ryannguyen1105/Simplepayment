package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const number = "0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	if min > max {
		return min
	}
	return min + rand.Int63n(max-min+1)
}

func RandomAlphabet(length int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomNumber(length int) string {
	var sb strings.Builder
	k := len(number)

	for i := 0; i < length; i++ {
		c := number[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

//func RandomEmail() string {
//	return RandomAlphabet(8) + RandomNumber(4) + "@gmail.com"
//}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomAlphabet(6))
}

func RandomUsername() string {
	return RandomAlphabet(5) + RandomNumber(4)
}

func RandomPassword() string {
	return RandomAlphabet(8) + RandomNumber(4)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomOwner() string {
	return RandomAlphabet(5)
}

func RandomStatus() string {
	ranStatus := []string{"pending", "completed", "failed"}
	n := len(ranStatus)
	return ranStatus[rand.Intn(n)]
}
