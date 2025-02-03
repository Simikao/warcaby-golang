package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak identyfikatora użytkownika (X-User-ID)"})
			c.Abort()
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Nieprawidłowy identyfikator użytkownika"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
