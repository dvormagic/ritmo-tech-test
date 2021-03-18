package models

import "libs.altipla.consulting/money"

type StatusType string

const (
	StatusTypeRejected   StatusType = "STATUS_REJECTED"
	StatusTypePending    StatusType = "STATUS_PENDING"
	StatusTypePreAproved StatusType = "STATUS_PREAPROVED"
	StatusTypeAproved    StatusType = "STATUS_APROVED"
)

type Offer struct {
	ID        string
	CompanyID string

	Status          StatusType
	Charges         int64
	SalesPercentage int64
	Accepted        bool

	RawAdvance string
	RawRefund  string
}

func (offer *Offer) Advance() money.Money {
	if offer.RawAdvance == "" {
		return money.FromCents(0)
	}
	advance, _ := money.Parse(offer.RawAdvance)
	return advance
}

func (offer *Offer) SetAdvance(advance money.Money) {
	offer.RawAdvance = advance.Format(money.FormatConfig{ForceDecimals: true})
}

func (offer *Offer) Refund() money.Money {
	if offer.RawRefund == "" {
		return money.FromCents(0)
	}
	price, _ := money.Parse(offer.RawRefund)
	return price
}

func (offer *Offer) SetRefund(price money.Money) {
	offer.RawRefund = price.Format(money.FormatConfig{ForceDecimals: true})
}
