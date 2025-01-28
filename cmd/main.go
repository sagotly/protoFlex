package main

import (
	"database/sql"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sagotly/protoFlex.git/src/client"
	"github.com/sagotly/protoFlex.git/src/controllers"
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

	serverRepo := repo.NewServerRepo(db)
	tunnelRepo := repo.NewTunnelRepo(db)
	addedExecutablesRepo := repo.NewAddedExecutablesRepo(db)

	tokenClient := client.NewServerClient()

	tokenController := controllers.NewTokenController(tokenClient)
	serverViewController := controllers.NewServerViewController(tunnelRepo, serverRepo)
	addedExecutablesController := controllers.NewAddedExcecutablesController(tunnelRepo, serverRepo, addedExecutablesRepo)

	a := app.New()
	w := a.NewWindow("Protoflex")
	Ui := Ui.NewUI(w, tokenController, serverViewController, addedExecutablesController, serverRepo, tunnelRepo, addedExecutablesRepo)

	mainContent, err := Ui.BuildUi()
	if err != nil {
		log.Fatalf("Error building UI: %v", err)
	}

	Ui.FyneWindow.SetContent(mainContent)
	Ui.FyneWindow.Resize(fyne.NewSize(1000, 600))
	Ui.FyneWindow.ShowAndRun()
}
