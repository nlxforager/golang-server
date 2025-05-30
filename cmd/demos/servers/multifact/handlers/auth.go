package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"golang-server/cmd/demos/servers/multifact/e2e_test/mocks"
	"golang-server/src/domain/auth"
	"golang-server/src/domain/email"
	"golang-server/src/log"
)

type AuthHandler struct {
	AuthService auth.Authenticator
	OtpEmailer  email.OTPEmailer
}

var _ auth.Authenticator = &mocks.MockAuth{}

func NewAuthHandler(authService auth.Authenticator, mailService email.OTPEmailer) (*AuthHandler, error) {
	if authService == nil {
		return nil, fmt.Errorf("auth service is required")
	}
	if mailService == nil {
		return nil, fmt.Errorf("mail service is required")
	}
	ah := &AuthHandler{
		AuthService: authService,
		OtpEmailer:  mailService,
	}
	return ah, nil
}

type RequestBody struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (h *AuthHandler) RegisterUsernamePassword() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "AuthHandler"))
	l.Info("AuthHandler::RegisterUsernamePassword")

	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					form := &RequestBody{}
					json.NewDecoder(r.Body).Decode(form)

					var err error
					var errStatusCode int
					switch {
					case form.Username == nil || form.Password == nil:
						err = fmt.Errorf("insufficent username or password")
						errStatusCode = http.StatusBadRequest
					default:
						err = h.AuthService.RegisterUsernamePassword(*form.Username, *form.Password)
						errStatusCode = http.StatusInternalServerError
					}

					if err != nil {
						w.WriteHeader(errStatusCode)
						w.Write(AsError(err).ToBytes())
						return
					}
					w.WriteHeader(http.StatusCreated)
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

type AuthUsernameResourceBody struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}
type TokenResponseData struct {
	Username    *string `json:"username"`
	WeakToken   string  `json:"weak_token"`
	RedirectUrl string  `json:"redirect_url"`
}
type TokenResponseBody_200 struct {
	Data TokenResponseData `json:"data"`
}

func (h *AuthHandler) AuthByUsernamePassword() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "AuthHandler"))
	l.Info("AuthHandler::AuthByUsernamePassword")

	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					form := &AuthUsernameResourceBody{}
					json.NewDecoder(r.Body).Decode(form)

					var err error
					var errStatusCode int
					var user auth.User
					switch {

					case form.Username == nil || form.Password == nil:
						err = fmt.Errorf("insufficent username or password %v", form)
						errStatusCode = http.StatusBadRequest
					default:
						_err, _user := h.AuthService.ByPasswordAndUsername(*form.Username, *form.Password)
						err = _err
						if _user != nil {
							user = *_user
						}
						errStatusCode = http.StatusUnauthorized
					}

					if err != nil {
						w.WriteHeader(errStatusCode)
						w.Write(AsError(err).ToBytes())
						return
					}
					switch user.AuthMode {
					case auth.AUTH_MODE_SIMPLE_PW:
						token, err := h.AuthService.CreateStrongToken(*form.Username, user.AuthMode)
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
					case auth.AUTH_MODE_2FA_PW_E:
						otpGen := h.AuthService.OtpGen()
						stri, err := h.AuthService.SetOTP(*form.Username, otpGen)
						email_, err := h.AuthService.GetEmail(*form.Username)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							w.Write(AsError(err).ToBytes())
							return
						}

						go h.OtpEmailer.SendOTP(email_, stri)

						weakToken, err := h.AuthService.CreateWeakToken(*form.Username, user.AuthMode)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						json.NewEncoder(w).Encode(TokenResponseBody_200{
							Data: TokenResponseData{
								Username:    form.Username,
								RedirectUrl: "/otp/",
								WeakToken:   weakToken,
							},
						})
					default:
						w.WriteHeader(http.StatusBadRequest)
						w.Write(AsError(fmt.Errorf("unknown authentication mode")).ToBytes())
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

type SubmitOtpResponseBody struct {
	Data SubmitOtpResponseData
}
type SubmitOtpResponseData struct {
	Token string `json:"token"`
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

					var claims map[string]string
					var token string

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

					claims, err = h.AuthService.ValidateAndGetClaims(*form.Token)
					if err != nil {
						goto prevalidation
					}
					token, err = h.AuthService.CreateStrongToken(claims["sub"], auth.AUTH_MODE(claims["auth_mode"]))
					if err != nil {
						goto prevalidation
					}
				prevalidation:
					{
						if err != nil {
							w.WriteHeader(http.StatusUnauthorized)
							w.Write(AsError(err).ToBytes())
						}
					}

					json.NewEncoder(w).Encode(SubmitOtpResponseBody{
						Data: SubmitOtpResponseData{
							Token: token,
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

type Patch struct {
	Op       string          `json:"op"`
	Username string          `json:"username"`
	Mode     *auth.AUTH_MODE `json:"auth_mode"`
	Email    *string         `json:"email"`
}
type PatchRequestBody []Patch

func (h *AuthHandler) PatchUser() func(w http.ResponseWriter, r *http.Request) {
	l := log.Logger.With(slog.String("handler", "authService"))
	l.Info("AuthHandler::PatchUser")

	options := Options{
		AcceptFuncsOpts: AcceptFuncsOpts{
			AcceptFuncs: map[string]AcceptFunc{
				"application/json": func(w http.ResponseWriter, r *http.Request) {
					l.Info("AuthHandler::PatchUser()")
					patches := PatchRequestBody{}
					var err error
					if err = json.NewDecoder(r.Body).Decode(&patches); err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					if len(patches) != 1 {
						err = fmt.Errorf("request can have one and only one patch")
						w.WriteHeader(http.StatusBadRequest)
						w.Write(AsError(err).ToBytes())
						return
					}
					p := patches[0]
					if p.Op == "modify" {
						if !p.Mode.IsValid() {
							w.WriteHeader(http.StatusBadRequest)
							return
						}

						err := h.AuthService.ModifyUser(p.Username, auth.ChangeSet{
							AuthMode: p.Mode,
							Email:    p.Email,
						})
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							w.Write(AsError(err).ToBytes())
							return
						}
						w.WriteHeader(http.StatusOK)
						return
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

//func (h *AuthHandler) GetOtpForm() func(w http.ResponseWriter, r *http.Request) {
//	l := log.Logger.With(slog.String("handler", "authService"))
//	l.Info("AuthHandler::GetOtpForm")
//
//	options := Options{
//		AcceptFuncsOpts: AcceptFuncsOpts{
//			AcceptFuncs: map[string]AcceptFunc{
//				"application/json": func(w http.ResponseWriter, r *http.Request) {
//					json.NewEncoder(w).Encode(struct {
//						TokenResponseBody_200 any `json:"data"`
//					}{
//						TokenResponseBody_200: struct {
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
