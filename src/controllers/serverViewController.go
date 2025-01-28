package controllers

import (
	"log"

	enteties "github.com/sagotly/protoFlex.git/src/entities"
	"github.com/sagotly/protoFlex.git/src/repo"
	"github.com/sagotly/protoFlex.git/src/system_scripts"
)

type ServerViewController struct {
	serverRepo *repo.ServerRepo
	tunnelRepo *repo.TunnelRepo
}

func NewServerViewController(tunnelRepo *repo.TunnelRepo, serverRepo *repo.ServerRepo) *ServerViewController {
	return &ServerViewController{
		serverRepo: serverRepo,
		tunnelRepo: tunnelRepo,
	}
}

func (n *ServerViewController) CreateNewServerBtn(serverName string, serverIp string, interface_name string) error {
	server := enteties.Server{
		Name:       serverName,
		Ip:         serverIp,
		TunnelList: interface_name,
	}

	err := n.serverRepo.CreateServer(server)
	if err != nil {
		log.Println("Error creating server:", err)
		return err
	}

	err = system_scripts.SetupNamespace(interface_name)
	if err != nil {
		return err
	}

	tunnel := enteties.Tunnel{
		ServerId:      server.Id,
		InterfaceName: interface_name,
	}

	err = n.tunnelRepo.CreateTunnel(tunnel)
	if err != nil {
		return err
	}
	return nil
}
