package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"ritmoexample/cmd/models"
	"ritmoexample/cmd/services/companies"
	"ritmoexample/cmd/services/offers"
)

func main() {
	if err := models.ConnectRepo(); err != nil {
		log.Fatal(err)
	}

	companiesServer := companies.NewServer()
	offersServer := offers.NewServer()
	router := gin.Default()

	// TODO(david): Here we would have to authenticate the requests with a middleware

	// Get company info
	router.GET("/companies/:id", companiesServer.Get)
	// Create a new company
	router.POST("/companies", companiesServer.Create)
	// Update company info
	router.PUT("/companies/:id", companiesServer.Update)

	// Get offers
	router.GET("/offers/:id", offersServer.Get)
	// Create a new offer
	router.POST("/offers", offersServer.Create)
	// Update offer status
	router.PUT("/offers/:id/status", offersServer.UpdateStatus)
	// Update if user accept the offer
	router.PUT("/offers/:id/accepted", offersServer.UpdateAccepted)

	router.Run(":8080")
}
