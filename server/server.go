package server

import (
	"fmt"
	"github.com/IgorChicherin/currency-service/config"
)

func Init() {
	conf := config.GetConfig()
	r := NewRouter()
	addr := fmt.Sprintf("%s:%s", conf.GetString("server.host"), conf.GetString("server.port"))
	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}
