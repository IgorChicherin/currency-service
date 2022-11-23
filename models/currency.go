package models

import (
	"encoding/json"
	"github.com/IgorChicherin/currency-service/db"
	"strings"
)

type CurrencyResponse struct {
	Raw     map[string]map[string]CurrencyInfoRaw     `json:"RAW"`
	Display map[string]map[string]CurrencyInfoDisplay `json:"DISPLAY"`
}

func (receiver *CurrencyResponse) Save() error {
	for fsym, tsyms := range receiver.Raw {
		for tsym, currencyInfoRaw := range tsyms {
			var currencyPair CurrencyPairInfo
			currencyPair.Symbol = strings.Join([]string{fsym, tsym}, ":")
			currencyPair.Raw = currencyInfoRaw
			currencyPair.Display = receiver.Display[fsym][tsym]

			err := currencyPair.Save()

			if err != nil {
				return err
			}
		}
	}
	return nil
}

type CurrencyRequestWs struct {
	Fsyms string   `json:"fsyms"`
	Tsyms []string `json:"tsyms"`
}

type CurrencyInfoRaw struct {
	CHANGE24HOUR    float64 `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR float64 `json:"CHANGEPCT24HOUR"`
	OPEN24HOUR      float64 `json:"OPEN24HOUR"`
	VOLUME24HOUR    float64 `json:"VOLUME24HOUR"`
	VOLUME24HOURTO  float64 `json:"VOLUME24HOURTO"`
	LOW24HOUR       float64 `json:"LOW24HOUR"`
	HIGH24HOUR      float64 `json:"HIGH24HOUR"`
	PRICE           float64 `json:"PRICE"`
	SUPPLY          float64 `json:"SUPPLY"`
	MKTCAP          float64 `json:"MKTCAP"`
}

type CurrencyInfoDisplay struct {
	CHANGE24HOUR    string `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR string `json:"CHANGEPCT24HOUR"`
	OPEN24HOUR      string `json:"OPEN24HOUR"`
	VOLUME24HOUR    string `json:"VOLUME24HOUR"`
	VOLUME24HOURTO  string `json:"VOLUME24HOURTO"`
	LOW24HOUR       string `json:"LOW24HOUR"`
	HIGH24HOUR      string `json:"HIGH24HOUR"`
	PRICE           string `json:"PRICE"`
	SUPPLY          string `json:"SUPPLY"`
	MKTCAP          string `json:"MKTCAP"`
}

type CurrencyPairInfo struct {
	Symbol  string              `json:"symbol"`
	Raw     CurrencyInfoRaw     `json:"raw"`
	Display CurrencyInfoDisplay `json:"display"`
}

func (receiver *CurrencyPairInfo) Save() error {
	database := db.GetDB()
	data, err := json.Marshal(&receiver)

	query := "INSERT INTO ticks (symbol, data) VALUES ($1, $2);"
	_, err = database.Exec(query, &receiver.Symbol, &data)

	return err
}

func (receiver *CurrencyPairInfo) Load(symbol string) error {
	database := db.GetDB()
	var data []byte
	row := database.QueryRow(
		"SELECT data FROM ticks WHERE symbol=$1 ORDER BY created_at DESC LIMIT 1",
		symbol)
	err := row.Scan(&data)

	err = json.Unmarshal(data, &receiver)
	return err
}
