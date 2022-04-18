package main

import (
	"IMConnection/conf"
	"IMConnection/model"
	"IMConnection/pkg/logging"
	"IMConnection/router"
	"fmt"
	"net/http"
)

func init() {
	conf.Setup()
	model.Setup()
	logging.Setup()
}

func main() {
	router := router.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    conf.ServerSetting.ReadTimeout,
		WriteTimeout:   conf.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
