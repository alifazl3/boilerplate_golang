package router

import (
	"boilerplate/api/handler"
	"boilerplate/internal/config"
	"boilerplate/internal/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	DB *gorm.DB
}

func (r *Router) Routes(conf *config.Config) error {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(func(c *gin.Context) {
		c.Set("DB", r.DB)
		c.Set("Config", conf)
		c.Next()
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	// Auth routes
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(conf, r.DB))
	protected.POST("/logout", handler.Logout)
	protected.POST("/changePassword", handler.ChangePassword)
	protected.GET("hello", handler.Hello)

	//router.Static("/cdn", "./tmp/")

	return router.Run("0.0.0.0:" + conf.AppPort)
}
