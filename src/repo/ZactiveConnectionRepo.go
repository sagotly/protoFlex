package repo

// import (
// 	"database/sql"

// 	enteties "github.com/sagotly/protoFlex.git/src/entities"
// )

// type ActiveConnectionRepo struct {
// 	db *sql.DB
// }

// func NewActiveConnectionRepo(db *sql.DB) *ActiveConnectionRepo {
// 	return &ActiveConnectionRepo{db: db}
// }

// func (a *ActiveConnectionRepo) CreateActiveConnection(activeConnection enteties.ActiveConnection) error {
// 	_, err := a.db.Exec("INSERT INTO active_connections (tunnel_id, pid, type_of_connection) VALUES ($1, $2, $3)", activeConnection.TunnelId, activeConnection.Pid, activeConnection.TypeOfConnection)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (a *ActiveConnectionRepo) GetActiveConnectionById(id int64) (enteties.ActiveConnection, error) {
// 	var activeConnection enteties.ActiveConnection
// 	err := a.db.QueryRow("SELECT * FROM active_connections WHERE id = $1", id).Scan(&activeConnection.Id, &activeConnection.TunnelId, &activeConnection.Pid, &activeConnection.TypeOfConnection)
// 	if err != nil {
// 		return activeConnection, err
// 	}
// 	return activeConnection, nil
// }
