package main

import (
	"github.com/gin-gonic/gin"

	"douyin/src/api"
)

func main() {
	r := gin.Default()

	api.InitRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
