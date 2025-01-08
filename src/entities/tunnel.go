package enteties

// Tunnel is a struct that represents a tunnel
type Tunnel struct {
	Id                   int64  `db:"id"`
	ServerId             int64  `db:"server_id"`
	InterfaceName        string `db:"interface_name"`
	ConnectedConnections string `db:"connected_connections"`
}
