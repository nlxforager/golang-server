package outlets

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"

	log2 "golang-server/cmd/product/makanplace/log"

	"golang-server/cmd/product/makanplace/controller/middlewares"
	"golang-server/cmd/product/makanplace/controller/response_types"
	"golang-server/cmd/product/makanplace/service/mk_outlet"
	"golang-server/cmd/product/makanplace/service/mk_user_session"
)

type Response struct {
	LoginUrls map[string]string `json:"login_urls"`

	UserInfo *mk_user_session.UserInfo `json:"user_info"`
}

type Link struct {
	Value string `json:"value"`
}
type Body struct {
	OutletName    string `json:"outlet_name"`
	OutletType    string `json:"outlet_type"`
	ProductName   string `json:"product_name"`
	Address       string `json:"address"`
	PostalCode    string `json:"postal_code"`
	OfficialLinks []Link `json:"official_links"`
}

func Register(mux *http.ServeMux, makanTokenCookieKey string, mkService *mk_user_session.Service, middlewares middlewares.MiddewareStack, outletService *mk_outlet_service.OutletService) {
	// middleware: isSuperUser
	mwsWithSuper := middlewares.Wrap(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, _ := r.Cookie(makanTokenCookieKey)
			session := mkService.GetSession(cookie.Value, false)
			log.Printf("checking IsSuperUser: %#v\n", session)

			if len(session.Gmails) == 0 {
				log.Printf("not IsSuperUser: no gmails to compare\n")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if slices.ContainsFunc(session.Gmails, func(gmail string) bool {
				return !mkService.IsSuperUser(gmail)
			}) {
				log.Printf("not IsSuperUser: some gmail is not super.\n")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			handler.ServeHTTP(w, r)
		})
	})

	mux.Handle("POST /outlet/", mwsWithSuper.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := fmt.Sprintf("%s [POST /outlet]\n", log2.HttpRequestPrefix(r.Context()))
		cookie, _ := r.Cookie(makanTokenCookieKey)
		session := mkService.GetSession(cookie.Value, false)
		if session == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var b Body
		err = json.Unmarshal(body, &b)

		var officialLinks []string
		for _, link := range b.OfficialLinks {
			officialLinks = append(officialLinks, link.Value)
		}
		err = outletService.AddOutlet(mk_outlet_service.ServiceBody{
			OutletName:    b.OutletName,
			OutletType:    b.OutletType,
			ProductName:   b.ProductName,
			Address:       b.Address,
			PostalCode:    b.PostalCode,
			OfficialLinks: officialLinks,
		})

		if err != nil {
			log.Printf("%s Error adding outlet: %v\n", prefix, err)
			response_types.ErrorNoBody(w, http.StatusBadRequest, err)
			return
		}

		response_types.OkEmptyJsonBody(w)
	})))

	mux.Handle("GET /outlet/", middlewares.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("outlets getting")
		cookie, _ := r.Cookie(makanTokenCookieKey)
		session := mkService.GetSession(cookie.Value, false)
		if session == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("session found %#v\n", session)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
		w.WriteHeader(http.StatusOK)
	})))
}
