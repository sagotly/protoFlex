package enteties

type AddedExecutable struct {
	Id        int64  `db:"id"`
	TunnelId  int64  `db:"tunnel_id"`
	Path      string `db:"path"`
	Arguments string `db:"arguments"`
	Active    bool   `db:"active"`
}

type AddExecutableRequest struct {
	Path      string `json:"path"`
	Arguments string `json:"arguments"`
	TunnelId  string `json:"tunnel_id"`
}

type ConnectExecutableRequest struct {
	TunnelId  int64  `json:"tunnel_id"`
	Path      string `json:"path"`
	Arguments string `json:"arguments"`
}
