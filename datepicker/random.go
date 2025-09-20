package datepicker

import (
	"crypto/rand"
	"encoding/base64"
)

func randomString(n int) string {
	if n <= 0 {
		n = 16
	}
	b := make([]byte, n)
	_, _ = rand.Read(b)
	// Base64 grows size; trim to requested length and keep only URL-safe chars
	s := base64.RawURLEncoding.EncodeToString(b)
	if len(s) > n {
		s = s[:n]
	}
	return s
}
