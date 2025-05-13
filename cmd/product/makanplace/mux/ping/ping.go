package ping

import (
	"encoding/json"
	"net/http"

	"golang-server/cmd/product/makanplace/service/mkusersessionservice"
)

type Response struct {
	LoginUrls map[string]string `json:"login_urls"`

	UserInfo *mkusersessionservice.UserInfo `json:"user_info"`
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Register(mux *http.ServeMux, makanTokenCookieKey string, mkService *mkusersessionservice.Service, goauthloginurl string) {
	// Checks current user state of the client
	// Provides server configuration values
	mux.Handle("/ping", withCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie(makanTokenCookieKey)
		var sessionId string
		if cookie != nil {
			sessionId = cookie.Value
		}

		resp := Response{LoginUrls: make(map[string]string)}

		resp.LoginUrls["google"] = goauthloginurl

		session := mkService.GetSession(sessionId, true)
		resp.UserInfo = session

		respB, _ := json.Marshal(resp)
		w.Write(respB)
	})))

}
