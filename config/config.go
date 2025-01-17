package config

import (
	"fmt"
	"os"

	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/nowpayments"
	"github.com/pagu-project/pagu/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BotName      string              `yaml:"bot_name"`
	Network      string              `yaml:"network"`
	NetworkNodes []string            `yaml:"network_nodes"`
	LocalNode    string              `yaml:"local_node"`
	Database     Database            `yaml:"database"`
	GRPC         *GRPC               `yaml:"grpc"` // ! TODO: config for modules should moved to the module.
	Wallet       *Wallet             `yaml:"wallet"`
	Logger       *log.Config         `yaml:"logger"`
	HTTP         *HTTP               `yaml:"http"`
	Phoenix      *PhoenixNetwork     `yaml:"phoenix"`
	Discord      *DiscordBot         `yaml:"discord"`
	Telegram     *Telegram           `yaml:"telegram"`
	Notification *Notification       `yaml:"notification"`
	NowPayments  *nowpayments.Config `yaml:"now_payments"`
}

type Database struct {
	URL string `yaml:"url"`
}

type Wallet struct {
	Address  string        `yaml:"address"`
	Path     string        `yaml:"path"`
	Password string        `yaml:"password"`
	Fee      amount.Amount `yaml:"fee"`
}

type DiscordBot struct {
	Token   string `yaml:"token"`
	GuildID string `yaml:"guild_id"`
}

type GRPC struct {
	Listen string `yaml:"listen"`
}

type HTTP struct {
	Listen string `yaml:"listen"`
}

type PhoenixNetwork struct {
	FaucetAmount amount.Amount `yaml:"faucet_amount"`
}

type Telegram struct {
	BotToken string `yaml:"bot_token"`
}

type Notification struct {
	Zoho *Zoho `yaml:"zoho"`
}

type Zoho struct {
	Mail ZapToMail `yaml:"mail"`
}

type ZapToMail struct {
	Host      string            `yaml:"host"`
	Port      int               `yaml:"port"`
	Username  string            `yaml:"username"`
	Password  string            `yaml:"password"`
	Templates map[string]string `yaml:"templates"`
}

func Load(path string) (*Config, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(payload, cfg); err != nil {
		return nil, err
	}

	// Check if the required configurations are set
	if err := cfg.BasicCheck(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// BasicCheck validate presence of required config variables.
func (cfg *Config) BasicCheck() error {
	if cfg.Wallet.Address == "" {
		return fmt.Errorf("config: Wallet address dose not set")
	}

	// Check if the WalletPath exists.
	if !utils.PathExists(cfg.Wallet.Path) {
		return fmt.Errorf("config: Wallet does not exist: %s", cfg.Wallet.Path)
	}

	if len(cfg.NetworkNodes) == 0 {
		return fmt.Errorf("config: network nodes is empty")
	}

	return nil
}
