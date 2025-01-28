package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (Ui *Ui) buildStatisticsView() (*fyne.Container, error) {
	titleLabel := widget.NewLabel("Protoflex")
	titleLabel.TextStyle.Bold = true

	executables, err := Ui.AddedExecutablesRepo.GetAllAddedExecutabless()
	if err != nil {
		log.Fatal("Failed to get active connections:", err)
	}
	addedExecutables := widget.NewLabel(fmt.Sprintf("Active connections: %d", len(executables)))

	servers, err := Ui.ServerRepo.GetAllServers()
	if err != nil {
		log.Println("Failed to get all servers:", err)
		return nil, err
	}

	numServers := widget.NewLabel(fmt.Sprintf("Servers: %d", len(servers)))

	middleCol := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		addedExecutables,
		numServers,
	)

	return middleCol, nil

}
