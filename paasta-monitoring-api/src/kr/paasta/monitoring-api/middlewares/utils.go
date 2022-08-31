package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSha256(reqMsg string) string {
	data := []byte(reqMsg)
	hash1 := sha256.New()
	hash1.Write(data)
	md := hash1.Sum(nil)
	mdStr := hex.EncodeToString(md)

	return mdStr
}
