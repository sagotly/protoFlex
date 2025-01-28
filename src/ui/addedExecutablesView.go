package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (ui *Ui) buildAddedExecutablesView() (*fyne.Container, error) {
	// Load the data from the database
	allExecutables, err := ui.AddedExecutablesRepo.GetAllAddedExecutabless()
	if err != nil {
		return nil, err
	}

	// Create the list widget
	executableList := widget.NewList(
		func() int {
			return len(allExecutables)
		},
		func() fyne.CanvasObject {
			// Each item has a title, description, and a "Connect" button
			title := widget.NewLabel("")
			desc := widget.NewLabel("")
			connectButton := widget.NewButton("Connect", nil)
			return container.NewBorder(nil, connectButton, nil, nil, container.NewVBox(title, desc))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			executable := allExecutables[i]
			c := o.(*fyne.Container)

			// Extract components from the container
			vbox := c.Objects[0].(*fyne.Container)
			titleLabel := vbox.Objects[0].(*widget.Label)
			descLabel := vbox.Objects[1].(*widget.Label)
			connectButton := c.Objects[1].(*widget.Button)

			// Set item data
			titleLabel.SetText(executable.Path)
			descLabel.SetText(fmt.Sprintf("Active: %v", executable.Active))

			// Define button behavior
			connectButton.OnTapped = func() {

				dialog.ShowForm("Connect Executable", "Connect", "Cancel", nil, func(ok bool) {
					if ok {
						// Perform the connect action
						err := ui.AddedExecutablesController.ClickOnExecutableBtn(executable.TunnelId, executable.Path, executable.Arguments)
						if err != nil {
							dialog.ShowError(fmt.Errorf("Failed to connect executable: %w", err), ui.FyneWindow)
							return
						}
						dialog.ShowInformation("Success", "Executable connected successfully!", ui.FyneWindow)
					}
				}, ui.FyneWindow)
			}
		},
	)

	// Button to add a new executable
	addExecutableBtn := widget.NewButton("+", func() {
		// Create input fields for the executable information
		pathEntry := widget.NewEntry()
		pathEntry.SetPlaceHolder("Absolute path to executable")

		argsEntry := widget.NewEntry()
		argsEntry.SetPlaceHolder("Arguments (optional)")

		tunnels, err := ui.TunnelRepo.GetAllTunnels()
		if err != nil {
			dialog.ShowError(err, ui.FyneWindow)
			return
		}

		// Create dropdown options for available tunnels
		tunnelOptions := make([]string, len(tunnels))
		for i, tunnel := range tunnels {
			tunnelOptions[i] = tunnel.InterfaceName
		}
		tunnelSelect := widget.NewSelect(tunnelOptions, nil)

		formItems := []*widget.FormItem{
			widget.NewFormItem("Path:", pathEntry),
			widget.NewFormItem("Arguments:", argsEntry),
			widget.NewFormItem("Interface Name:", tunnelSelect),
		}

		// Show a form dialog for adding a new executable
		dialog.ShowForm("Add Executable", "Add", "Cancel", formItems, func(ok bool) {
			if ok {
				args := strings.Split(argsEntry.Text, " ")
				err := ui.AddedExecutablesController.AddExecutableBtn(pathEntry.Text, args, tunnelSelect.Selected)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed to add executable: %w", err), ui.FyneWindow)
					return
				}

				// Refresh the list by reloading data from the database
				allExecutables, err = ui.AddedExecutablesRepo.GetAllAddedExecutabless()
				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed to refresh executables list: %w", err), ui.FyneWindow)
					return
				}
				executableList.Refresh()
			}
		}, ui.FyneWindow)
	})

	// Create the layout for the executables view
	executablesView := container.NewBorder(nil, addExecutableBtn, nil, nil, executableList)
	return executablesView, nil
}
