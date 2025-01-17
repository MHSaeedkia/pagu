package utils

import (
	"math/rand"
	"time"
)

const (
	Alphabet              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers               = "0123456789"
	CapitalAlphanumerical = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandomString(n int, input string) string {
	src := rand.NewSource(time.Now().UnixNano())
	out := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(input) {
			out[i] = input[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(out)
}
