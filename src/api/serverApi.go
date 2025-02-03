package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagotly/protoFlex.git/src/controllers"
	enteties "github.com/sagotly/protoFlex.git/src/entities"
)

type ServerApi struct {
	tokenController  *controllers.TokenController
	ServerController *controllers.ServerViewController
}

func NewServerApi(serverController *controllers.ServerViewController) *ServerApi {
	return &ServerApi{
		ServerController: serverController,
	}
}

func (h *ServerApi) GetServers(c *gin.Context) {
	servers, err := h.ServerController.GetAllServers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, servers)
}

func (h *ServerApi) AddServer(c *gin.Context) {
	var req enteties.AddServerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.ServerController.CreateNewServerBtn(req.Name, req.Ip, req.TunnelList)
	if err != nil && err.Error() == "tunnel already exists" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tunnel already exists"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Server added successfully"})
}
