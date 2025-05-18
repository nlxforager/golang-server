package middlewares

import (
	"log"
	"net/http"

	log2 "golang-server/cmd/product/makanplace/log"
)

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s [middleware cors]\n", log2.HttpRequestPrefix(r.Context()))
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)

		log.Printf("served next handler. acao: %s \n ", w.Header().Get("Access-Control-Allow-Origin"))
	})
}
