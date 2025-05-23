package oauth_google

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang-server/cmd/product/makanplace/controller/response_types"
	mklog "golang-server/cmd/product/makanplace/httplog"
	"golang-server/cmd/product/makanplace/service/mk_user_session"
	goauthservice "golang-server/cmd/product/makanplace/service/oauth/google"

	"google.golang.org/api/oauth2/v2"
)

type Session struct {
	Session string `json:"session_id"`
}

func Register(mux *http.ServeMux, gOAuthService *goauthservice.Service, mkService *mk_user_session.Service, goauthloginurl string) {
	mux.HandleFunc(goauthloginurl, func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		if origin == "" {
			referer := r.Header.Get("Referer")
			parsedURL, err := url.Parse(referer)
			if err != nil {
				response_types.ErrorNoBody(w, http.StatusForbidden, fmt.Errorf("unkown client host"))
			}

			// Extract the scheme, host, and port (if present)
			origin = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
		}
		redirUrl, _ := gOAuthService.AuthCodeURL(origin)
		http.Redirect(w, r, redirUrl, http.StatusTemporaryRedirect)
	})

	// client will provide authCode and state.
	// if we are able to exchange a valid access token and use it to obtain user info, we will associate the google credential to a makanplace user.
	mux.HandleFunc(gOAuthService.AuthCodeSuccessCallbackPath(), func(w http.ResponseWriter, r *http.Request) {
		authCode := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		gmailUserInfo, origin, err := gOAuthService.UserInfo(state, authCode)
		fmt.Printf("%s gmailUserInfo %#v origin %s\n", mklog.SPrintHttpRequestPrefix(r), gmailUserInfo, origin)
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
		fmt.Printf("%s sessionId %s\n", mklog.SPrintHttpRequestPrefix(r), sessionId)

		//w.Header().Set("Set-Cookie", makanTokenCookieKey+"="+sessionId+"; path=/; HttpOnly; SameSite=None; Secure;")
		//w.Header().Set("Content-Type", "application/json")
		//
		//var resp response_types.Response[Session]
		//resp.Data = Session{Session: sessionId}
		//resp.Error = nil
		//b, _ := json.Marshal(r)
		//w.Write(b)

		referrer := origin + "/auth_callback?session_id=" + sessionId
		fmt.Printf("%s referrer %s\n", mklog.SPrintHttpRequestPrefix(r), referrer)
		http.Redirect(w, r, referrer, http.StatusTemporaryRedirect)
	})
}
