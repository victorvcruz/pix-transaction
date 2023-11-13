package cfg

import "os"

type DynamodbConfig struct {
	AccountTable string
	UserTable    string
	KeyTable     string
	PixTable     string
}

type Config struct {
	DynamodbConfig DynamodbConfig
}

func Get() (*Config, error) {
	return &Config{
		DynamodbConfig{
			AccountTable: os.Getenv("ACCOUNT_TABLE"),
			UserTable:    os.Getenv("USER_TABLE"),
			KeyTable:     os.Getenv("KEY_TABLE"),
			PixTable:     os.Getenv("PIX_TABLE"),
		},
	}, nil
}
