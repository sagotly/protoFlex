package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/sagotly/protoFlex.git/src/controllers"
	"github.com/sagotly/protoFlex.git/src/repo"
)

type Ui struct {
	FyneWindow fyne.Window

	TunnelRepo           *repo.TunnelRepo
	ServerRepo           *repo.ServerRepo
	AddedExecutablesRepo *repo.AddedExecutablesRepo

	AddedExecutablesController *controllers.AddedExecutablesController
	ServerViewController       *controllers.ServerViewController
	tokenController            *controllers.TokenController
}

func NewUI(fyneWindow fyne.Window, tokenController *controllers.TokenController, serverViewController *controllers.ServerViewController, addedExecutablesController *controllers.AddedExecutablesController, serverRepo *repo.ServerRepo, tunnelRepo *repo.TunnelRepo, addedExecutablesRepo *repo.AddedExecutablesRepo) *Ui {
	return &Ui{
		FyneWindow:                 fyneWindow,
		TunnelRepo:                 tunnelRepo,
		ServerRepo:                 serverRepo,
		AddedExecutablesRepo:       addedExecutablesRepo,
		ServerViewController:       serverViewController,
		AddedExecutablesController: addedExecutablesController,
		tokenController:            tokenController,
	}
}

func (u *Ui) BuildUi() (*fyne.Container, error) {
	// Build the containers for the bottom section (two columns)
	serverViewContainer, err := u.buildServerView()
	if err != nil {
		return nil, fmt.Errorf("Error building server view: %v", err)
	}

	addedExecutablesViewContainer, err := u.buildAddedExecutablesView()
	if err != nil {
		return nil, fmt.Errorf("Error building active connections view: %v", err)
	}

	// Create the main container with a two-column layout for the bottom section
	lowerContent := container.NewGridWithColumns(2,
		serverViewContainer,
		addedExecutablesViewContainer,
	)
	lowerContent.MinSize().AddWidthHeight(800, 400) // set the minimum size

	text := canvas.NewText("ProtoFlex", color.Opaque) // set the text and color
	text.TextSize = 100                               // set the desired text size
	text.Alignment = fyne.TextAlignCenter             // set the text alignment to center
	text.TextStyle = fyne.TextStyle{Bold: true}

	// Place the title in the center and the statsContainer on the right side
	topBannerOverlay := container.NewHBox(
		layout.NewSpacer(), // add spacer to push the title away from the left edge
		text,
		layout.NewSpacer(), // add spacer as a separator between the title and statsContainer
	)

	// Combine all the containers into a single container, with the top as the banner and the bottom as the two columns.
	// Use VSplit to specify that the top part occupies 1/4 of the height
	split := container.NewBorder(
		topBannerOverlay,
		nil,
		nil,
		nil,
		lowerContent,
	)

	return split, nil
}
