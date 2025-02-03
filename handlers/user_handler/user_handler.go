package main

import (
	"net/http"
	db "warcaby/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func registerUser(c *gin.Context) {
	var newUser db.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Błąd dekodowania JSON"})
		return
	}

	if newUser.Nick == "" || newUser.Email == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nick, email i hasło są wymagane"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy haszowaniu hasła"})
		return
	}
	newUser.Password = string(hashedPassword)

	if err := db.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas tworzenia użytkownika"})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

func login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Błąd dekodowania JSON"})
		return
	}

	var user db.User
	if err := db.DB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Niepoprawny email lub hasło"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Niepoprawny email lub hasło"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logowanie udane", "user": user})
}
