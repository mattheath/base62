package base62

import (
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var result string

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
		t.Logf("Decoded %s to %v", tc.encoded, v)
		assert.Equal(t, tc.num, v)
	}
}

func BenchmarkEncodeInt64Medium(b *testing.B) {
	var id string
	for n := 0; n < b.N; n++ {
		id = EncodeInt64(4815162342)
	}
	result = id
}

func BenchmarkEncodeInt64Long(b *testing.B) {
	var id string
	for n := 0; n < b.N; n++ {
		id = EncodeInt64(9223372036854775807)
	}
	result = id
}

var paddedTestcases = []struct {
	num     int64
	encoded string
}{
	{1, "000000000000001"},
	{9, "000000000000009"},
	{10, "00000000000000A"},
	{35, "00000000000000Z"},
	{36, "00000000000000a"},
	{61, "00000000000000z"},
	{62, "000000000000010"},
	{99, "00000000000001b"},
	{3844, "000000000000100"},
	{3860, "00000000000010G"},
	{4815162342, "0000000005Frvgk"},
	{9223372036854775807, "0000AzL8n0Y58m7"},
}

func TestEncodeInt64WithPadding(t *testing.T) {
	e := NewEncoding(encodeStd).Option(Padding(15))

	for _, tc := range paddedTestcases {
		v := e.EncodeInt64(tc.num)
		t.Logf("Encoded %v as %s", tc.num, v)
		assert.Equal(t, tc.encoded, v)
	}
}

func TestEncodeInt64WithVaryingPadding(t *testing.T) {
	testlens := []int{11, 23, 30}
	for _, tl := range testlens {
		e := NewEncoding(encodeStd).Option(Padding(tl))

		for _, tc := range paddedTestcases {
			v := e.EncodeInt64(tc.num)
			t.Logf("Encoded %v as %s", tc.num, v)
			assert.Equal(t, tl, len(v))
		}
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

	{"24467927614188555520896788267013", "8HFaR8qWtRlGDHnO57"}, // a few boundary flake id tests
	{"24467927614170108776823078715395", "8HFaR8qAulTgCBd6Wp"},
	{"24467927614170108776823078715394", "8HFaR8qAulTgCBd6Wo"},
	{"24467927614170108776823078715393", "8HFaR8qAulTgCBd6Wn"},
	{"24467927614170108776823078715392", "8HFaR8qAulTgCBd6Wm"},

	{"170141183460469231731687303715884105727", "3tX16dB2jpss4tZORYcqo3"}, // max signed 128bit int
	{"170141183460469231731687303715884105757", "3tX16dB2jpss4tZORYcqoX"}, // max signed 128bit int + 30
	{"340282366920938463463374607431768211455", "7n42DGM5Tflk9n8mt7Fhc7"}, // max unsigned 128bit int

	{"2707803647802660400290261537185326956543", "zzzzzzzzzzzzzzzzzzzzzz"}, // max 22 character when encoded
}

func TestEncodeBigInt(t *testing.T) {
	for _, tc := range bigTestcases {
		var (
			n  = new(big.Int)
			ok bool
		)

		n, ok = n.SetString(tc.num, 10)
		require.True(t, ok)

		v := EncodeBigInt(n)
		t.Logf("Encoded %v as %s", tc.num, v)
		assert.Equal(t, tc.encoded, v)
	}
}

func BenchmarkEncodeBigIntVeryLong(b *testing.B) {
	var (
		v = new(big.Int)
		s string
	)
	v.SetString("340282366920938463463374607431768211455", 10)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s = EncodeBigInt(v)
	}
	result = s
}

func TestDecodeToBigInt(t *testing.T) {
	for _, tc := range bigTestcases {
		v := DecodeToBigInt(tc.encoded)
		t.Logf("Decoded %v to %s", tc.encoded, v.String())
		assert.Equal(t, tc.num, v.String())
	}
}

// Ensure than padded base62 strings are correctly decoded

func TestPaddedDecodeToInt64(t *testing.T) {
	testcases := []struct {
		encoded string
		result  int64
	}{
		{"000005Frvgk", 4815162342},
		{"000000000AzL8n0Y58m7", 9223372036854775807},
	}

	for _, tc := range testcases {
		v := DecodeToInt64(tc.encoded)
		t.Logf("Decoded %s to %v", tc.encoded, v)
		assert.Equal(t, tc.result, v)
	}
}

func TestPaddedDecodeToBigInt(t *testing.T) {
	testcases := []struct {
		encoded string
		result  string
	}{
		{"0000000000000000000000000000000000005Frvgk", "4815162342"},
		{"000000000000000000003tX16dB2jpss4tZORYcqoX", "170141183460469231731687303715884105757"},
		{"000000000000000000007n42DGM5Tflk9n8mt7Fhc7", "340282366920938463463374607431768211455"},
	}

	for _, tc := range testcases {
		v := DecodeToBigInt(tc.encoded)
		t.Logf("Decoded %s to %v", tc.encoded, v)
		assert.Equal(t, tc.result, v.String())
	}
}

// TestLexicalPaddedSort tests that numbers encoded as base62 strings
// are correctly lexically sorted with the original order preserved
// if these are left padded to the same length.
//
// An alternative sort method which could be used to avoid padding
// would be a Shortlex sort, which sorts by cardinality, then lexicographically
func TestLexicalPaddedSort(t *testing.T) {

	var (
		lexicalOrder  sort.StringSlice = make([]string, 0)
		originalOrder                  = make([]string, 0)
	)

	// Generate lots of numbers, and encode them
	var i int64
	var modifier int64 = 1
	for i = 0; i < 100000; i++ {
		if i%10000 == 0 {
			modifier = modifier * 30
		}

		v := EncodeInt64(i + modifier)

		lexicalOrder = append(lexicalOrder, v)
		originalOrder = append(originalOrder, v)
	}

	// Find longest string & pad encoded strings to this length
	maxlen := len(originalOrder[len(originalOrder)-1])
	originalOrder = padStringArray(originalOrder, maxlen)
	lexicalOrder = padStringArray(lexicalOrder, maxlen)

	// Sort string array
	lexicalOrder.Sort()

	// Compare ordering with original
	var mismatch int64
	for i, v := range originalOrder {
		// t.Logf("%s %s", v, lexicalOrder[i])
		if lexicalOrder[i] != v {
			mismatch++
		}
	}
	assert.Equal(t, int64(0), mismatch, fmt.Sprintf("Expected zero mismatches, got %v", mismatch))
}

func padStringArray(s []string, minlen int) []string {
	for i, v := range s {
		s[i] = StdEncoding.pad(v, minlen)
	}
	return s
}
