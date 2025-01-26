package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	emailservice "golang-server/cmd/servers/multifact/e2e_test/mocks/email"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authservice "golang-server/cmd/servers/multifact/e2e_test/mocks"
	"golang-server/cmd/servers/multifact/mux"
)

// "2FA_PW_E"
func TestHandler_Password_2FA_OK(t *testing.T) {
	mockAuthService := authservice.NewMockAuth()
	mockAuthService.UserByUsernames["user1"] = authservice.MockUser{
		Username: "user1",
		Password: "password1",
		Email:    "some.com.dummy",
		Mode:     "2FA_PW_E",
	}

	mockMailService := emailservice.NewMockOtpSingleSendReceiver()
	mux := mux.NewMux(&mux.MuxOpts{
		AuthMuxOpts: &mux.AuthMuxOpts{
			Auth: mockAuthService,
			Mail: mockMailService,
		},
	})

	type RegisterBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	registerBody, err := json.Marshal(&RegisterBody{
		Username: "user1",
		Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(registerBody)
	req := httptest.NewRequestWithContext(context.TODO(), http.MethodPost, "/register/", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}

	type PostTokenBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	postTokenBody, err := json.Marshal(&PostTokenBody{
		Username: "user1",
		Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}

	// FIXME
	// add routes for changing mode and email once completed
	user := mockAuthService.UserByUsernames["user1"]
	user.Mode = "2FA_PW_E"
	user.Email = "some@email"
	mockAuthService.UserByUsernames["user1"] = user

	reader = bytes.NewReader(postTokenBody)
	req = httptest.NewRequestWithContext(context.TODO(), http.MethodPost, "/token/", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")

	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	type response struct {
		Data  any    `json:"data"`
		Error string `json:"error"`
	}

	var resp response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		t.Errorf("expected error to be nil got %v, data %s", err, string(data))
	}

	_wantErr := ""
	if resp.Error != _wantErr {
		t.Errorf("expected %s got %v", _wantErr, resp.Error)
	}

	d, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Errorf("expected error to be nil got %#v, data %s", resp, string(data))
	}
	if d["username"] != "user1" {
		t.Errorf("expected username to be user1 got %v", d["username"])
	}
	if d["password"] != nil {
		t.Errorf("expected password to be password1 got %v", d["password"])
	}

	if d["weak_token"] == nil {
		t.Errorf("expected token got %v", d["weak_token"])
	}

	pwOkToken := d["weak_token"].(string)

	token, _ := jwt.Parse(pwOkToken, nil) // skip verify err since secret is server-side
	claims, _ := token.Claims.(jwt.MapClaims)

	_, ok = claims["is_auth"]
	if ok {
		t.Fatal("expected token to be authenticated")
	}

	if d["redirect_url"] == nil {
		t.Errorf("expected password to be redirect_url got %v", d["redirect_url"])
	}

	dd, ok := d["redirect_url"].(string)
	if !ok {
		t.Errorf("expected redirect_url to be nil got %v", d["redirect_url"])
	}
	if dd == "" {
		t.Errorf("expected redirect_url to be non-empty got %v", dd)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	otp, err := mockMailService.OTP(ctx)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	type SubmitOtpForm struct {
		Token string `json:"weak_token"`
		Otp   string `json:"otp"`
	}

	f := &SubmitOtpForm{
		Token: pwOkToken,
		Otp:   otp,
	}

	ff, eee := json.Marshal(f)
	if eee != nil {
		t.Errorf("expected error to be nil got %v", eee)
	}
	reader = bytes.NewReader(ff)
	otpSubmitReq := httptest.NewRequestWithContext(context.TODO(), http.MethodPost, dd, reader)
	otpSubmitReq.Header.Set("Accept", "application/json")

	w = httptest.NewRecorder()
	mux.ServeHTTP(w, otpSubmitReq)

	otpSubmitRes := w.Result()
	defer otpSubmitRes.Body.Close()

	type OtpResponse struct {
		Data  any    `json:"data"`
		Error string `json:"error"`
	}

	var v OtpResponse
	err = json.NewDecoder(otpSubmitRes.Body).Decode(&v)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	if otpSubmitRes.StatusCode != http.StatusOK {
		t.Errorf("otpSubmitRes expected status code to be %v got %v, body %#v", http.StatusOK, otpSubmitRes.StatusCode, v)
	}

	_wantErr = ""
	if resp.Error != _wantErr {
		t.Errorf("expected %s got %v", _wantErr, resp.Error)
	}

	d, ok = v.Data.(map[string]interface{})
	if !ok {
		t.Errorf("expected error to be nil got %#v, data %s", resp, string(data))
	}
	if d["weak_token"] != nil {
		t.Errorf("expected nil weak_token got %v", d["weak_token"])
	}
	if d["token"] == nil {
		t.Errorf("expected token got %v", d["token"])
	}

	otpOkToken, _ := d["token"].(string)
	token, _ = jwt.Parse(otpOkToken, nil)
	claims, _ = token.Claims.(jwt.MapClaims)
	is, ok := claims["is_auth"].(string)
	if !ok {
		t.Errorf("expected is_auth to be non-nil got %v", is)
	}
	if is != "true" {
		t.Errorf("expected is_auth to be true got %v", is)
	}
}

// "SIMPLE_PW"
func TestHandler_Password_Simple(t *testing.T) {
	type Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	b, err := json.Marshal(&Body{
		Username: "user1",
		Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(b)
	req := httptest.NewRequestWithContext(context.TODO(), http.MethodPost, "/token/", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	mockAuthService := authservice.NewMockAuth()
	mockAuthService.UserByUsernames["user1"] = authservice.MockUser{
		Username: "user1",
		Password: "password1",
		Email:    "some.com.dummy",
		Mode:     "SIMPLE_PW",
	}

	mockMailService := emailservice.NewMockOtpSingleSendReceiver()
	mux := mux.NewMux(&mux.MuxOpts{
		AuthMuxOpts: &mux.AuthMuxOpts{
			Auth: mockAuthService,
			Mail: mockMailService,
		},
	})
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	type response struct {
		Data  any    `json:"data"`
		Error string `json:"error"`
	}

	var resp response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		t.Errorf("expected error to be nil got %v, data %s", err, string(data))
	}

	_wantErr := ""
	if resp.Error != _wantErr {
		t.Errorf("expected %s got %v", _wantErr, resp.Error)
	}

	d, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Errorf("expected error to be nil got %#v, data %s", resp, string(data))
	}
	if d["username"] != "user1" {
		t.Errorf("expected username to be user1 got %v", d["username"])
	}
	if d["password"] != nil {
		t.Errorf("expected password to be empty got %v", d["password"])
	}

	if d["weak_token"] != nil {
		t.Errorf("expected nil weak_token got %v", d["weak_token"])
	}

	if d["token"] == nil {
		t.Errorf("expected token got %v", d["token"])

	}

	pwOkToken, _ := d["token"].(string)

	token, _ := jwt.Parse(pwOkToken, nil)
	claims, _ := token.Claims.(jwt.MapClaims)
	is, ok := claims["is_auth"].(string)
	if !ok {
		t.Errorf("expected is_auth to be non-nil got %v", is)
	}
	if is != "true" {
		t.Errorf("expected is_auth to be true got %v", is)
	}
}

func TestHandler_Password_NOTOK(t *testing.T) {
	type Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	b, err := json.Marshal(&Body{
		Username: "user1",
		Password: "passwoasdfrd1",
	})
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(b)
	req := httptest.NewRequestWithContext(context.TODO(), http.MethodPost, "/token/", reader)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	mockAuth := authservice.NewMockAuth()
	mockAuth.UserByUsernames["user1"] = authservice.MockUser{
		Username: "user1",
		Password: "corridged",
		Email:    "user1@example.com",
		Mode:     "SIMPLE_PW",
	}
	mockMailService := emailservice.NewMockOtpSingleSendReceiver()
	mux.NewMux(&mux.MuxOpts{
		AuthMuxOpts: &mux.AuthMuxOpts{
			Auth: mockAuth,
			Mail: mockMailService,
		},
	}).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	type response struct {
		Data  any    `json:"data"`
		Error string `json:"error"`
	}

	var resp response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		t.Errorf("expected error to be nil got %v, data %s", err, string(data))
	}

	_wantErr := ""
	if resp.Error == _wantErr {
		t.Errorf("expected error %s got %v", _wantErr, resp.Error)
	}

	d, ok := resp.Data.(map[string]interface{})
	if ok {
		t.Errorf("expected error to be nil got %#v, data %s", resp, string(data))
	}
	if d["username"] != nil {
		t.Errorf("expected username to be null got %v", d["username"])
	}
	if d["password"] != nil {
		t.Errorf("expected password to be null got %v", d["password"])
	}
	for _, v := range []string{"weak_token", "token"} {
		pwOkToken, ok := d[v].(string)
		if ok {
			t.Errorf("expected token to be null got %v", pwOkToken)
		}
	}
}

// "SIMPLE_PW" -> "2FA_PW_E"
// after changing, the first login attempt should obtain the `weak_token` only. that is, 2FA does not complete after 1FA
func TestHandler_PatchAuthMode(t *testing.T) {
	mockAuthService := authservice.NewMockAuth()
	mockAuthService.UserByUsernames["user1"] = authservice.MockUser{
		Username: "user1",
		Password: "password1",
		Email:    "some.com.dummy",
		Mode:     "SIMPLE_PW",
	}

	mockMailService := emailservice.NewMockOtpSingleSendReceiver()
	mux := mux.NewMux(&mux.MuxOpts{
		AuthMuxOpts: &mux.AuthMuxOpts{
			Auth: mockAuthService,
			Mail: mockMailService,
		},
	})

	type Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	b, err := json.Marshal(&Body{
		Username: "user1",
		Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(b)
	req := httptest.NewRequestWithContext(context.TODO(), http.MethodPost, "/register/", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status code to be %v got %v", http.StatusCreated, res.StatusCode)
	}

	// login

	b, err = json.Marshal(&Body{
		Username: "user1",
		Password: "password1",
	})
	if err != nil {
		t.Fatal(err)
	}
	reader = bytes.NewReader(b)

	req = httptest.NewRequestWithContext(context.TODO(), http.MethodPost, "/token/", reader)
	req.Header.Set("Accept", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	res = w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	//
	type RespBody struct {
		Data any `json:"data"`
	}
	var respbody RespBody
	err = json.Unmarshal(resp, &respbody)
	if err != nil {
		t.Fatalf("expected error to be nil got %v", err)
	}
	d, ok := respbody.Data.(map[string]interface{})
	if d["token"] == nil {
		t.Errorf("expected token got %v", d["token"])
	}

	token, ok := d["token"].(string)
	if !ok {
		t.Fatalf("expected token to be non-nil got %v", d["token"])
	}

	{ // PATCH
		type Patch struct {
			Op       string `json:"op"`
			Username string `json:"username"`
			Mode     string `json:"auth_mode"`
			Email    string `json:"email"`
		}
		type PatchBody []Patch

		patches := PatchBody{
			Patch{
				Op:       "modify",
				Username: "user1",
				Mode:     "2FA_PW_E",
				Email:    "dummy@some.com",
			},
		}
		bodypatch, _ := json.Marshal(&patches)
		reader = bytes.NewReader(bodypatch)
		req = httptest.NewRequestWithContext(context.TODO(), http.MethodPatch, "/user/", reader)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		res = w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
		}
	}

	// AFTER CHANGE
	reader = bytes.NewReader(b)
	req = httptest.NewRequestWithContext(context.TODO(), http.MethodPost, "/token/", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")

	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	type response struct {
		Data  any    `json:"data"`
		Error string `json:"error"`
	}

	var firstFAresp response
	err = json.Unmarshal(data, &firstFAresp)
	if err != nil {
		t.Errorf("expected error to be nil got %v, data %s", err, string(data))
	}

	_wantErr := ""
	if firstFAresp.Error != _wantErr {
		t.Errorf("expected %s got %v", _wantErr, firstFAresp.Error)
	}

	d, ok = firstFAresp.Data.(map[string]interface{})
	if !ok {
		t.Errorf("expected error to be nil got %#v, data %s", resp, string(data))
	}
	if d["username"] != "user1" {
		t.Errorf("expected username to be user1 got %v", d["username"])
	}
	if d["password"] != nil {
		t.Errorf("expected password to be password1 got %v", d["password"])
	}

	if d["weak_token"] == nil {
		t.Errorf("expected weak_token %v", d["weak_token"])
	}

	if d["token"] != nil {
		t.Errorf("expected nil token got %v", d["token"])
	}
}
