package repo

import (
	"database/sql"

	e "github.com/sagotly/protoFlex.git/src/entities"
)

type ServerRepo struct {
	db *sql.DB
}

func NewServerRepo(db *sql.DB) *ServerRepo {
	return &ServerRepo{db: db}
}

func (s *ServerRepo) CreateServer(server e.Server) error {
	_, err := s.db.Exec("INSERT INTO servers (ip, name, tunnel_list) VALUES ($1, $2, $3)", server.Ip, server.Name, server.TunnelList)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerRepo) GetServerById(id int64) (e.Server, error) {
	var server e.Server
	err := s.db.QueryRow("SELECT * FROM servers WHERE id = $1", id).Scan(&server.Id, &server.Ip, &server.Name, &server.TunnelList)
	if err != nil {
		return server, err
	}
	return server, nil
}

func (s *ServerRepo) GetAllServers() ([]e.Server, error) {
	var servers []e.Server
	rows, err := s.db.Query("SELECT * FROM servers")
	if err != nil {
		return servers, err
	}
	defer rows.Close()
	for rows.Next() {
		var server e.Server
		err := rows.Scan(&server.Id, &server.Ip, &server.Name, &server.TunnelList)
		if err != nil {
			return servers, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}
