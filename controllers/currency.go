package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/IgorChicherin/currency-service/config"
	"github.com/IgorChicherin/currency-service/forms"
	"github.com/IgorChicherin/currency-service/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"strings"
)

type CurrencyController struct {
	Upgrader websocket.Upgrader
}

func (cc *CurrencyController) Home(c *gin.Context) {
	conf := config.GetConfig()
	ws := fmt.Sprintf("ws://%s:%s/ws", conf.GetString("server.host"), conf.GetString("server.port"))
	c.HTML(http.StatusOK, "home.html", gin.H{"host": ws})
}

func (cc *CurrencyController) GetCurrencyWs(c *gin.Context) {
	cc.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := cc.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var currency models.CurrencyRequestWs
		err = json.Unmarshal(message, &currency)

		if err != nil {
			err = conn.WriteJSON(map[string]string{"error": "Decoding error"})
			break
		}

		errors, currencyResponse, respErr := cc.getCurrency(currency.Fsyms, &currency.Tsyms)

		if len(errors) == 0 && respErr == nil {
			err = conn.WriteJSON(&currencyResponse)
		}

		if len(errors) != 0 {
			err = conn.WriteJSON(map[string]any{"errors": errors})
		}

		if respErr != nil {
			err = conn.WriteJSON(map[string]string{"error": "Internal error"})
		}

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (cc *CurrencyController) GetCurrency(c *gin.Context) {
	var request forms.CurrencyRequest

	if c.ShouldBind(&request) != nil {
		c.IndentedJSON(http.StatusBadRequest, map[string]string{"error": "wrong parameters"})
		return
	}

	tsyms := strings.Split(request.Tsyms, ",")

	errors, currencyResponse, err := cc.getCurrency(request.Fsyms, &tsyms)

	if len(errors) != 0 {
		c.IndentedJSON(http.StatusBadRequest, errors)
		return
	}

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}

	c.IndentedJSON(http.StatusOK, currencyResponse)
}

func (cc *CurrencyController) getCurrency(fsyms string, tsyms *[]string) ([]map[string]string, models.CurrencyResponse, error) {
	var errors []map[string]string
	conf := config.GetConfig()
	tsymsConf := conf.GetStringSlice("currencies.tsyms")
	fsymsConf := conf.GetStringSlice("currencies.fsyms")

	if !slices.Contains(fsymsConf, fsyms) {
		errors = append(errors, map[string]string{fsyms: "Symbol not configured"})
	}

	for _, tsym := range *tsyms {
		if !slices.Contains(tsymsConf, tsym) {
			errors = append(errors, map[string]string{tsym: "Symbol not configured"})
		}
	}

	if len(errors) != 0 {
		return errors, models.CurrencyResponse{}, nil
	}

	data, err := cc.getCurrencyFromHttp(fsyms, tsyms)
	if err != nil {
		dbData, err := cc.getCurrencyFromDb(fsyms, tsyms)
		if err != nil {
			return errors, models.CurrencyResponse{}, err
		}
		return errors, dbData, nil
	}
	return errors, data, nil
}

func (cc *CurrencyController) getCurrencyFromDb(fsym string, tsyms *[]string) (models.CurrencyResponse, error) {
	var data models.CurrencyResponse
	raw := make(map[string]map[string]models.CurrencyInfoRaw)
	display := make(map[string]map[string]models.CurrencyInfoDisplay)
	for _, tsym := range *tsyms {
		var currencyPair = models.CurrencyPairInfo{}
		symbol := strings.Join([]string{fsym, tsym}, ":")
		err := currencyPair.Load(symbol)
		if err != nil {
			return data, err
		}

		if rawData, found := raw[fsym]; !found {
			raw[fsym] = map[string]models.CurrencyInfoRaw{tsym: currencyPair.Raw}
		} else {
			rawData[tsym] = currencyPair.Raw
		}

		if displayData, found := display[fsym]; !found {
			display[fsym] = map[string]models.CurrencyInfoDisplay{tsym: currencyPair.Display}
		} else {
			displayData[tsym] = currencyPair.Display
		}
	}
	data.Raw = raw
	data.Display = display
	return data, nil
}

func (cc *CurrencyController) getCurrencyFromHttp(fsym string, tsyms *[]string) (models.CurrencyResponse, error) {
	return GetCurrencyFromHttp(fsym, tsyms)
}

func GetCurrencyFromHttp(fsym string, tsyms *[]string) (models.CurrencyResponse, error) {
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
