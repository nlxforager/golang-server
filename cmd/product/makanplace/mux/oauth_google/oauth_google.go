package oauth_google

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	goauthservice "golang-server/cmd/product/makanplace/service/oauth/google"
	"google.golang.org/api/oauth2/v2"
	"log"
	"net/http"
	"sync"
)

type UserInfo struct {
	Id                int64
	GoogleCredentials []GoogleCredential
}

type GoogleCredential struct {
	UserInfo *oauth2.Userinfo
}

type SessionMap struct {
	m map[string]UserInfo
	l sync.Mutex
}

var sessionMap SessionMap

var userId int64

func Register(mux *http.ServeMux, service *goauthservice.Service) {
	mux.HandleFunc("/auth/google/login", func(w http.ResponseWriter, r *http.Request) {
		redirUrl := service.AuthCodeURL()
		http.Redirect(w, r, redirUrl, http.StatusTemporaryRedirect)
	})

	// client will provide authCode and state.
	// if we are able to obtain the service,
	mux.HandleFunc(service.AuthCodeSuccessCallbackPath(), func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic: %v", r)
			}
		}()
		log.Println(r.RequestURI)

		authCode := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		userInfo, err := service.UserInfo(state, authCode)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
		}
		infoB, err := json.Marshal(userInfo)
		if err != nil {
			log.Printf("%#v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sessionMap.l.Lock()
		cookieVal := uuid.New().String()
		userId++
		sessionMap.m[cookieVal] = UserInfo{
			Id: userId,
			GoogleCredentials: []GoogleCredential{
				{UserInfo: userInfo},
			},
		}
		sessionMap.l.Unlock()

		w.Header().Set("Set-Cookie", "makantoken"+"="+cookieVal+"; path=/")
		w.Write([]byte(fmt.Sprintf("%s", infoB)))
	})
}
