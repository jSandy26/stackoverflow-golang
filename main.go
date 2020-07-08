package main

import (
	"net/http"

	"github.com/jsandy26/stackoverflow-golang/controllers"
	"github.com/jsandy26/stackoverflow-golang/middlewares"
	"github.com/jsandy26/stackoverflow-golang/models"

	"github.com/gin-gonic/gin"
)

var (
	r = gin.Default()
)

func main() {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	r.Use(middlewares.UserLoaderMiddleware())

	r.POST("/login", middlewares.Login)

	r.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		r.POST("/post", controllers.CreatePost)
		r.PUT("/posts/:id", controllers.UpdatePost)
		r.POST("/answer/:id", controllers.CreateAnswer)
	}

	r.GET("/posts/:id", controllers.ListPosts)

	models.ConnectDataBase()
	r.Run()
}
