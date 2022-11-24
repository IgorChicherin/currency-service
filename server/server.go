package server

import (
	"fmt"
	"github.com/IgorChicherin/currency-service/config"
	"github.com/IgorChicherin/currency-service/controllers"
	"github.com/go-co-op/gocron"
	"log"
	"net/http"
	"time"
)

func Run() (*http.Server, *gocron.Scheduler, error) {
	shed, err := RunCronJob(fetchSymbols)

	if err != nil {
		return nil, nil, err
	}

	return NewHttpServer(), shed, nil
}

func fetchSymbols() {
	conf := config.GetConfig()
	tsymsConf := conf.GetStringSlice("currencies.tsyms")
	fsymsConf := conf.GetStringSlice("currencies.fsyms")

	for _, fsym := range fsymsConf {
		response, err := controllers.GetCurrencyFromHttp(fsym, &tsymsConf)

		if err != nil {
			log.Println("fetchSymbols fetch task error")
		}

		err = response.Save()

		if err != nil {
			log.Println("fetchSymbols save task error")
		}
	}
}

func NewHttpServer() *http.Server {
	conf := config.GetConfig()
	router := NewRouter()
	addr := fmt.Sprintf("%s:%s", conf.GetString("server.host"), conf.GetString("server.port"))

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func RunCronJob(job func()) (*gocron.Scheduler, error) {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Minutes().Do(func() {
		job()
	})
	if err != nil {
		return s, err
	}

	s.StartAsync()

	return s, nil
}
