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
	RPC     *RootChainRPCConfig `json:"rpc"`
	Address string              `json:"address"`
}

type RootChainRPCConfig struct {
	HTTP string `json:"http"`
	WS   string `json:"ws"`
}
