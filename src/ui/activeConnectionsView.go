package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sagotly/protoFlex.git/src/system_scripts"
)

func (Ui *Ui) buildActiveConnectionsView() (*fyne.Container, error) {
	connections, err := system_scripts.GetActiveConnections()
	if err != nil {
		log.Println(`Failed to get active connections:`, err)
	}

	allTunnels, err := Ui.TunnelsRepo.GetAllTunnels()
	if err != nil {
		log.Println(`Failed to get all tunnels:`, err)
		return nil, err
	}
	options := make([]string, len(allTunnels))
	for i, tunnel := range allTunnels {
		options[i] = fmt.Sprint(tunnel.InterfaceName)
	}

	rightList := widget.NewList(
		func() int {
			return len(connections)
		},
		func() fyne.CanvasObject {
			title := widget.NewLabel("")
			desc := widget.NewLabel("")
			plusBtn := widget.NewButton("+", nil)
			return container.NewBorder(nil, nil, nil, plusBtn, container.NewVBox(title, desc))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			item := connections[i]
			co := o.(*fyne.Container)
			// Container layout: container.NewBorder(nil, nil, nil, plusBtn, container.NewVBox(title, desc))
			// Objects: VBox(title, desc), plusBtn
			vbox := co.Objects[0].(*fyne.Container)
			titleLabel := vbox.Objects[0].(*widget.Label)
			descLabel := vbox.Objects[1].(*widget.Label)
			plusBtn := co.Objects[1].(*widget.Button)

			titleLabel.SetText(fmt.Sprint(item.Id))
			descLabel.SetText(item.Pid)

			plusBtn.OnTapped = func() {
				optionSelect := widget.NewSelect(options, nil)
				dialog.ShowForm("Choose Option", "Confirm", "Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("Option:", optionSelect),
					}, func(ok bool) {
						if ok {
							err := Ui.TunnelsRepo.AddConnectionToTunnel(optionSelect.Selected, item.FiveTuple)
							if err != nil {
								log.Println(`Failed to add connection to tunnel:`, err)
								dialog.ShowError(err, Ui.FyneWindow)
							}
						}
					}, Ui.FyneWindow)
			}
		},
	)

	rightCol := container.NewGridWithColumns(1, rightList)

	return rightCol, nil
}
