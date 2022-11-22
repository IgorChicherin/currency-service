package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/IgorChicherin/currency-service/config"
	"github.com/IgorChicherin/currency-service/forms"
	"github.com/IgorChicherin/currency-service/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"strings"
)

type CurrencyController struct{}

func (receiver CurrencyController) GetCurrency(c *gin.Context) {
	var request forms.CurrencyRequest
	var errors []map[string]string
	conf := config.GetConfig()
	tsymsConf := conf.GetStringSlice("currencies.tsyms")
	fsymsConf := conf.GetStringSlice("currencies.fsyms")

	if c.ShouldBind(&request) != nil {
		c.IndentedJSON(http.StatusBadRequest, map[string]string{"error": "wrong parameters"})
		return
	}

	if !slices.Contains(fsymsConf, request.Fsyms) {
		errors = append(errors, map[string]string{request.Fsyms: "Symbol not configured"})
	}

	tsyms := strings.Split(request.Tsyms, ",")

	for _, tsym := range tsyms {
		if !slices.Contains(tsymsConf, tsym) {
			errors = append(errors, map[string]string{tsym: "Symbol not configured"})
		}
	}

	if len(errors) == 0 {
		data, err := receiver.getCurrencyFromHttp(request.Fsyms, &tsyms)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "")
			return
		}
		c.IndentedJSON(http.StatusOK, data)
		return
	}

	c.IndentedJSON(http.StatusBadRequest, errors)
}

func (receiver CurrencyController) getCurrencyFromDb(fsym string, tsyms *[]string) (map[string]interface{}, error) {
	//TODO: implement this
	return map[string]interface{}{}, nil
}

func (receiver CurrencyController) getCurrencyFromHttp(fsym string, tsyms *[]string) (models.CurrencyResponse, error) {
	var resp models.CurrencyResponse
	URL := fmt.Sprintf(
		"https://min-api.cryptocompare.com/data/pricemultifull?fsyms=%s&tsyms=%s",
		fsym, strings.Join(*tsyms, ","))
	response, err := http.Get(URL)
	if err != nil {
		return resp, err
	}

	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	return resp, err
}
