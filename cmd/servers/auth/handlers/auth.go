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

type AUTH_MODE string

const AUTH_MODE_SIMPLE_PW AUTH_MODE = "SIMPLE_PW"
const AUTH_MODE_2FA_PW_E AUTH_MODE = "2FA_PW_E"

func (h *AuthHandler) AuthByUsernamePassword() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "AuthHandler"))
	l.Info("AuthHandler::AuthByUsernamePassword")

	type RequestBody struct {
		Username *string    `json:"username"`
		Password *string    `json:"password"`
		Mode     *AUTH_MODE `json:"auth_mode"`
	}

	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					form := &RequestBody{}
					json.NewDecoder(r.Body).Decode(form)

					var err error
					var errStatusCode int
					switch {
					case form.Mode == nil:
						err = fmt.Errorf("unknown authentication mode")
						errStatusCode = http.StatusBadRequest
					case form.Username == nil || form.Password == nil:
						err = fmt.Errorf("insufficent username or password")
						errStatusCode = http.StatusBadRequest
					default:
						err = h.AuthService.ByPasswordAndUsername(*form.Username, *form.Password)
						errStatusCode = http.StatusUnauthorized
					}

					if err != nil {
						w.WriteHeader(errStatusCode)
						w.Write(AsError(err).ToBytes())
						return
					}

					switch *form.Mode {
					case AUTH_MODE_SIMPLE_PW:
						token, err := h.AuthService.CreateTokenUsernameOnly(*form.Username)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
						json.NewEncoder(w).Encode(struct {
							Data any `json:"data"`
						}{
							Data: struct {
								Username *string `json:"username"`
								Token    string  `json:"token"`
							}{
								Username: form.Username,
								Token:    token,
							},
						})
					case AUTH_MODE_2FA_PW_E:
						otp := h.AuthService.OtpGen()
						err = h.AuthService.SetOTP(*form.Username, func() string {
							return otp
						})
						email_, err := h.AuthService.GetEmail(*form.Username)

						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							w.Write(AsError(err).ToBytes())
							return
						}

						go h.MailService.SendOTP(email_, otp)

						weakToken, err := h.AuthService.CreateTokenUsernameOnly(*form.Username)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
						json.NewEncoder(w).Encode(struct {
							Data any `json:"data"`
						}{
							Data: struct {
								Username    *string `json:"username"`
								WeakToken   string  `json:"weak_token"`
								RedirectUrl string  `json:"redirect_url"`
							}{
								Username:    form.Username,
								RedirectUrl: "/otp/",
								WeakToken:   weakToken,
							},
						})
					}

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
	Token *string `json:"weak_token"`
}

func (h *AuthHandler) SubmitOtp() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "authService"))
	l.Info("AuthHandler::SubmitOtp")

	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					form := &OTP{}

					var err error
					if err = json.NewDecoder(r.Body).Decode(form); err != nil {
						goto prevalidation
					}
					if form.Otp == nil || form.Token == nil {
						err = fmt.Errorf("insufficent otp or token")
						goto prevalidation
					}
					if err = h.AuthService.VerifyOTP(*form.Otp, *form.Token); err != nil {
						err = fmt.Errorf("otp and token invalid")
					}

				prevalidation:
					{
						if err != nil {
							w.WriteHeader(http.StatusUnauthorized)
							w.Write(AsError(err).ToBytes())
						}
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

//func (h *AuthHandler) GetOtpForm() func(w http.ResponseWriter, r *http.Request) {
//	l := log.Logger.With(slog.String("handler", "authService"))
//	l.Info("AuthHandler::GetOtpForm")
//
//	options := Options{
//		AcceptFuncsOpts: AcceptFuncsOpts{
//			AcceptFuncs: map[string]AcceptFunc{
//				"application/json": func(w http.ResponseWriter, r *http.Request) {
//					json.NewEncoder(w).Encode(struct {
//						Data any `json:"data"`
//					}{
//						Data: struct {
//							SubmitUrl string `json:"submit_url"`
//						}{
//							SubmitUrl: "/otp/",
//						},
//					})
//				},
//			},
//			DefaultFunc: func(w http.ResponseWriter, r *http.Request) {
//				w.WriteHeader(http.StatusNotAcceptable)
//				w.Write([]byte("{\"error\":\"Invalid Accept Header\"}"))
//			},
//		},
//	}
//	return func(w http.ResponseWriter, r *http.Request) {
//		accepts := r.Header["Accept"]
//		options.GetAcceptFunc(accepts)(w, r)
//	}
//}
