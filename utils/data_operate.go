package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"strconv"
)

var (
	eps float64 = 0.000001 //设置容忍度
)

// Hash ...
func Hash(data []byte) []byte {
	dst := make([]byte, 0x40)
	src := sha256.Sum256(data)
	hex.Encode(dst, src[:])
	return dst
}

// GetUID ...
func GetUID() string {
	u := make([]byte, 16)

	io.ReadFull(rand.Reader, u)
	return hex.EncodeToString(u)
}

// FloatEquals ...
func FloatEquals(a, b float64) bool {
	return math.Abs(a-b) < eps
}

// Decimal ...
func Decimal(value float64, decimal int) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf(fmt.Sprintf("%%.%vf", decimal), value), 64)
	return value
}

// MD5 ...
func MD5(str string) string {
	sum := md5.Sum([]byte(str))
	return hex.EncodeToString(sum[:])
}
