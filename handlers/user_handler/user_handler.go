package user_hanler

import (
	"net/http"
	"strconv"
	"strings"
	db "warcaby/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
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

func Login(c *gin.Context) {
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

func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID"})
		return
	}
	var user db.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Użytkownik nie znaleziony"})
		return
	}

	publicUser := gin.H{
		"ID":   user.ID,
		"Nick": user.Nick,
		"Bio":  user.Bio,
	}
	c.JSON(http.StatusOK, publicUser)
}

func GetMyUser(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak identyfikatora użytkownika"})
		return
	}
	userID := userIDVal.(int)
	var user db.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Użytkownik nie znaleziony"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nieautoryzowany dostęp"})
		return
	}
	currentUserID := userIDVal.(int)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID"})
		return
	}

	if currentUserID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Możesz modyfikować tylko swoje konto"})
		return
	}

	var updateData db.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Błąd dekodowania JSON"})
		return
	}

	var user db.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Użytkownik nie znaleziony"})
		return
	}

	if updateData.Nick != "" {
		user.Nick = updateData.Nick
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Bio != "" {
		user.Bio = updateData.Bio
	}
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy haszowaniu hasła"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy aktualizacji użytkownika"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nieautoryzowany dostęp"})
		return
	}
	currentUserID := userIDVal.(int)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID"})
		return
	}
	if currentUserID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Możesz usuwać tylko swoje konto"})
		return
	}
	if err := db.DB.Delete(&db.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy usuwaniu użytkownika"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Użytkownik usunięty"})
}

func SearchUsers(c *gin.Context) {
	query := c.Query("search")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parametr 'search' jest wymagany"})
		return
	}
	var users []db.User
	if err := db.DB.Where("lower(nick) LIKE ?", "%"+strings.ToLower(query)+"%").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas wyszukiwania"})
		return
	}
	publicUsers := make([]gin.H, 0, len(users))
	for _, user := range users {
		publicUser := gin.H{
			"ID":   user.ID,
			"Nick": user.Nick,
			"Bio":  user.Bio,
		}
		publicUsers = append(publicUsers, publicUser)
	}
	c.JSON(http.StatusOK, publicUsers)
}
