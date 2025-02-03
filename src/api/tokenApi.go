package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagotly/protoFlex.git/src/controllers"
	enteties "github.com/sagotly/protoFlex.git/src/entities"
)

type TokenApi struct {
	TokenController *controllers.TokenController
}

func NewTokenApi(tokenController *controllers.TokenController) *TokenApi {
	return &TokenApi{
		TokenController: tokenController,
	}
}

func (t *TokenApi) GenerateToken(c *gin.Context) {
	var req enteties.GenerateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	err := t.TokenController.GenerateToken(req.Ip, req.Port)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	c.JSON(200, gin.H{"message": "Token generated successfully"})
}

func (t *TokenApi) ValidateToken(c *gin.Context) {
	var req enteties.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	fmt.Println("req.Ip", req.Ip)
	err := t.TokenController.ValidateToken(req.Ip, req.Port, req.Token)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	c.JSON(200, gin.H{"message": "Token validated successfully"})
}
