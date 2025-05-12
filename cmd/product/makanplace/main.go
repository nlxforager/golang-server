package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	oauth "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func antiCsrfState() string {
	return "SOME_STATE"
}

const COOKIE_NAME_GOOGLE_AUTHED_BEFORE = "loginoncebefore"

func main() {
	if err := Init(); err != nil {
		log.Fatal(err)
	}
	eCh := make(chan os.Signal, 1)
	signal.Notify(eCh, syscall.SIGTERM, syscall.SIGINT)

	var googleOauthConfig = &oauth.Config{
		RedirectURL:  "http://localhost" + LISTENING_PORT + AUTH_CODE_SUCCESS_CALLBACK_PATH,
		ClientID:     CLIENT_ID_PREFIX + ".apps.googleusercontent.com",
		ClientSecret: CLIENT_SECRET,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	redirUrl := googleOauthConfig.AuthCodeURL(antiCsrfState(), oauth.SetAuthURLParam("prompt", "consent select_account"))

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/google/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirUrl, http.StatusTemporaryRedirect)
	})

	hc := http.DefaultClient
	if ENABLE_LOG_REQUEST {
		hc = &http.Client{
			Transport: &LoggingRoundTripper{
				rt: http.DefaultTransport,
			},
		}
	}

	mux.HandleFunc(AUTH_CODE_SUCCESS_CALLBACK_PATH, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic: %v", r)
				//log.Printf("Stack trace:\n%s", debug.Stack())

			}
		}()
		log.Println(r.RequestURI)

		authCode := r.URL.Query().Get("code")
		thisState := r.URL.Query().Get("state")
		if thisState != antiCsrfState() {
			log.Printf("state mismatch\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		log.Printf("authcode: %v\n", authCode)

		token, err := googleOauthConfig.Exchange(context.Background(), authCode)
		if err != nil {
			log.Printf("error: %v\n", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		log.Printf("token: %v\n", token)

		authHc := option.WithHTTPClient(hc)
		authService, err := oauth2.NewService(context.Background(), authHc, option.WithTokenSource(googleOauthConfig.TokenSource(context.Background(), token)))
		if err != nil {
			log.Printf("%#v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		userInfoService := oauth2.NewUserinfoService(authService)
		req := userInfoService.Get()
		userInfo, err := req.Do(googleapi.QueryParameter("access_token", token.AccessToken))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("%#v\n", userInfo.Id)

		infoB, err := json.Marshal(userInfo)
		if err != nil {
			log.Printf("%#v\n", err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Set-Cookie", COOKIE_NAME_GOOGLE_AUTHED_BEFORE+"=yes; path=/")
		//w.Header().Set("Set-Cookie", COOKIE_NAME_GOOGLE_AUTHED_BEFORE+"=yes")
		w.Write([]byte(fmt.Sprintf("%s", infoB)))
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie(COOKIE_NAME_GOOGLE_AUTHED_BEFORE)
		var cookieVal string
		if cookie != nil {
			cookieVal = cookie.Value
		}
		w.Write([]byte("pong. the client browser has cookie. " + COOKIE_NAME_GOOGLE_AUTHED_BEFORE + "=" + cookieVal))
	})

	go func() {
		log.Println("Listening on " + LISTENING_PORT)
		http.ListenAndServe(LISTENING_PORT, mux)
	}()
	recvSig := <-eCh
	log.Println("Received signal: " + recvSig.String() + " ; exiting...")
}

func Init() (err error) {
	err = InitConfig()
	if err != nil {
		return
	}
	return
}
