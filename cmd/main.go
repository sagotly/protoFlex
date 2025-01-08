package main

import (
	"database/sql"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "github.com/mattn/go-sqlite3"
	enteties "github.com/sagotly/protoFlex.git/src/entities"
	"github.com/sagotly/protoFlex.git/src/repo"
	Ui "github.com/sagotly/protoFlex.git/src/ui"
	"github.com/sagotly/protoFlex.git/src/utils"
)

// Just for demonstration, if we wanted dynamic stats updates, we could create a function
// that updates them on some interval, but here it's static for simplicity.

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

	//create server repo
	serverRepo := repo.NewServerRepo(db)
	err = serverRepo.CreateServer(enteties.Server{
		Ip:         "1.2.3.4.",
		Name:       "Server1",
		TunnelList: `["tun0", "wg0"]`,
	})
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	serverInstance, err := serverRepo.GetServerById(1)
	if err != nil {
		log.Fatalf("Error getting server by id: %v", err)
	}
	//create tunnel repo
	tunnelRepo := repo.NewTunnelRepo(db)
	err = tunnelRepo.CreateTunnel(enteties.Tunnel{
		ServerId:      serverInstance.Id,
		InterfaceName: "tun0",
	})
	if err != nil {
		log.Fatalf("Error creating tunnel: %v", err)
	}
	fmt.Println("Server: ", serverInstance)

	a := app.New()
	w := a.NewWindow("Protoflex")

	Ui := Ui.NewUI(w, serverRepo, tunnelRepo)
	mainContent, err := Ui.BuildUi()
	if err != nil {
		log.Fatalf("Error building UI: %v", err)
	}

	Ui.FyneWindow.SetContent(mainContent)
	Ui.FyneWindow.Resize(fyne.NewSize(1000, 600))
	Ui.FyneWindow.ShowAndRun()
}
