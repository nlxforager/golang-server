package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"golang-server/src/domain/auth"
	"golang-server/src/infrastructure/messaging/email"
	"golang-server/src/log"
)

type AuthHandler struct {
	AuthService auth.AuthService
	MailService email.EmailService
}

var _ auth.AuthService = &auth.MockAuth{}

func NewAuthHandler(authService auth.AuthService, mailService email.EmailService) (*AuthHandler, error) {
	if authService == nil {
		return nil, fmt.Errorf("auth service is required")
	}
	ah := &AuthHandler{
		AuthService: authService,
		MailService: mailService,
	}
	return ah, nil
}

func (h *AuthHandler) AuthByUsernamePassword() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "authService"))
	l.Info("auth by password")

	type RequestBody struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
	}

	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					form := &RequestBody{}
					json.NewDecoder(r.Body).Decode(form)

					var err error
					if form.Username == nil || form.Password == nil {
						err = fmt.Errorf("insufficent username or password")
					} else if err = h.AuthService.ByPasswordAndUsername(*form.Username, *form.Password); err != nil {
						err = fmt.Errorf("username and password invalid")
					}

					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write(AsError(err).ToBytes())
						return
					}

					otp := h.AuthService.OtpGen()
					err = h.AuthService.SetOTP(*form.Username, func() string {
						return otp
					})

					go func() {
						email, eerr := h.AuthService.GetEmail(*form.Username)
						if eerr == nil {
							h.MailService.SendOTP(email, otp)
						}
					}()

					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					token, err := h.AuthService.CreateTokenUsernameOnly(*form.Username)
					if err != nil {
						fmt.Printf("err %v", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					json.NewEncoder(w).Encode(struct {
						Data any `json:"data"`
					}{
						Data: struct {
							Username    *string `json:"username"`
							Token       string  `json:"token"`
							RedirectUrl string  `json:"redirect_url"`
						}{
							Username:    form.Username,
							RedirectUrl: "/otp_form/",
							Token:       token,
						},
					})
				},
			},
			DefaultFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte("{\"error\":\"Invalid Accept Header\"}"))
			},
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		accepts := r.Header["Accept"]
		options.GetAcceptFunc(accepts)(w, r)
	}
}

type OTP struct {
	Otp   *string `json:"otp"`
	Token *string `json:"token"`
}

func (h *AuthHandler) GetOtpForm() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "authService"))
	l.Info("AuthHandler::GetOtpForm")

	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(struct {
						Data any `json:"data"`
					}{
						Data: struct {
							SubmitUrl string `json:"submit_url"`
						}{
							SubmitUrl: "/otp/",
						},
					})
				},
			},
			DefaultFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte("{\"error\":\"Invalid Accept Header\"}"))
			},
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		accepts := r.Header["Accept"]
		options.GetAcceptFunc(accepts)(w, r)
	}
}

func (h *AuthHandler) SubmitOtp() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "authService"))
	l.Info("AuthHandler::SubmitOtp")

	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					form := &OTP{}
					json.NewDecoder(r.Body).Decode(form)

					var err error
					if form.Otp == nil || form.Token == nil {
						err = fmt.Errorf("insufficent otp or token")
					} else if err = h.AuthService.VerifyOTP(*form.Otp, *form.Token); err != nil {
						err = fmt.Errorf("username and token invalid")
					}
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write(AsError(err).ToBytes())
					}

					json.NewEncoder(w).Encode(struct {
						Data any `json:"data"`
					}{
						Data: struct {
							Token string `json:"token"`
						}{},
					})
				},
			},
			DefaultFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte("{\"error\":\"Invalid Accept Header\"}"))
			},
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		accepts := r.Header["Accept"]
		options.GetAcceptFunc(accepts)(w, r)
	}
}
