package mk_outlet_service

import (
	"encoding/json"
	"errors"
	"log"

	"golang-server/cmd/product/makanplace/client/onemap"
	outletrepo "golang-server/cmd/product/makanplace/repositories/outlet"
)

type AddOutletBody struct {
	OutletName    string
	OutletType    string
	MenuItems     []string
	Address       string
	PostalCode    string
	OfficialLinks []string
	ReviewLinks   []outletrepo.ReviewLink
}

type PutOutletBody struct {
	Id            *int64
	OutletName    string
	OutletType    string
	ProductNames  []string
	Address       string
	PostalCode    string
	OfficialLinks []string
	ReviewLinks   []outletrepo.ReviewLink
}

type LatLongGetter interface {
	GetLetLong() (*LatLong, error)
}

type Service struct {
	repo *outletrepo.Repo
}

func (s *Service) AddOutlet(b AddOutletBody) (int64, error) {
	return s.repo.NewOutletWithMenu(b.OutletName, b.Address, b.PostalCode, b.OfficialLinks, b.ReviewLinks, b.MenuItems)
}

func (s *Service) PutOutlet(b PutOutletBody) error {
	if b.Id == nil {
		return errors.New("id is required")
	}
	return s.repo.UpdateOutletWithMenu(b.Id, b.OutletName, b.Address, b.PostalCode, b.OfficialLinks, b.ReviewLinks, b.ProductNames)
}

type LatLong = onemap.LatLong
type MenuItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Outlet struct {
	Name          string
	Address       string
	PostalCode    string
	OfficialLinks []string
	ReviewLinks   []outletrepo.ReviewLink

	*LatLong
	Id        int64
	MenuItems []MenuItem `json:"menu_items"`
}

func (s *Service) GetOutlets(postalCode *string, id *int) ([]Outlet, error) {
	outletsDb, err := s.repo.GetOutlets(postalCode, id)
	if err != nil {
		return []Outlet{}, err
	}

	var outlets []Outlet
	c := onemap.NewClient()

	for _, outlet := range outletsDb {
		var latlong *LatLong
		if len(outlet.LatLong) != 2 {
			_latlong, lErr := c.GetLatLong(outlet.PostalCode)
			if lErr != nil || _latlong == nil {
				log.Printf("Error getting latlong or nil latlong: %v\n", lErr)
			} else {
				sErr := s.repo.SetLatLong(outlet.PostalCode, _latlong.Latitude, _latlong.Longitude)
				if sErr != nil {
					log.Printf("Error setting latlong: %v\n", sErr)
				}
			}
		} else {
			latlong = &LatLong{
				Latitude:  outlet.LatLong[0],
				Longitude: outlet.LatLong[1],
			}
		}

		var mi []MenuItem
		json.Unmarshal(outlet.MenuItems, &mi)

		var rl []outletrepo.ReviewLink

		json.Unmarshal(outlet.ReviewLinks, &rl)

		o := Outlet{
			Name:          outlet.Name,
			Address:       outlet.Address,
			PostalCode:    outlet.PostalCode,
			OfficialLinks: outlet.OfficialLinks,
			LatLong:       latlong,
			Id:            outlet.Id,
			MenuItems:     mi,
			ReviewLinks:   rl,
		}

		outlets = append(outlets, o)
	}
	return outlets, nil
}

func NewOutletService(repo *outletrepo.Repo) *Service {
	return &Service{repo: repo}
}
