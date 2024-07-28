package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// LoggingMiddleware logs the details of each HTTP request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request context
		log.Println("Context:", r.Context())

		// Log the request headers
		log.Println("Headers:")
		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("%s: %s\n", name, value)
			}
		}

		// Log the request body
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println("Error reading body:", err)
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body
				log.Println("Body:", string(bodyBytes))
			}
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
