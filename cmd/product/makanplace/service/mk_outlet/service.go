package mk_outlet_service

import (
	"log"

	"golang-server/cmd/product/makanplace/client/onemap"
	"golang-server/cmd/product/makanplace/repositories/outlet"
)

type ServiceBody struct {
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

func (s *Service) AddOutlet(b ServiceBody) error {
	return s.repo.NewOutletWithMenu(b.OutletName, b.Address, b.PostalCode, b.OfficialLinks, b.ReviewLinks, []string{b.ProductName})
}

type LatLong = onemap.LatLong
type Outlet struct {
	Name          string   `json:"name"`
	Address       string   `json:"address"`
	PostalCode    string   `json:"postalCode"`
	OfficialLinks []string `json:"officialLinks"`
	ReviewLinks   []string `json:"reviewLinks"`

	*LatLong `json:"latlong"`
	Id       int64 `json:"id"`
}

func (s *Service) GetOutlets() ([]Outlet, error) {
	outletsDb, err := s.repo.GetOutlets()
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
