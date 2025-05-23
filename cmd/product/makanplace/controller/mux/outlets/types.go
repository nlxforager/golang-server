package outlets

import (
	outletrepo "golang-server/cmd/product/makanplace/repositories/outlet"
	mk_outlet_service "golang-server/cmd/product/makanplace/service/mk_outlet"
)

type ReviewLink struct {
	Link     string `json:"link"`
	Platform string `json:"platform"`
	Creator  string `json:"creator"`
}

func toControllerType(links []outletrepo.ReviewLink) (v []ReviewLink) {
	for _, link := range links {
		v = append(v, ReviewLink(link))
	}
	return
}

type Outlet struct {
	Name                       string   `json:"name"`
	Address                    string   `json:"address"`
	PostalCode                 string   `json:"postal_code"`
	OfficialLinks              []string `json:"official_links"`
	*mk_outlet_service.LatLong `json:"latlong"`
	ReviewLinks                []ReviewLink                 `json:"review_links"`
	Id                         int64                        `json:"id"`
	MenuItem                   []mk_outlet_service.MenuItem `json:"menu"`
}
