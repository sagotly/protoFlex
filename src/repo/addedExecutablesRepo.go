package repo

import (
	"database/sql"

	enteties "github.com/sagotly/protoFlex.git/src/entities"
)

type AddedExecutablesRepo struct {
	db *sql.DB
}

func NewAddedExecutablesRepo(db *sql.DB) *AddedExecutablesRepo {
	return &AddedExecutablesRepo{db: db}
}

func (a *AddedExecutablesRepo) CreateAddedExecutable(AddedExecutable enteties.AddedExecutable) error {
	_, err := a.db.Exec("INSERT INTO AddedExecutables (tunnel_id, path, arguments, active) VALUES ($1, $2, $3, $4)", AddedExecutable.TunnelId, AddedExecutable.Path, AddedExecutable.Arguments, AddedExecutable.Active)
	if err != nil {
		return err
	}
	return nil
}

func (a *AddedExecutablesRepo) GetAllAddedExecutabless() ([]enteties.AddedExecutable, error) {
	var AddedExecutables []enteties.AddedExecutable
	rows, err := a.db.Query("SELECT * FROM AddedExecutables")
	if err != nil {
		return AddedExecutables, err
	}
	defer rows.Close()
	for rows.Next() {
		var AddedExecutable enteties.AddedExecutable
		err := rows.Scan(&AddedExecutable.Id, &AddedExecutable.TunnelId, &AddedExecutable.Path, &AddedExecutable.Arguments, &AddedExecutable.Active)
		if err != nil {
			return AddedExecutables, err
		}
		AddedExecutables = append(AddedExecutables, AddedExecutable)
	}
	return AddedExecutables, nil
}
