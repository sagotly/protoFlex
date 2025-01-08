package enteties

type Server struct {
	Id         int64  `db:"id"`
	Name       string `db:"name"`
	Ip         string `db:"ip"`
	TunnelList string `db:"tunnel_list"`
}
