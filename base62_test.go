package base62

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testcases = []struct {
	num     int64
	encoded string
}{
	{1, "1"},
	{9, "9"},
	{10, "A"},
	{35, "Z"},
	{36, "a"},
	{61, "z"},
	{62, "10"},
	{99, "1b"},
	{3844, "100"},
	{3860, "10G"},
	{4815162342, "5Frvgk"},
	{9223372036854775807, "AzL8n0Y58m7"},
}

func TestEncodeInt64(t *testing.T) {
	for _, tc := range testcases {
		v := EncodeInt64(tc.num)
		t.Logf("Encoded %v as %s", tc.num, v)
		assert.Equal(t, tc.encoded, v)
	}
}

func TestDecodeToInt64(t *testing.T) {
	for _, tc := range testcases {
		v := DecodeToInt64(tc.encoded)
		t.Logf("Decoded %s as %v", tc.encoded, v)
		assert.Equal(t, tc.num, v)
	}
}

var bigTestcases = []struct {
	num     string
	encoded string
}{
	{"1", "1"},
	{"9", "9"},
	{"10", "A"},
	{"35", "Z"},
	{"36", "a"},
	{"61", "z"},
	{"62", "10"},
	{"99", "1b"},
	{"3844", "100"},
	{"3860", "10G"},
	{"4815162342", "5Frvgk"},
	{"9223372036854775807", "AzL8n0Y58m7"},
}

func TestEncodeBigInt(t *testing.T) {
	for _, tc := range bigTestcases {
		var (
			n  *big.Int = new(big.Int)
			ok bool
		)

		n, ok = n.SetString(tc.num, 10)
		require.True(t, ok)

		v := EncodeBigInt(n)
		t.Logf("Encoded %v as %s", tc.num, v)
		assert.Equal(t, tc.encoded, v)
	}
}
