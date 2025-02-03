package enteties

type validateTokenResponse struct {
	WgConfig string `json:"wg_config"`
}

type GenerateTokenRequest struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type ValidateTokenRequest struct {
	Ip    string `json:"ip"`
	Port  string `json:"port"`
	Token string `json:"token"`
}
