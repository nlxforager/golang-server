package ping

import (
	"encoding/json"
	"golang-server/cmd/product/makanplace/httplog"
	"log"
	"net/http"

	"golang-server/cmd/product/makanplace/controller/middlewares"
	"golang-server/cmd/product/makanplace/service/mk_user_session"

	"golang-server/cmd/product/makanplace/repositories/auth"
)

type OutletForm struct {
	ProductName []string `json:"product_names"`
}
type Response struct {
	LoginUrls map[string]string `json:"login_urls"`

	UserInfo *auth.UserWithGmail `json:"user_info"`

	OutletForm `json:"outlet_form"`
}

func Register(mux *http.ServeMux, mkService *mk_user_session.Service, goauthloginurl string, mws middlewares.MiddewareStack) {
	// Checks current user state of the client
	// Provides server configuration values

	mux.Handle("/ping", mws.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId := middlewares.GetSessionFromRequest(r)
		log.Printf("%s sessionId: %s\n", httplog.SPrintHttpRequestPrefix(r), sessionId)

		resp := Response{LoginUrls: make(map[string]string)}
		resp.LoginUrls["google"] = goauthloginurl

		session := mkService.GetSession(sessionId, false)
		resp.UserInfo = session

		resp.OutletForm = OutletForm{
			ProductName: []string{"Fried Hokkien Mee"},
		}
		respB, _ := json.Marshal(resp)

		log.Printf("%s response: %s\n", httplog.SPrintHttpRequestPrefix(r), string(respB))
		w.Write(respB)
	})))

}
