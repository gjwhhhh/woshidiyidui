package main

import (
	"douyin/src/api"
	"douyin/src/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := gin.Default()
	api.InitRouter(router)

	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
