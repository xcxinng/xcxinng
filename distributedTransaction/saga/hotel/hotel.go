package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/hotels/booking", DO)
	router.POST("/hotels/booking/undo", Undo)

	router.Run("localhost:19091")
}

func DO(c *gin.Context) {
	c.JSON(200, "success")
	// c.AbortWithError(500, errors.New("internal server error"))
}

func Undo(c *gin.Context) {
	c.JSON(200, "success")
}
