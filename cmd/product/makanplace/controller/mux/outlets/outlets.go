package outlets

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"

	log2 "golang-server/cmd/product/makanplace/httplog"

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
	ReviewLinks   []Link `json:"review_links"`
}

func Register(mux *http.ServeMux, makanTokenCookieKey string, mkService *mk_user_session.Service, mws middlewares.MiddewareStack, outletService *mk_outlet_service.Service) {
	// middleware: isSuperUser
	mwsWithSuper := mws.Wrap(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionId := middlewares.GetSessionIdFromRequest(r)
			session := mkService.GetSession(sessionId, false)
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
		prefix := fmt.Sprintf("%s [POST /outlet]\n", log2.SPrintHttpRequestPrefix(r))

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

		var reviewLinks []string
		for _, link := range b.ReviewLinks {
			reviewLinks = append(reviewLinks, link.Value)
		}
		err = outletService.AddOutlet(mk_outlet_service.ServiceBody{
			OutletName:    b.OutletName,
			OutletType:    b.OutletType,
			ProductName:   b.ProductName,
			Address:       b.Address,
			PostalCode:    b.PostalCode,
			OfficialLinks: officialLinks,
			ReviewLinks:   reviewLinks,
		})

		if err != nil {
			log.Printf("%s Error adding outlet: %v\n", prefix, err)
			response_types.ErrorNoBody(w, http.StatusBadRequest, err)
			return
		}

		response_types.OkEmptyJsonBody(w)
	})))

	mux.Handle("GET /outlets/", mws.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("outlets getting")
		cookie, _ := r.Cookie(makanTokenCookieKey)
		session := mkService.GetSession(cookie.Value, false)
		if session == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		outletsS, err := outletService.GetOutlets()
		if err != nil {
			log.Printf("Error getting outlets: %v\n", err)
			response_types.ErrorNoBody(w, http.StatusInternalServerError, err)
			return
		}

		var out []Outlet
		for _, o := range outletsS {
			out = append(out, Outlet{
				Name:          o.Name,
				Address:       o.Address,
				PostalCode:    o.PostalCode,
				OfficialLinks: o.OfficialLinks,
				ReviewLinks:   o.ReviewLinks,
				LatLong:       o.LatLong,
			})
		}

		log.Printf("outlets got %d\n", len(out))
		response_types.OkJsonBody(w, struct {
			Outlets []Outlet `json:"outlets"`
		}{out})
	})))
}

type Outlet struct {
	Name                       string   `json:"name"`
	Address                    string   `json:"address"`
	PostalCode                 string   `json:"postal_code"`
	OfficialLinks              []string `json:"official_links"`
	*mk_outlet_service.LatLong `json:"latlong"`
	ReviewLinks                []string `json:"review_links"`
}
