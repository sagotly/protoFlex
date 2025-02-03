package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
	"github.com/sagotly/protoFlex.git/src/api"
	"github.com/sagotly/protoFlex.git/src/client"
	"github.com/sagotly/protoFlex.git/src/controllers"
	"github.com/sagotly/protoFlex.git/src/repo"
	"github.com/sagotly/protoFlex.git/src/utils"
)

// Executable - структура для хранения информации о "исполняемом файле"
type Executable struct {
	ID        int      `json:"id"`        // Уникальный идентификатор
	Path      string   `json:"path"`      // Путь к исполняемому файлу
	Arguments []string `json:"arguments"` // Аргументы запуска
	TunnelId  string   `json:"tunnel_id"` // ID (или имя) туннеля
	Active    bool     `json:"active"`    // Флаг, подключено ли исполнение
}

// Tunnel - структура для туннелей
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

	// Указываем папку с HTML-шаблонами
	r.LoadHTMLFiles("src/templates/index.html", "src/templates/server.html")

	// Обработчик для главной страницы
	r.GET("/", func(c *gin.Context) {
		// Рендерим `index.html`
		c.HTML(200, "index.html", gin.H{
			"title": "Welcome to Main Page",
		})
	})
	// Обработчик для главной страницы
	r.GET("/s", func(c *gin.Context) {
		// Рендерим `index.html`
		c.HTML(200, "server.html", gin.H{
			"title": "Welcome to Main Page",
		})
	})
	// 1. Получить список всех Executables
	r.GET("/executables", executableApi.GetExecutables)

	// 2. Добавить Executable
	r.POST("/executables", executableApi.AddExecutable)

	// 3. Подключение (Connect) к Executable
	r.POST("/executables/connect", executableApi.ConnectExecutable)

	// 4. Получить список туннелей
	r.GET("/tunnels", executableApi.GetAllTunnels)

	// Получить список серверов
	r.GET("/servers", serverApi.GetServers)

	// Генерация токена
	r.POST("/connections/generate-token", tokenApi.GenerateToken)

	// Проверка токена
	r.POST("/connections/validate-token", tokenApi.ValidateToken)

	// Добавить сервер
	r.POST("/servers", serverApi.AddServer)

	r.Run(":8080")
}
