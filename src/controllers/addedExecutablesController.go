package controllers

import (
	"fmt"
	"strings"

	enteties "github.com/sagotly/protoFlex.git/src/entities"
	"github.com/sagotly/protoFlex.git/src/repo"
	"github.com/sagotly/protoFlex.git/src/system_scripts"
)

type AddedExecutablesController struct {
	serverRepo          *repo.ServerRepo
	tunnelRepo          *repo.TunnelRepo
	addedExecutableRepo *repo.AddedExecutablesRepo
}

func NewAddedExcecutablesController(tunnelRepo *repo.TunnelRepo, serverRepo *repo.ServerRepo, addedExecutablesRepo *repo.AddedExecutablesRepo) *AddedExecutablesController {
	return &AddedExecutablesController{
		serverRepo:          serverRepo,
		tunnelRepo:          tunnelRepo,
		addedExecutableRepo: addedExecutablesRepo,
	}
}

func (ax *AddedExecutablesController) AddExecutableBtn(executablePath string, executableArgs []string, interfaceName string) error {
	tunnels, err := ax.tunnelRepo.GetAllTunnels()
	if err != nil {
		return fmt.Errorf("Error getting tunnels: %w", err)
	}

	wantedTunnel := enteties.Tunnel{}
	for _, tunnel := range tunnels {
		if tunnel.InterfaceName == interfaceName {
			wantedTunnel = tunnel
			break
		}
	}
	executable := enteties.AddedExecutable{
		TunnelId:  wantedTunnel.Id,
		Path:      executablePath,
		Arguments: strings.Join(executableArgs, " "),
		Active:    false,
	}
	err = ax.addedExecutableRepo.CreateAddedExecutable(executable)
	if err != nil {
		return fmt.Errorf("error while creating executable: %w", err)
	}

	return nil
}

func (ax *AddedExecutablesController) ClickOnExecutableBtn(tunnelId int64, path string, args string) error {
	tunnel, err := ax.tunnelRepo.GetTunnelById(tunnelId)
	if err != nil {
		return fmt.Errorf("failed to get tunnel by id %w", err)
	}

	argSlice := strings.Split(args, " ")
	err = system_scripts.RunExecutable(tunnel.InterfaceName+"_namespace", path, argSlice)
	if err != nil {
		return fmt.Errorf("error while clkick on the executable: %w", err)
	}
	return nil
}

func (ax *AddedExecutablesController) GetAllExecutables() ([]enteties.AddedExecutable, error) {
	executables, err := ax.addedExecutableRepo.GetAllAddedExecutabless()
	if err != nil {
		return nil, fmt.Errorf("error while getting all executables: %w", err)
	}
	for i := range executables {
		tunnel, err := ax.tunnelRepo.GetTunnelById(executables[i].TunnelId)
		if err != nil {
			return nil, fmt.Errorf("error while getting tunnel by id: %w", err)
		}
		executables[i].Interface = tunnel.InterfaceName
	}
	fmt.Println("executables: ", executables)
	return executables, nil
}

func (ax *AddedExecutablesController) GetAllTunnels() ([]enteties.Tunnel, error) {
	tunnels, err := ax.tunnelRepo.GetAllTunnels()
	if err != nil {
		return nil, fmt.Errorf("error while getting all tunnels: %w", err)
	}
	return tunnels, nil
}
