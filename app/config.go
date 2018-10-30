package app

type Config struct {
	Port     int             `json:"port"`
	Operator *OperatorConfig `json:"operator"`
}

type OperatorConfig struct {
	PrivateKey string `json:"privkey"`
}
