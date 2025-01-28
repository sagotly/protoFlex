package controllers

import (
	"github.com/sagotly/protoFlex.git/src/client"
)

type TokenController struct {
	client *client.ServerClient
}

func NewTokenController(client *client.ServerClient) *TokenController {
	return &TokenController{
		client: client,
	}
}

func (t *TokenController) GenerateToken(ip string, port string) error {
	err := t.client.GenerateToken(ip, port)
	if err != nil {
		return err
	}

	return nil
}

func (t *TokenController) ValidateToken(ip string, port string, token string) error {
	err := t.client.ValidateToken(ip, port, token)
	if err != nil {
		return err
	}

	return nil
}
