package ping

import (
	"encoding/json"
	"net/http"

	"golang-server/cmd/product/makanplace/controller/middlewares"
	"golang-server/cmd/product/makanplace/service/mkusersessionservice"

	"golang-server/cmd/product/makanplace/repositories/auth"
)

type Response struct {
	LoginUrls map[string]string `json:"login_urls"`

	UserInfo *auth.UserWithGmail `json:"user_info"`
}

func Register(mux *http.ServeMux, makanTokenCookieKey string, mkService *mkusersessionservice.Service, goauthloginurl string, mws middlewares.MiddewareStack) {
	// Checks current user state of the client
	// Provides server configuration values

	mux.Handle("/ping", mws.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
