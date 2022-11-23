package server

import (
	"github.com/IgorChicherin/currency-service/controllers"
	"github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.HTMLRender = gintemplate.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	upgrader := websocket.Upgrader{}

	health := new(controllers.HealthController)
	currency := controllers.CurrencyController{Upgrader: upgrader}

	router.GET("/", currency.Home)
	router.GET("/ping", health.Ping)
	router.GET("/alive", health.Alive)
	router.GET("/currency", currency.GetCurrency)
	router.GET("/ws", currency.GetCurrencyWs)
	return router

}
