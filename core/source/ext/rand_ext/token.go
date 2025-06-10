package rand_ext

import (
	"cardamom/core/source/ext/log_ext"
	"crypto/rand"
	"encoding/base64"
)

func GetRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(log_ext.Errorf("unable to read random bytes -- %w", err))
	}
	return base64.StdEncoding.EncodeToString(b)
}
