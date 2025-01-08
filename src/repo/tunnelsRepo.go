package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"

	enteties "github.com/sagotly/protoFlex.git/src/entities"
)

// TunnelRepo is a struct that contains a pointer to the database
type TunnelRepo struct {
	db *sql.DB
}

// NewTunnelRepo is a function that returns a pointer to a TunnelRepo struct
func NewTunnelRepo(db *sql.DB) *TunnelRepo {
	return &TunnelRepo{db: db}
}

// CreateTunnel is a function that inserts a new tunnel into the database
func (t *TunnelRepo) CreateTunnel(tunnel enteties.Tunnel) error {
	_, err := t.db.Exec("INSERT INTO tunnels (server_id, interface_name, connected_connections) VALUES ($1, $2, $3)", tunnel.ServerId, tunnel.InterfaceName, tunnel.ConnectedConnections)
	if err != nil {
		return err
	}
	return nil
}

// GetTunnelById is a function that retrieves a tunnel from the database by its id
func (t *TunnelRepo) GetTunnelById(id int64) (enteties.Tunnel, error) {
	var tunnel enteties.Tunnel
	err := t.db.QueryRow("SELECT * FROM tunnels WHERE id = $1", id).Scan(&tunnel.Id, &tunnel.ServerId, &tunnel.InterfaceName, &tunnel.ConnectedConnections)
	if err != nil {
		return tunnel, err
	}
	return tunnel, nil
}

// GetTunnelByInterfaceName is a function that retrieves a tunnel from the database by its interface name
func (t *TunnelRepo) getTunnelByInterfaceName(name string) (enteties.Tunnel, error) {
	var tunnel enteties.Tunnel
	err := t.db.QueryRow("SELECT * FROM tunnels WHERE interface_name = $1", name).Scan(&tunnel.Id, &tunnel.ServerId, &tunnel.InterfaceName, &tunnel.ConnectedConnections)
	if err != nil {
		return tunnel, err
	}
	return tunnel, nil
}

// GetAllTunnels is a function that retrieves all tunnels from the database
func (t *TunnelRepo) GetAllTunnels() ([]enteties.Tunnel, error) {
	var tunnels []enteties.Tunnel
	rows, err := t.db.Query("SELECT * FROM tunnels")
	if err != nil {
		return tunnels, err
	}
	defer rows.Close()
	for rows.Next() {
		var tunnel enteties.Tunnel
		err := rows.Scan(&tunnel.Id, &tunnel.ServerId, &tunnel.InterfaceName, &tunnel.ConnectedConnections)
		if err != nil {
			return tunnels, err
		}
		tunnels = append(tunnels, tunnel)
	}
	return tunnels, nil
}

func (t *TunnelRepo) AddConnectionToTunnel(tunnelName string, connectionTuple string) error {
	// Получение туннеля с использованием существующей функции
	tunnel, err := t.getTunnelByInterfaceName(tunnelName)
	if err != nil {
		return fmt.Errorf("failed to fetch tunnel by ID: %w", err)
	}

	// Преобразование строки JSON в Go срез
	var connections []string
	if tunnel.ConnectedConnections == "[]" {
		connections = []string{}
	} else {
		if err := json.Unmarshal([]byte(tunnel.ConnectedConnections), &connections); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}
	}

	// Добавление нового PID в массив
	connections = append(connections, connectionTuple)

	// Преобразование обратно в JSON
	updatedConnections, err := json.Marshal(connections)
	if err != nil {
		return fmt.Errorf("failed to serialize connections: %w", err)
	}

	// Обновление базы данных с новым массивом
	_, err = t.db.Exec(`
        UPDATE tunnels 
        SET connected_connections = $1
        WHERE interface_name = $2;
    `, string(updatedConnections), tunnelName)
	if err != nil {
		return fmt.Errorf("failed to update connections: %w", err)
	}

	return nil
}
