package oauth_google

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang-server/cmd/product/makanplace/service/mkusersessionservice"
	goauthservice "golang-server/cmd/product/makanplace/service/oauth/google"

	"google.golang.org/api/oauth2/v2"
)

func Register(mux *http.ServeMux, makanTokenCookieKey string, goAuthService *goauthservice.Service, mkService *mkusersessionservice.Service) {
	mux.HandleFunc("/auth/google/login", func(w http.ResponseWriter, r *http.Request) {
		redirUrl := goAuthService.AuthCodeURL()
		http.Redirect(w, r, redirUrl, http.StatusTemporaryRedirect)
	})

	// client will provide authCode and state.
	// if we are able to exchange a valid access token and use it to obtain user info, we will associate the google credential to a makanplace user.

	log.Println(goAuthService.AuthCodeSuccessCallbackPath())
	mux.HandleFunc(goAuthService.AuthCodeSuccessCallbackPath(), func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic: %v", r)
			}
		}()

		authCode := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		userInfo, err := goAuthService.UserInfo(state, authCode)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
		}
		infoB, err := json.Marshal(userInfo)
		if err != nil {
			log.Printf("%#v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sessionId, err := mkService.CreateUserSession([]*oauth2.Userinfo{userInfo})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Set-Cookie", makanTokenCookieKey+"="+sessionId+"; path=/")
		w.Write([]byte(fmt.Sprintf("%s", infoB)))
	})
}
