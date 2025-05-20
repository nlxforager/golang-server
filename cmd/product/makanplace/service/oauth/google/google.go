package google

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang-server/cmd/product/makanplace/config"

	oauth "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type Service struct {
	config *oauth.Config
	client *http.Client

	authCodeSuccessCallbackPath     string
	authCodeSuccessCallbackEndpoint string
}

func (s *Service) Exchange(ctx context.Context, code string, opts ...oauth.AuthCodeOption) (*oauth.Token, error) {
	return s.config.Exchange(ctx, code, opts...)
}

func (s *Service) AuthCodeURL(origin string) (string, error) {
	//return s.config.AuthCodeURL(s.antiCsrfState(), oauth.SetAuthURLParam("prompt", ""))
	csrf, err := s.setAntiCsrf(origin)
	if err != nil {
		return "", err
	}
	return s.config.AuthCodeURL(csrf, oauth.SetAuthURLParam("prompt", "consent select_account")), nil
}

func (s *Service) AuthCodeSuccessCallbackPath() string {
	return s.authCodeSuccessCallbackPath
}

var antCsrfOrigin = make(map[ /*csrf*/ string] /*origin*/ string)
var lock sync.RWMutex

func (s *Service) antiCsrfState(csrf string) (string, error) {
	lock.Lock()
	defer lock.Unlock()

	o, ok := antCsrfOrigin[csrf]
	log.Printf("antiCsrfState() %#v \n", antCsrfOrigin)

	if !ok {
		return "", errors.New("csrf not found in antiCsrfOrigin")
	}

	return o, nil
}

func (s *Service) setAntiCsrf(origin string) (string, error) {
	log.Printf("Setting antiCsrfOrigin\n")
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.New("anti-csrf generation failed")
	}
	antiCsrf := base64.RawURLEncoding.EncodeToString(b)
	log.Printf("Setting antiCsrfOrigin locking %s\n ", antiCsrf)

	lock.Lock()
	antCsrfOrigin[antiCsrf] = origin
	log.Printf("Setting antiCsrfOrigin unlocking\n ")

	lock.Unlock()
	log.Printf("Setting antiCsrfOrigin end %#v\n", antCsrfOrigin)

	return antiCsrf, nil
}

var ErrStateMismatch = errors.New("state mismatch")
var ErrExchangeTokenFailed = errors.New("exchange token failed")

func (s *Service) UserInfo(state string, authCode string) (_ *oauth2.Userinfo, origin string, _ error) {
	origin, err := s.antiCsrfState(state)
	if err != nil {
		return nil, "", ErrStateMismatch
	}
	log.Printf("[trying to get].. \n")
	{
		c := http.Client{}
		res, rErr := c.Do(&http.Request{
			Method: "GET",
			URL:    &url.URL{Scheme: "https", Host: "www.googleapis.com", Path: "/"},
		})
		if rErr == nil {
			log.Printf("[trying to get] want %d got %s \n", 404, res.Status)
		} else {
			log.Printf("[trying to get] err %s\n", rErr)
		}

	}
	log.Printf("[trying to get] completed \n")
	token, err := s.Exchange(context.Background(), authCode)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, "", ErrExchangeTokenFailed
	}
	log.Printf("token: %v\n", token)
	authHc := option.WithHTTPClient(s.client)
	authService, err := oauth2.NewService(context.Background(), authHc, option.WithTokenSource(s.config.TokenSource(context.Background(), token)))
	if err != nil {
		return nil, "", err
	}

	userInfoService := oauth2.NewUserinfoService(authService)
	req := userInfoService.Get()
	userInfo, err := req.Do(googleapi.QueryParameter("access_token", token.AccessToken))
	if err != nil {
		return nil, "", err
	}
	return userInfo, origin, nil
}

func (s *Service) FrontEndHomePageURL() string {
	return s.authCodeSuccessCallbackEndpoint
}

func NewService(c config.GoogleAuthConfig) Service {
	authCodeSuccessCallbackPath := c.AUTH_CODE_SUCCESS_ENDPOINT_PATH // to be binded with mux and used during config.Exchange.
	authCodeSuccessCallbackEndpoint := c.AUTH_CODE_SUCCESS_ENDPOINT_HOST
	var config = &oauth.Config{
		RedirectURL:  authCodeSuccessCallbackEndpoint + authCodeSuccessCallbackPath,
		ClientID:     c.CLIENT_ID_PREFIX + ".apps.googleusercontent.com",
		ClientSecret: c.CLIENT_SECRET,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	fmt.Printf("config %#v\n", config)

	hc := http.DefaultClient
	if c.ENABLE_LOG_REQUEST {
		hc = &http.Client{
			Transport: &LoggingRoundTripper{
				rt: http.DefaultTransport,
			},
		}
	}

	return Service{
		config:                          config,
		client:                          hc,
		authCodeSuccessCallbackPath:     authCodeSuccessCallbackPath,
		authCodeSuccessCallbackEndpoint: authCodeSuccessCallbackEndpoint,
	}
}

type LoggingRoundTripper struct {
	rt http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Log request
	fmt.Printf("[HTTP Request] %s %s\n", req.Method, req.URL.String())
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body
		fmt.Printf("[Request Body] %s\n", string(bodyBytes))
	}

	// Perform the actual request
	resp, err := lrt.rt.RoundTrip(req)
	if err != nil {
		fmt.Printf("[HTTP Error] %v\n", err)
		return nil, err
	}

	// Log response
	fmt.Printf("[HTTP Response] Status: %s in %v\n", resp.Status, time.Since(start))
	if resp.Body != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body
		fmt.Printf("[Response Body] %s\n", string(bodyBytes))
	}

	return resp, nil
}
