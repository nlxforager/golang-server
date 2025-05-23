package outlets

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	mklog "golang-server/cmd/product/makanplace/httplog"

	outletrepo "golang-server/cmd/product/makanplace/repositories/outlet"

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
	Link     string `json:"link"`
	Platform string `json:"platform"`
	Creator  string `json:"creator"`
}

type PostOutletBody struct {
	OutletName    string   `json:"name"`
	OutletType    string   `json:"outlet_type"`
	Address       string   `json:"address"`
	PostalCode    string   `json:"postal_code"`
	OfficialLinks []Link   `json:"official_links"`
	ReviewLinks   []Link   `json:"review_links"`
	MenuItems     []string `json:"menu"`
}

type PutOutletBody struct {
	Id            *int64   `json:"id"`
	OutletName    string   `json:"name"`
	OutletType    string   `json:"outlet_type"`
	ProductName   string   `json:"product_name"`
	Address       string   `json:"address"`
	PostalCode    string   `json:"postal_code"`
	OfficialLinks []Link   `json:"official_links"`
	ReviewLinks   []Link   `json:"review_links"`
	ProductNames  []string `json:"menu"`
}

func Register(mux *http.ServeMux, mkService *mk_user_session.Service, mws middlewares.MiddewareStack, outletService *mk_outlet_service.Service) {
	// middleware: isSuperUser
	mwsWithSuper := mws.Wrap(middlewares.SuperUserMiddleware(mkService))

	// create outlet
	mux.Handle("POST /outlet/", mwsWithSuper.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := fmt.Sprintf("%s [POST /outlet]\n", mklog.SPrintHttpRequestPrefix(r))

		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var b PostOutletBody
		err = json.Unmarshal(body, &b)

		var officialLinks []string
		for _, link := range b.OfficialLinks {
			officialLinks = append(officialLinks, link.Link)
		}

		var reviewLinks []outletrepo.ReviewLink
		for _, link := range b.ReviewLinks {
			reviewLinks = append(reviewLinks, outletrepo.ReviewLink{
				Link:     link.Link,
				Platform: link.Platform,
				Creator:  link.Creator,
			})
		}
		outletId, err := outletService.AddOutlet(mk_outlet_service.AddOutletBody{
			OutletName:    b.OutletName,
			OutletType:    b.OutletType,
			MenuItems:     b.MenuItems,
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

		response_types.OkJsonBody(w, struct {
			Id int64 `json:"id"`
		}{
			Id: outletId,
		})
	})))

	// Put outlet
	mux.Handle("PUT /outlet/", mwsWithSuper.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := fmt.Sprintf("%s [PUT /outlet]\n", mklog.SPrintHttpRequestPrefix(r))

		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var b PutOutletBody
		err = json.Unmarshal(body, &b)

		var officialLinks []string
		for _, link := range b.OfficialLinks {
			officialLinks = append(officialLinks, link.Link)
		}

		var reviewLinks []outletrepo.ReviewLink
		for _, link := range b.ReviewLinks {
			reviewLinks = append(reviewLinks, outletrepo.ReviewLink{
				Link:     link.Link,
				Platform: link.Platform,
				Creator:  link.Creator,
			})
		}
		err = outletService.PutOutlet(mk_outlet_service.PutOutletBody{
			Id:            b.Id,
			OutletName:    b.OutletName,
			OutletType:    b.OutletType,
			ProductNames:  b.ProductNames,
			Address:       b.Address,
			PostalCode:    b.PostalCode,
			OfficialLinks: officialLinks,
			ReviewLinks:   reviewLinks,
		})

		if err != nil {
			log.Printf("%s Error editing outlet: %v\n", prefix, err)
			response_types.ErrorNoBody(w, http.StatusBadRequest, err)
			return
		}

		response_types.OkEmptyJsonBody(w)
	})))

	// get outlets
	mux.Handle("GET /outlets", mws.Finalize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId := middlewares.GetSessionIdFromRequest(r)
		session := mkService.GetSession(sessionId, false)
		if session == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		prefix := mklog.SPrintHttpRequestPrefix(r)

		postalCode := r.URL.Query().Get("postal_code")
		_id := r.URL.Query().Get("id")

		var postCodeParam *string
		if postalCode != "" {
			postCodeParam = &postalCode
		}

		var id *int
		if _id != "" {
			_idInt, err := strconv.Atoi(_id)
			if err != nil {
				response_types.ErrorNoBody(w, http.StatusBadRequest, err)
			}
			id = &_idInt
		}

		outletsS, err := outletService.GetOutlets(postCodeParam, id)

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
				ReviewLinks:   toControllerType(o.ReviewLinks),
				LatLong:       o.LatLong,
				Id:            o.Id,
				MenuItem:      o.MenuItems,
			})
		}

		log.Printf("%s outlets got %d\n", prefix, len(out))
		response_types.OkJsonBody(w, struct {
			Outlets []Outlet `json:"outlets"`
		}{out})
	})))
}
