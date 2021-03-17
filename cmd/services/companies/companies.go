package companies

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"

	"ritmoexample/cmd/errors"
	"ritmoexample/cmd/models"
)

type Server struct{}

func NewServer() *Server {
	return new(Server)
}

func serializeCompany(company *models.Company) *Company {
	return &Company{
		ID:             company.ID,
		Name:           company.Name,
		FiscalID:       company.FiscalID,
		Address:        company.Address,
		AddressLineTwo: company.AddressLineTwo,
		Region:         company.Region,
		City:           company.City,
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

	company, ok := models.Repo.Companies()[req.ID]
	if !ok {
		ctx.JSON(http.StatusNotFound, errors.NewError(fmt.Sprintf("company not found: %s", req.ID)))
		return
	}

	ctx.JSON(http.StatusOK, serializeCompany(company))
}

type Company struct {
	ID             string `uri:"id,omitempty"`
	Name           string `json:"name,required"`
	FiscalID       string `json:"fiscalId,required"`
	Address        string `json:"address,required"`
	AddressLineTwo string `json:"addressLineTwo,required"`
	Region         string `json:"region,required"`
	City           string `json:"city,required"`
}

func (s *Server) Create(ctx *gin.Context) {
	var req Company
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}

	company := &models.Company{
		ID:             ksuid.New().String(),
		Name:           req.Name,
		FiscalID:       req.FiscalID,
		Address:        req.Address,
		AddressLineTwo: req.AddressLineTwo,
		Region:         req.Region,
		City:           req.City,
	}
	if err := models.Repo.SetCompany(company); err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, serializeCompany(company))
}

func (s *Server) Update(ctx *gin.Context) {
	var uriReq Company
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}
	company, ok := models.Repo.Companies()[uriReq.ID]
	if !ok {
		ctx.JSON(http.StatusNotFound, errors.NewError(fmt.Sprintf("company not found: %s", uriReq.ID)))
		return
	}

	var jsonReq Company
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError(err.Error()))
		return
	}

	company.Name = jsonReq.Name
	company.FiscalID = jsonReq.FiscalID
	company.Address = jsonReq.Address
	company.AddressLineTwo = jsonReq.AddressLineTwo
	company.Region = jsonReq.Region
	company.City = jsonReq.City
	if err := models.Repo.SetCompany(company); err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, serializeCompany(company))
}
