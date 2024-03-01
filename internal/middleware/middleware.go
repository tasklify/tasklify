package middleare

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

func generateRandomString(length int) string {

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func CSPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxNonce := generateRandomString(16)
		responseTargetsNonce := generateRandomString(16)
		cssNonce := generateRandomString(16)

		// set then in context
		ctx := context.WithValue(r.Context(), "htmxNonce", htmxNonce)
		ctx = context.WithValue(ctx, "responseTargetsNonce", responseTargetsNonce)
		ctx = context.WithValue(ctx, "cssNonce", cssNonce)

		cspHeader := fmt.Sprintf("base-uri 'self'; default-src 'self'; script-src 'nonce-%s' 'nonce-%s'; style-src 'nonce-%s'; object-src 'none';", htmxNonce, responseTargetsNonce, cssNonce)

		fmt.Println(cspHeader)

		w.Header().Set("Content-Security-Policy", cspHeader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TextHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
