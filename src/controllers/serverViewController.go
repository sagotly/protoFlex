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
	_, exists, err := n.tunnelRepo.GetTunnelByInterfaceName(interface_name)
	if err != nil {
		log.Println("Error creating server:", err)
	}

	id, err := n.serverRepo.CreateServer(server)
	if err != nil {
		log.Println("Error creating server:", err)
		return err
	}

	err = system_scripts.SetupNamespace(interface_name)
	if err != nil {
		return err
	}

	tunnel := enteties.Tunnel{
		ServerId:      id,
		InterfaceName: interface_name,
	}

	if !exists {
		err = n.tunnelRepo.CreateTunnel(tunnel)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *ServerViewController) GetAllServers() ([]enteties.Server, error) {
	servers, err := n.serverRepo.GetAllServers()
	if err != nil {
		log.Println("Error getting servers:", err)
		return nil, err
	}
	return servers, nil
}
