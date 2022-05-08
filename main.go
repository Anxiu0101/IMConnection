package main

import (
	"IMConnection/cache"
	"IMConnection/conf"
	"IMConnection/model"
	"IMConnection/pkg/logging"
	"IMConnection/router"
	"IMConnection/service"
	"fmt"
	"net/http"
)

func init() {
	conf.Setup()
	model.Setup()
	cache.Setup()
	logging.Setup()
}

func main() {
	router := router.InitRouter()

	// 开启服务管理监听
	go service.Manager.Listen()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    conf.ServerSetting.ReadTimeout,
		WriteTimeout:   conf.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
