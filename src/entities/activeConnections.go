package enteties

type ActiveConnection struct {
	Id int64 `db:"id"`
	// TunnelId         int64  `db:"tunnel_id"`
	Pid              string `db:"pid"`
	TypeOfConnection string `db:"type_of_connection"`
	FiveTuple        string `db:"five_tuple"` // Добавлено новое поле для хранения 5-tuple
}

type ActiveConnections []*ActiveConnection
