package enteties

type AddedExecutable struct {
	Id        int64  `db:"id"`
	TunnelId  int64  `db:"tunnel_id"`
	Path      string `db:"path"`
	Arguments string `db:"arguments"`
	Active    bool   `db:"active"`
}
