package mk_outlet_service

import "golang-server/cmd/product/makanplace/repositories/outlet"

type ServiceBody struct {
	OutletName    string
	OutletType    string
	ProductName   string
	Address       string
	PostalCode    string
	OfficialLinks []string
}

type OutletService struct {
	repo *outlet.Repo
}

func (s *OutletService) AddOutlet(b ServiceBody) error {
	return s.repo.NewOutletWithMenu(b.OutletName, b.Address, b.PostalCode, b.OfficialLinks, []string{b.ProductName})
}

func NewOutletService(repo *outlet.Repo) *OutletService {
	return &OutletService{repo: repo}
}
