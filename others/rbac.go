package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/storyicon/grbac"
)

func QueryRolesByHeaders(header http.Header) (roles []string, err error) {
	// role := header.Get("role")
	// if role == "" {
	// 	return nil, errors.New("role not found")
	// }
	roles = append(roles, "editor")
	return roles, err
}

func Authentication() gin.HandlerFunc {
	rbac, err := grbac.New(grbac.WithJSON("rules.json", time.Minute*10))
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		roles, err := QueryRolesByHeaders(c.Request.Header)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		state, err := rbac.IsRequestGranted(c.Request, roles)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !state.IsGranted() {
			c.AbortWithStatusJSON(200, "Forbidden")
			return
		}
	}
}

func xx() {
	c := gin.New()
	c.Use(Authentication())

	c.POST("/articles", func(ctx *gin.Context) {
		ctx.JSON(200, "success")
	})
	c.PUT("/articles", func(ctx *gin.Context) {
		ctx.JSON(200, "success")
	})
	c.GET("/articles", func(ctx *gin.Context) {
		ctx.JSON(200, "success")
	})

	c.Run(":8080")
}
