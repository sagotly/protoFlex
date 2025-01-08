package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	enteties "github.com/sagotly/protoFlex.git/src/entities"
)

func (u *Ui) buildServerView() (*fyne.Container, error) {
	// Load the data from the db
	allServersFromDB, err := u.ServerRepo.GetAllServers()
	if err != nil {
		return nil, err
	}

	// Create the server list widget
	serverList := widget.NewList(
		func() int {
			return len(allServersFromDB)
		},
		func() fyne.CanvasObject {
			// Each item is a container with name and description labels
			title := widget.NewLabel("")
			desc := widget.NewLabel("")
			return container.NewVBox(title, desc)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			s := allServersFromDB[i]
			c := o.(*fyne.Container)
			title := c.Objects[0].(*widget.Label)
			desc := c.Objects[1].(*widget.Label)

			title.SetText(s.Ip)
			desc.SetText(fmt.Sprintf("IP: %s | Tunnels: %s", s.Ip, s.TunnelList))
		},
	)

	// Button to add a new server
	addServerBtn := widget.NewButton("+", func() {
		// Create input fields for the server information
		nameEntry := widget.NewEntry()
		ipEntry := widget.NewEntry()
		wantedEntry := widget.NewEntry()

		formItems := []*widget.FormItem{
			widget.NewFormItem("Name of the server:", nameEntry),
			widget.NewFormItem("IP of the server:", ipEntry),
			widget.NewFormItem("Wanted tunnels:", wantedEntry),
		}

		// Show a form dialog for adding a server
		dialog.ShowForm("Add Server", "Confirm", "Cancel", formItems, func(ok bool) {
			if ok {
				newItem := enteties.Server{
					Ip:         ipEntry.Text,
					Name:       nameEntry.Text,
					TunnelList: wantedEntry.Text,
				}

				// Save the server to the database
				err := u.ServerRepo.CreateServer(newItem)
				if err != nil {
					dialog.ShowError(err, u.FyneWindow)
					return
				}

				// Refresh the server list after adding a new server
				serverList.Refresh()
			}
		}, u.FyneWindow)
	})

	// Create the server view layout
	serverView := container.NewBorder(nil, addServerBtn, nil, nil, serverList)
	return serverView, nil

}
