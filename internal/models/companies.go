package models

type Company struct {
	ID string

	Name     string
	FiscalID string

	Address        string
	AddressLineTwo string
	Region         string
	City           string
}
