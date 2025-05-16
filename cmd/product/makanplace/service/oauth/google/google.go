package google

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
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

	authCodeSuccessCallbackPath string
}

func (s *Service) Exchange(ctx context.Context, code string, opts ...oauth.AuthCodeOption) (*oauth.Token, error) {
	return s.config.Exchange(ctx, code, opts...)
}

func (s *Service) AuthCodeURL() string {
	return s.config.AuthCodeURL(s.antiCsrfState(), oauth.SetAuthURLParam("prompt", ""))
	return s.config.AuthCodeURL(s.antiCsrfState(), oauth.SetAuthURLParam("prompt", "consent select_account"))
}

func (s *Service) AuthCodeSuccessCallbackPath() string {
	return s.authCodeSuccessCallbackPath
}

func (s *Service) antiCsrfState() string {
	return "SOME_STATE"
}

var ErrStateMismatch = errors.New("state mismatch")
var ErrExchangeTokenFailed = errors.New("exchange token failed")

func (s *Service) UserInfo(state string, authCode string) (*oauth2.Userinfo, error) {
	if state != s.antiCsrfState() {
		return nil, ErrStateMismatch
	}

	token, err := s.Exchange(context.Background(), authCode)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, ErrExchangeTokenFailed
	}
	log.Printf("token: %v\n", token)
	authHc := option.WithHTTPClient(s.client)
	authService, err := oauth2.NewService(context.Background(), authHc, option.WithTokenSource(s.config.TokenSource(context.Background(), token)))
	if err != nil {
		return nil, err
	}

	userInfoService := oauth2.NewUserinfoService(authService)
	req := userInfoService.Get()
	userInfo, err := req.Do(googleapi.QueryParameter("access_token", token.AccessToken))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (s *Service) FrontEndHomePageURL() string {
	return "http://localhost:5173"
}

func NewService(c config.GoogleAuthConfig) Service {
	authCodeSuccessCallbackPath := c.AUTH_CODE_SUCCESS_ENDPOINT_PATH // to be binded with mux and used during config.Exchange.

	var config = &oauth.Config{
		RedirectURL:  "http://localhost" + c.Port + authCodeSuccessCallbackPath,
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
		config:                      config,
		client:                      hc,
		authCodeSuccessCallbackPath: authCodeSuccessCallbackPath,
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
