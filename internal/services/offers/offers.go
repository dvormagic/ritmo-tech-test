package offers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"libs.altipla.consulting/money"

	"ritmoexample/internal/errors"
	"ritmoexample/internal/models"
)

type Server struct{}

func NewServer() *Server {
	return new(Server)
}

type Offer struct {
	ID              string `json:"id,omitempty"`
	CompanyID       string `json:"companyId"`
	Status          string `json:"status"`
	Charges         int64  `json:"charges"`
	SalesPercentage int64  `json:"salesPercentage"`
	Accepted        bool   `json:"accepted"`
	Advance         string `json:"advance"`
	Refund          string `json:"refund"`
}

func serializeOffer(offer *models.Offer) *Offer {
	return &Offer{
		ID:              offer.ID,
		CompanyID:       offer.CompanyID,
		Status:          string(offer.Status),
		Charges:         offer.Charges,
		SalesPercentage: offer.SalesPercentage,
		Advance:         offer.Advance().Format(money.EUR),
		Refund:          offer.Refund().Format(money.EUR),
		Accepted:        offer.Accepted,
	}
}

type getRequest struct {
	ID string `uri:"id"`
}

func (s *Server) Get(ctx *gin.Context) {
	var req getRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}

	offer, ok := models.Repo.Offers()[req.ID]
	if !ok {
		ctx.JSON(http.StatusNotFound, errors.NewError(fmt.Sprintf("offer not found: %s", req.ID)))
		return
	}

	ctx.JSON(http.StatusOK, serializeOffer(offer))
}

type getByCompanyRequest struct {
	CompanyID string `uri:"company-id"`
}

func (s *Server) GetByCompany(ctx *gin.Context) {
	var req getByCompanyRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}

	offer, ok := models.Repo.OffersByCompany()[req.CompanyID]
	if !ok {
		ctx.JSON(http.StatusNotFound, errors.NewError(fmt.Sprintf("not offer with company found: %s", req.CompanyID)))
		return
	}

	ctx.JSON(http.StatusOK, serializeOffer(offer))
}

type createRequest struct {
	CompanyID       string `json:"companyId,required"`
	Charges         int64  `json:"charges,required"`
	SalesPercentage int64  `json:"salesPercentage,required"`
	Advance         string `json:"advance,required"`
	Refund          string `json:"refund,required"`
}

func (s *Server) Create(ctx *gin.Context) {
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}

	if _, ok := models.Repo.Companies()[req.CompanyID]; !ok {
		ctx.JSON(http.StatusNotFound, errors.NewError(fmt.Sprintf("company not found: %s", req.CompanyID)))
		return
	}

	offer := &models.Offer{
		ID:              ksuid.New().String(),
		CompanyID:       req.CompanyID,
		Status:          models.StatusTypePending,
		Charges:         req.Charges,
		SalesPercentage: req.SalesPercentage,
	}
	advance, err := money.Parse(req.Advance)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewError(fmt.Sprintf("bad advance format: %s", req.Advance)))
		return
	}
	offer.SetAdvance(advance)

	refund, err := money.Parse(req.Refund)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewError(fmt.Sprintf("bad refund format: %s", req.Refund)))
		return
	}
	offer.SetRefund(refund)

	if err := models.Repo.SetOffer(offer); err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, serializeOffer(offer))
}

type updateStatusRequest struct {
	ID     string `uri:"id,required"`
	Status string `json:"status,required"`
}

func (s *Server) UpdateStatus(ctx *gin.Context) {
	var uriReq updateStatusRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}

	offer, ok := models.Repo.Offers()[uriReq.ID]
	if !ok {
		ctx.JSON(http.StatusNotFound, errors.NewError(fmt.Sprintf("offer not found: %s", uriReq.ID)))
		return
	}

	var jsonReq updateStatusRequest
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}

	switch jsonReq.Status {
	case string(models.StatusTypeRejected):
		offer.Status = models.StatusTypeRejected
	case string(models.StatusTypePreAproved):
		offer.Status = models.StatusTypePreAproved
	case string(models.StatusTypeAproved):
		offer.Status = models.StatusTypeAproved
	default:
		ctx.JSON(http.StatusBadRequest, errors.NewError(fmt.Sprintf("invalid status type: %s", jsonReq.Status)))
		return
	}
	if err := models.Repo.SetOffer(offer); err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, serializeOffer(offer))
}

type updateAcceptedRequest struct {
	ID       string `uri:"id,required"`
	Accepted bool   `json:"accepted,required"`
}

func (s *Server) UpdateAccepted(ctx *gin.Context) {
	var uriReq updateAcceptedRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}
	offer, ok := models.Repo.Offers()[uriReq.ID]
	if !ok {
		ctx.JSON(http.StatusNotFound, errors.NewError(fmt.Sprintf("offer not found: %s", uriReq.ID)))
		return
	}

	var jsonReq updateAcceptedRequest
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}
	offer.Accepted = jsonReq.Accepted
	if err := models.Repo.SetOffer(offer); err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, serializeOffer(offer))
}
