package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sagotly/protoFlex.git/src/api"
	"github.com/sagotly/protoFlex.git/src/client"
	"github.com/sagotly/protoFlex.git/src/controllers"
	"github.com/sagotly/protoFlex.git/src/repo"
	"github.com/sagotly/protoFlex.git/src/utils"
)

// Executable - structure to store information about the executable file
type Executable struct {
	ID        int      `json:"id"`        // Unique identifier
	Path      string   `json:"path"`      // Path to the executable file
	Arguments []string `json:"arguments"` // Launch arguments
	TunnelId  string   `json:"tunnel_id"` // ID (or name) of the tunnel
	Active    bool     `json:"active"`    // Flag indicating whether the execution is connected
}

// Tunnel - structure for tunnels
type Tunnel struct {
	ID            int    `json:"id"`
	InterfaceName string `json:"interface_name"`
}

func main() {
	fmt.Println("What a long night ahead... ")
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		log.Fatal(err) // Log the error if the database connection fails
	}
	defer db.Close()
	if err := utils.SetupDatabase(db); err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	serverRepo := repo.NewServerRepo(db)
	tunnelRepo := repo.NewTunnelRepo(db)
	addedExecutablesRepo := repo.NewAddedExecutablesRepo(db)

	tokenClient := client.NewServerClient()

	tokenController := controllers.NewTokenController(tokenClient)
	serverViewController := controllers.NewServerViewController(tunnelRepo, serverRepo)
	addedExecutablesController := controllers.NewAddedExcecutablesController(tunnelRepo, serverRepo, addedExecutablesRepo)

	tokenApi := api.NewTokenApi(tokenController)
	executableApi := api.NewExecutableApi(addedExecutablesController)
	serverApi := api.NewServerApi(serverViewController)

	r := gin.Default()

	// Specify the folder with HTML templates
	r.LoadHTMLFiles("src/templates/index.html", "src/templates/server.html")

	// Handler for the main page
	r.GET("/", func(c *gin.Context) {
		// Render `index.html`
		c.HTML(200, "index.html", gin.H{
			"title": "Welcome to Main Page",
		})
	})
	// Handler for the server page
	r.GET("/s", func(c *gin.Context) {
		// Render `server.html`
		c.HTML(200, "server.html", gin.H{
			"title": "Welcome to Server Page",
		})
	})
	// 1. Get the list of all Executables
	r.GET("/executables", executableApi.GetExecutables)

	// 2. Add an Executable
	r.POST("/executables", executableApi.AddExecutable)

	// 3. Connect to an Executable
	r.POST("/executables/connect", executableApi.ConnectExecutable)

	// 4. Get the list of tunnels
	r.GET("/tunnels", executableApi.GetAllTunnels)

	// Get the list of servers
	r.GET("/servers", serverApi.GetServers)

	// Generate a token
	r.POST("/connections/generate-token", tokenApi.GenerateToken)

	// Validate a token
	r.POST("/connections/validate-token", tokenApi.ValidateToken)
	// Add a server
	r.POST("/servers", serverApi.AddServer)

	r.Run(":8080")
}
