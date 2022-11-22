package server

import (
	"github.com/IgorChicherin/currency-service/controllers"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)
	currency := new(controllers.CurrencyController)

	router.GET("/ping", health.Ping)
	router.GET("/alive", health.Alive)
	router.GET("/currency", currency.GetCurrency)

	return router

}
