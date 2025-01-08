package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/sagotly/protoFlex.git/src/repo"
)

type Ui struct {
	FyneWindow fyne.Window

	ServerRepo  *repo.ServerRepo
	TunnelsRepo *repo.TunnelRepo
}

func NewUI(fyneWindow fyne.Window, serverRepo *repo.ServerRepo, tunnelsRepo *repo.TunnelRepo) *Ui {
	return &Ui{
		FyneWindow:  fyneWindow,
		ServerRepo:  serverRepo,
		TunnelsRepo: tunnelsRepo,
	}
}

func (Ui *Ui) BuildUi() (*fyne.Container, error) {
	serverViewContainer, err := Ui.buildServerView()
	if err != nil {
		log.Println("Error building server view:", err)
		return nil, err
	}

	statisticsViewContainer, err := Ui.buildStatisticsView()
	if err != nil {
		log.Println("Error building statistics view:", err)
		return nil, err
	}

	activeConnectionViewContainer, err := Ui.buildActiveConnectionsView()
	if err != nil {
		log.Println("Error building active connections view:", err)
		return nil, err
	}

	content := container.NewGridWithColumns(3,
		serverViewContainer,
		statisticsViewContainer,
		activeConnectionViewContainer,
	)
	return content, nil
}
