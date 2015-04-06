# base62

[![Build Status](https://travis-ci.org/mattheath/base62.svg?branch=master)](https://travis-ci.org/mattheath/base62) [![GoDoc](https://godoc.org/github.com/mattheath/base62?status.svg)](https://godoc.org/github.com/mattheath/base62)

`base62` implements base62 encoding for integer numbers of arbitrary precision.

Currently both `int64` and `*big.Int` are supported

## Encoding

The encoding is different to base64 encoding in two ways:
 * Reduced character set, without symbols - `0-9A-Za-z`
 * Character ordering - base64 is `A-Za-z0-9`

Note that this character ordering is also different to most other base62 encodings, but _preserves the sort order of encoded values_.

As a result, the encoded strings can be lexically sorted if either a sorting order such as [shortlex](http://en.wikipedia.org/wiki/Shortlex_order) is used, or if the resulting strings are left padded to ensure a consistent length.

## Usage

First `go get` the package:
```
go get github.com/mattheath/base62
```

This can then be used as below:
```go
package main

import (
    "fmt"
    "math/big"
    "github.com/mattheath/base62"
)

func main() {
	// Encoding 64 bit integers
	var n int64 = 4815162342
    encoded := base62.EncodeInt64(n)
    fmt.Println(encoded) // prints 5Frvgk

    // Arbitrary precision integers can be specified using the math/big pkg
    var b *big.Int = new(big.Int)
    b.SetString("340282366920938463463374607431768211455") // 128bit unsigned int
    bigEncoded := base62.EncodeBigInt(b)
    fmt.Println(bigEncoded) // prints 7n42DGM5Tflk9n8mt7Fhc7

    // Padding can be specified using an option on an Encoding
    // eg. to pad strings to a minimum length of 15 chars:
    e := base62.NewStdEncoding().Option(Padding(15))
    encoded = e.EncodeInt64(n)
    fmt.Println(encoded) // prints 0000000005Frvgk
}
```
