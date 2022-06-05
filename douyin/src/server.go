package main

import (
	"douyin/src/api"
	"douyin/src/global"
	"douyin/src/pkg/setting"
	"douyin/src/pojo/entity"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

//项目启动会自动调用此初始化函数
func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
}

//初始化配置文件读取
func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	global.JWTSecret = []byte(global.JWTSecret)
	return nil
}

//初始化数据库连接
func setupDBEngine() error {
	var err error
	global.DBEngine, err = entity.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}

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
