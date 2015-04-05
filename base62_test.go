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

	{"9223372036854775807", "AzL8n0Y58m7"},       // max signed int64
	{"9223372036854775809", "AzL8n0Y58m9"},       // beyond int64
	{"9223372036854775861", "AzL8n0Y58mz"},       //
	{"18446744073709551615", "LygHa16AHYF"},      // max uint64
	{"571849066284996100034", "AzL8n0Y58m70"},    // max int64 * 62
	{"35454642109669758202168", "AzL8n0Y58m70y"}, // (max int64 * 62^2) + 60

	{"170141183460469231731687303715884105727", "3tX16dB2jpss4tZORYcqo3"}, // max signed 128bit int
	{"170141183460469231731687303715884105757", "3tX16dB2jpss4tZORYcqoX"}, // max signed 128bit int + 30
	{"340282366920938463463374607431768211455", "7n42DGM5Tflk9n8mt7Fhc7"}, // max unsigned 128bit int
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
