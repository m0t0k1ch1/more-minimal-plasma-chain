package app

type Config struct {
	Port      int              `json:"port"`
	Operator  *OperatorConfig  `json:"operator"`
	RootChain *RootChainConfig `json:"rootchain"`
}

type OperatorConfig struct {
	PrivateKey string `json:"privkey"`
}

type RootChainConfig struct {
	RPC     string `json:"rpc"`
	Address string `json:"address"`
}
