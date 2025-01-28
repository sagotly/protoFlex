package client

import (
	"fmt"
	"net/http"
)

// ServerClient defines a client for interacting with the server.
type ServerClient struct {
	httpClient *http.Client
}

// NewServerClient creates a new instance of ServerClient.
func NewServerClient() *ServerClient {
	return &ServerClient{
		httpClient: &http.Client{},
	}
}

// GenerateToken sends a request to the server to generate a token.
// It takes the server's IP and port as arguments and returns the generated token or an error.
func (sc *ServerClient) GenerateToken(serverIP string, serverPort string) error {
	url := fmt.Sprintf("http://%s:%s/generate", serverIP, serverPort)

	resp, err := sc.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send request to generate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to generate token, server returned: %s", resp.Status)
	}
	return nil
}

// ValidateToken sends a request to validate the token on the server.
// It takes the server's IP, port, and the token to be validated.
func (sc *ServerClient) ValidateToken(serverIP string, serverPort string, token string) error {
	url := fmt.Sprintf("http://%s:%s/validate?token=%s", serverIP, serverPort, token)

	resp, err := sc.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send request to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token validation failed, server returned: %s", resp.Status)
	}

	return nil
}
