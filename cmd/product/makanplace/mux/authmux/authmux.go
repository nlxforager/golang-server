package ping

import (
	"net/http"

	"golang-server/cmd/product/makanplace/service/mkusersessionservice"
)

type Response struct {
	LoginUrls map[string]string `json:"login_urls"`

	UserInfo *mkusersessionservice.UserInfo `json:"user_info"`
}

func Register(mux *http.ServeMux, makanTokenCookieKey string, mkService *mkusersessionservice.Service, goauthloginurl string) {
	// Revoke Session
	mux.HandleFunc("/revoke_session", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie(makanTokenCookieKey)
		var sessionId string
		if cookie != nil {
			sessionId = cookie.Value
		}

		resp := Response{LoginUrls: make(map[string]string)}

		resp.LoginUrls["google"] = goauthloginurl
		_ = mkService.RemoveSession(sessionId)
	})
}
