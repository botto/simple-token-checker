package main

type config struct {
	ListenPort int    `env:"LISTEN_PORT" envDefault:"3000"`
	ListenAddr string `env:"LISTEN_ADDR" envDefault:"0.0.0.0"`
	VaultToken string `env:"VAULT_TOKEN"`
	CAFilePath string `env:"CA_FILE_PATH"`
	VaultAddr  string `env:"VAULT_ADDR" envDefault:"http://127.0.0.1:8200"`
}
