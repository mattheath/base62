// Package base62 implements base62 encoding
package base62

import (
	"math"
	"strings"
)

const base = 62
const encodeStd = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// EncodeInt64 returns the base62 encoding of n
func EncodeInt64(n int64) string {
	var (
		s   []byte = make([]byte, 0)
		rem int64
	)

	// Progressively divide by base, store remainder each time
	// Prepend as an additional character is the higher power
	for n > 0 {
		rem = n % 62
		n = n / 62
		s = append([]byte{encodeStd[rem]}, s...)
	}

	return string(s)
}

// DecodeToInt64 decodes a base62 encoded string
func DecodeToInt64(s string) int64 {
	var (
		n     int64
		c     int64
		idx   int
		power int
	)

	for i, v := range s {
		idx = strings.IndexRune(encodeStd, v)

		// Work downwards through powers of our base
		power = len(s) - (i + 1)

		// Calculate value at this position and add
		c = int64(idx) * int64(math.Pow(float64(base), float64(power)))
		n = n + c
	}

	return int64(n)
}
