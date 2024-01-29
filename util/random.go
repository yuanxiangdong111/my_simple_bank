package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//r.Seed()
	rand.Seed(time.Now().UnixNano())
}

// 返回min-max
func RandomInt(min, max int64) int64 {
	// rand.Int63n(max - min + 1) 返回 0 -> max - min
	// 所以整体返回min -> max
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// 随机使用者名字
func RandomOwner() string {
	return RandomString(6)
}

// 随机的钱
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// 随机的货币名
func RandomCurrency() string {
	currenties := []string{"EUR", "USD", "CAD"}
	n := len(currenties)
	return currenties[rand.Intn(n)]
}
