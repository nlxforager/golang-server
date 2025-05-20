package oauth_google

import (
	"log"
	"net/http"

	mklog "golang-server/cmd/product/makanplace/httplog"
	"golang-server/cmd/product/makanplace/service/mk_user_session"
	goauthservice "golang-server/cmd/product/makanplace/service/oauth/google"

	"google.golang.org/api/oauth2/v2"
)

func Register(mux *http.ServeMux, makanTokenCookieKey string, gOAuthService *goauthservice.Service, mkService *mk_user_session.Service, goauthloginurl string) {
	mux.HandleFunc(goauthloginurl, func(w http.ResponseWriter, r *http.Request) {
		redirUrl := gOAuthService.AuthCodeURL()
		http.Redirect(w, r, redirUrl, http.StatusTemporaryRedirect)
	})

	// client will provide authCode and state.
	// if we are able to exchange a valid access token and use it to obtain user info, we will associate the google credential to a makanplace user.
	mux.HandleFunc(gOAuthService.AuthCodeSuccessCallbackPath(), func(w http.ResponseWriter, r *http.Request) {
		authCode := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		gmailUserInfo, err := gOAuthService.UserInfo(state, authCode)
		if err != nil {
			log.Printf("%s Error getting user info (UserInfo): %v", mklog.SPrintHttpRequestPrefix(r), err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		sessionId, err := mkService.CreateUserSession([]*oauth2.Userinfo{gmailUserInfo})
		if err != nil {
			log.Printf("%s Error getting user info (CreateUserSession): %v", mklog.SPrintHttpRequestPrefix(r), err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Set-Cookie", makanTokenCookieKey+"="+sessionId+"; path=/; HttpOnly; SameSite=None; Secure;")
		referrer := r.Header.Get("Referer")
		http.Redirect(w, r, referrer, http.StatusTemporaryRedirect)
	})
}
