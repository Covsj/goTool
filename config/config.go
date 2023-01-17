package config

import (
	"runtime"

	"goTool/log"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/ethclient"
)

var DefaultToolConfig *ToolConfig

type Web3Config struct {
	EthRpc      string            `json:"eth_rpc" toml:"eth_rpc"`
	BscRpc      string            `json:"bsc_rpc" toml:"bsc_rpc"`
	MaticRpc    string            `json:"matic_rps" toml:"matic_rpc"`
	EthClient   *ethclient.Client `json:"eth_client"`
	BscClient   *ethclient.Client `json:"bsc_client"`
	MaticClient *ethclient.Client `json:"matic_client"`
}

type ToolConfig struct {
	LogLevel   string `json:"log_level" toml:"log_level"`
	Web3Config *Web3Config
}

func (c *ToolConfig) LoadFromFile(path string) error {
	_, err := toml.DecodeFile(path, c)
	return err
}

func init() {

	DefaultToolConfig = &ToolConfig{}
	_, currentPath, _, _ := runtime.Caller(0)
	currentPath = currentPath[:len(currentPath)-10]

	if err := DefaultToolConfig.LoadFromFile(currentPath + "/goTool.toml"); err != nil {
		log.Error("msg", "init config failed", "error", err.Error())
	}

	log.SetLevel(DefaultToolConfig.LogLevel)

	if DefaultToolConfig != nil && DefaultToolConfig.Web3Config != nil {
		if DefaultToolConfig.Web3Config.EthRpc != "" {
			for i := 0; i < 3; i++ {
				client, err := ethclient.Dial(DefaultToolConfig.Web3Config.EthRpc)
				if err != nil {
					continue
				} else {
					DefaultToolConfig.Web3Config.EthClient = client
					log.Info("msg", "init eth client success")
					break
				}
			}
		}
		if DefaultToolConfig.Web3Config.BscRpc != "" {
			for i := 0; i < 3; i++ {
				client, err := ethclient.Dial(DefaultToolConfig.Web3Config.BscRpc)
				if err != nil {
					continue
				} else {
					DefaultToolConfig.Web3Config.BscClient = client
					log.Info("msg", "init bsc client success")
					break
				}
			}
		}
		if DefaultToolConfig.Web3Config.MaticRpc != "" {
			for i := 0; i < 3; i++ {
				client, err := ethclient.Dial(DefaultToolConfig.Web3Config.MaticRpc)
				if err != nil {
					continue
				} else {
					DefaultToolConfig.Web3Config.MaticClient = client
					log.Info("msg", "init matic client success")
					break
				}
			}
		}

	}

}
