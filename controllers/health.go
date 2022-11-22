package controllers

import (
	"github.com/IgorChicherin/currency-service/config"
	"github.com/IgorChicherin/currency-service/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HealthController struct{}

func (h HealthController) basePingResponse() map[string]interface{} {
	conf := config.GetConfig()
	return map[string]interface{}{
		"component": conf.GetString("service.component"),
		"version":   conf.GetString("service.version"),
		"host":      conf.GetString("server.host"),
		"time":      time.Now().Unix(),
	}
}

func (h HealthController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, h.basePingResponse())
}

func (h HealthController) Alive(c *gin.Context) {
	response := h.basePingResponse()
	if err := db.GetDB().Ping(); err != nil {
		response["alive"] = false
	}

	response["alive"] = true

	c.JSON(http.StatusOK, response)
}
