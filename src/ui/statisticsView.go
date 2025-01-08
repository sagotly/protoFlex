package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sagotly/protoFlex.git/src/system_scripts"
)

func (Ui *Ui) buildStatisticsView() (*fyne.Container, error) {
	titleLabel := widget.NewLabel("Protoflex")
	titleLabel.TextStyle.Bold = true

	connections, err := system_scripts.GetActiveConnections()
	if err != nil {
		log.Fatal("Failed to get active connections:", err)
	}
	activeConnections := widget.NewLabel(fmt.Sprintf("Active connections: %d", len(connections)))

	servers, err := Ui.ServerRepo.GetAllServers()
	if err != nil {
		log.Println("Failed to get all servers:", err)
		return nil, err
	}

	numServers := widget.NewLabel(fmt.Sprintf("Servers: %d", len(servers)))

	middleCol := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		activeConnections,
		numServers,
	)

	return middleCol, nil

}
