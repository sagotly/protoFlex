package handlers

import (
	"github.com/sagotly/protoFlex.git/src/entities"
	"github.com/sagotly/protoFlex.git/src/repo"
	"github.com/sagotly/protoFlex.git/src/system_scripts"
)

type namespaceHandler struct {
	tunnelRepo *repo.TunnelRepo
}

func (n *namespaceHandler) createNewNamespace(interface_name string, server enteties.Server) error {
	err := system_scripts.SetupNamespace(interface_name)
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
