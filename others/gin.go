package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type VO struct {
	Data string `form:"data" json:"data"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		var v VO
		c.MustBindWith(&v, binding.Query)
		c.JSON(http.StatusOK, v)
	})
	r.POST("/ping", func(c *gin.Context) {
		var v VO
		c.ShouldBindJSON(&v)
		fmt.Println(v.Data)
		c.JSON(http.StatusOK, v)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
