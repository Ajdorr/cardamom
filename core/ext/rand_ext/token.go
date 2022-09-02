package rand_ext

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GetRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Errorf("unable to read random bytes -- %w", err))
	}
	return base64.StdEncoding.EncodeToString(b)
}
