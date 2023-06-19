package controllers

import (
	"minecraft_searcher/api/rest/token"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	type loginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var input loginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.Compare(input.Username, "wojtess") == 0 && bcrypt.CompareHashAndPassword([]byte("$2a$12$.w9OX99rT5b1PYrU4eiJfuEh9gSkjgq/ydelHPCWvWX3q3ltWoo7q"), []byte(input.Password)) == nil {
		tokens, err := token.GenerateTokens(0)
		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, tokens)
	} else {
		c.JSON(http.StatusAccepted, gin.H{"error": "wrong username or password"})
	}

}

func Refresh(c *gin.Context) {
	type refreshToken struct {
		RefreshToken string `json:"refresh_token"`
	}
	var input refreshToken

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := token.GetToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		userId := int(claims["user_id"].(float64))
		tokens, err := token.GenerateTokens(userId)
		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, tokens)
	}
}
