package tools

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)

var Salt = "random_salt"

func HashPassword(password, salt string) string {
    h := hmac.New(sha256.New, []byte(salt))
    h.Write([]byte(password))
    return hex.EncodeToString(h.Sum(nil))
}