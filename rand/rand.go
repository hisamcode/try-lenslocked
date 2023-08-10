package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("rand.Bytes: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("rand.Bytes: didn't read enough random bytes")
	}
	return b, nil
}

// n is the number of bytes being used to generate the random string
func String(n int) (*string, error) {
	b, err := Bytes(n)
	if err != nil {
		return nil, fmt.Errorf("rand.String: %w", err)
	}
	encode := base64.URLEncoding.EncodeToString(b)
	return &encode, nil
}
