package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Hello(c *gin.Context) {
	user, exists := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user, "exists": exists})
}
