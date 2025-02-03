package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sagotly/protoFlex.git/src/controllers"
	enteties "github.com/sagotly/protoFlex.git/src/entities"
)

type executableApi struct {
	ExecutableController *controllers.AddedExecutablesController
}

func NewExecutableApi(executableController *controllers.AddedExecutablesController) *executableApi {
	return &executableApi{
		ExecutableController: executableController,
	}
}

func (h *executableApi) GetExecutables(c *gin.Context) {
	executables, err := h.ExecutableController.GetAllExecutables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, executables)
}

func (h *executableApi) AddExecutable(c *gin.Context) {
	// Ожидаемый Request:
	// {
	//   "path": "/usr/bin/tool",
	//   "arguments": "--verbose --debug",
	//   "tunnel_id": "tun0"
	// }
	var req enteties.AddExecutableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	args := strings.Split(req.Arguments, "")

	err := h.ExecutableController.AddExecutableBtn(req.Path, args, req.TunnelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Executable added successfully"})
}

func (h *executableApi) GetAllTunnels(c *gin.Context) {
	tunnels, err := h.ExecutableController.GetAllTunnels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	fmt.Print("tunnels: ", tunnels)
	c.JSON(http.StatusOK, tunnels)
}

func (h *executableApi) ConnectExecutable(c *gin.Context) {
	// Ожидаемый Request:
	// {
	//   "tunnel_id": 1,
	//   "path": "/usr/bin/tool",
	//   "arguments": "--verbose --debug"
	// }
	var req enteties.ConnectExecutableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err := h.ExecutableController.ClickOnExecutableBtn(req.TunnelId, req.Path, req.Arguments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Executable connected successfully"})
}
