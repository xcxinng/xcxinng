package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/cars/booking", DO)
	router.POST("/cars/booking/undo", Undo)

	router.Run("localhost:19090")
}

func DO(c *gin.Context) {
	c.JSON(200, "success")
}

func Undo(c *gin.Context) {
	c.JSON(200, "success")
}
