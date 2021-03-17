package models

import (
	"encoding/json"
	"io/ioutil"

	"github.com/juju/errors"
)

var Repo *LocalStorage

type LocalStorage struct {
	StoredCompanies []*Company `json:"companies"`
	StoredOffers    []*Offer   `json:"offers"`
}

func ConnectRepo() error {
	Repo = new(LocalStorage)

	file, err := ioutil.ReadFile("storage/storage.json")
	if err != nil {
		return errors.Trace(err)
	}
	if err := json.Unmarshal(file, &Repo); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (repo *LocalStorage) Companies() map[string]*Company {
	companies := map[string]*Company{}
	for _, company := range Repo.StoredCompanies {
		companies[company.ID] = company
	}

	return companies
}

func (repo *LocalStorage) SetCompany(company *Company) error {
	if !checkCompany(Repo.StoredCompanies, company) {
		Repo.StoredCompanies = append(Repo.StoredCompanies, company)
	}

	bytes, err := json.Marshal(Repo)
	if err != nil {
		return errors.Trace(err)
	}
	if err := ioutil.WriteFile("storage/storage.json", bytes, 0644); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func checkCompany(companies []*Company, company *Company) bool {
	for _, c := range companies {
		if c.ID == company.ID {
			c = company
			return true
		}
	}

	return false
}

func (repo *LocalStorage) Offers() map[string]*Offer {
	offers := map[string]*Offer{}
	for _, offer := range Repo.StoredOffers {
		offers[offer.ID] = offer
	}

	return offers
}

func (repo *LocalStorage) SetOffer(offer *Offer) error {
	if !checkOffer(Repo.StoredOffers, offer) {
		Repo.StoredOffers = append(Repo.StoredOffers, offer)
	}

	bytes, err := json.Marshal(Repo)
	if err != nil {
		return errors.Trace(err)
	}
	if err := ioutil.WriteFile("storage/storage.json", bytes, 0644); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func checkOffer(offers []*Offer, offer *Offer) bool {
	for _, o := range offers {
		if o.ID == offer.ID {
			o = offer
			return true
		}
	}

	return false
}
