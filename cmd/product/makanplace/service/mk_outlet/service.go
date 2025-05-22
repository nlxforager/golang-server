package mk_outlet_service

import (
	"errors"
	"log"

	"golang-server/cmd/product/makanplace/client/onemap"
	"golang-server/cmd/product/makanplace/repositories/outlet"
)

type AddOutletBody struct {
	OutletName    string
	OutletType    string
	ProductName   string
	Address       string
	PostalCode    string
	OfficialLinks []string
	ReviewLinks   []string
}

type PutOutletBody struct {
	Id            *int64
	OutletName    string
	OutletType    string
	ProductName   string
	Address       string
	PostalCode    string
	OfficialLinks []string
	ReviewLinks   []string
}

type LatLongGetter interface {
	GetLetLong() (*LatLong, error)
}

type Service struct {
	repo *outlet.Repo
}

func (s *Service) AddOutlet(b AddOutletBody) error {
	return s.repo.NewOutletWithMenu(b.OutletName, b.Address, b.PostalCode, b.OfficialLinks, b.ReviewLinks, []string{b.ProductName})
}

func (s *Service) PutOutlet(b PutOutletBody) error {
	if b.Id == nil {
		return errors.New("id is required")
	}
	return s.repo.UpdateOutletWithMenu(b.Id, b.OutletName, b.Address, b.PostalCode, b.OfficialLinks, b.ReviewLinks, []string{b.ProductName})
}

type LatLong = onemap.LatLong
type Outlet struct {
	Name          string
	Address       string
	PostalCode    string
	OfficialLinks []string
	ReviewLinks   []string

	*LatLong
	Id int64
}

func (s *Service) GetOutlets(postalCode *string, id *int) ([]Outlet, error) {
	outletsDb, err := s.repo.GetOutlets(postalCode, id)
	if err != nil {
		return []Outlet{}, err
	}

	var outlets []Outlet
	for _, outlet := range outletsDb {
		c := onemap.OneMapClient{}

		var latlong *LatLong
		if len(outlet.LatLong) != 2 {
			_latlong, lErr := c.GetLatLong(outlet.PostalCode)
			if lErr != nil {
				log.Printf("Error getting latlong: %v\n", lErr)
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

		outlets = append(outlets, Outlet{
			Name:          outlet.Name,
			Address:       outlet.Address,
			PostalCode:    outlet.PostalCode,
			OfficialLinks: outlet.OfficialLinks,
			LatLong:       latlong,
			ReviewLinks:   outlet.ReviewLinks,
			Id:            outlet.Id,
		})
	}
	return outlets, nil
}

func NewOutletService(repo *outlet.Repo) *Service {
	return &Service{repo: repo}
}
