package session

import (
	"net/http"

	"golang-server/cmd/product/makanplace/controller/middlewares"
	"golang-server/cmd/product/makanplace/service/mk_user_session"
)

type Response struct {
	LoginUrls map[string]string `json:"login_urls"`

	UserInfo *mk_user_session.UserInfo `json:"user_info"`
}

func Register(mux *http.ServeMux, mkService *mk_user_session.Service, goauthloginurl string, mws middlewares.MiddewareStack) {
	// Revoke Session
	mux.Handle("POST /revoke_session", mws.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId := middlewares.GetSessionIdFromRequest(r)
		_ = mkService.RemoveSession(sessionId)
	})))
}
