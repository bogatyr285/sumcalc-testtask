package sha256

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

type SHA256Hasher struct{}

func NewSHA256Hasher() *SHA256Hasher {
	return &SHA256Hasher{}
}

func SHA256Byte(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (h SHA256Hasher) Hash(n int64) string {
	b := big.NewInt(n).Bytes()
	return SHA256Byte(b)
}
