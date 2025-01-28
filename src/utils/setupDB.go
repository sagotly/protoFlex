package utils

import "database/sql"

func SetupDatabase(db *sql.DB) error {
	// Create the 'servers' table
	serversTable := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		name TEXT NOT NULL,
		tunnel_list TEXT DEFAULT '[]'
	);`
	if _, err := db.Exec(serversTable); err != nil {
		return err
	}

	// Create the 'tunnels' table
	tunnelsTable := `
	CREATE TABLE IF NOT EXISTS tunnels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		server_id INTEGER NOT NULL,
		interface_name TEXT NOT NULL,
		connected_connections TEXT DEFAULT '[]',
		FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE
	);`
	if _, err := db.Exec(tunnelsTable); err != nil {
		return err
	}

	// Create the 'AddedExecutables' table
	AddedExecutablesTable := `
	CREATE TABLE IF NOT EXISTS AddedExecutables (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tunnel_id INTEGER NOT NULL,
		path TEXT NOT NULL,
		arguments TEXT DEFAULT '',
		active BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (tunnel_id) REFERENCES tunnels(id) ON DELETE CASCADE
	);`
	if _, err := db.Exec(AddedExecutablesTable); err != nil {
		return err
	}

	return nil
}
