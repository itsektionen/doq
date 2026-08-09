package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	StorageType           string `mapstructure:"storage_type"`
	FilePath              string `mapstructure:"file_path"`
	GDriveCredentialsFile string `mapstructure:"gdrive_credentials_file"`
	GDriveFolderID        string `mapstructure:"gdrive_folder_id"`
	Port                  string `mapstructure:"port"`
}

func (c *Config) Validate() error {
	switch strings.ToLower(c.StorageType) {
	case "gdrive", "googledrive", "drive":
		if c.GDriveCredentialsFile == "" {
			return errors.New("gdrive_credentials_file is required when storage_type is 'gdrive'")
		}
	case "file", "local":
		if c.FilePath == "" {
			c.FilePath = filepath.Join("tmp", "img")
		}
	case "nop", "none", "":
		// valid
	default:
		return fmt.Errorf("unknown storage_type: %s (must be 'gdrive', 'file', or 'nop')", c.StorageType)
	}

	if c.Port == "" {
		c.Port = ":8080"
	}

	return nil
}

func Load() (Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("storage_type", "nop")
	v.SetDefault("file_path", "tmp/img")
	v.SetDefault("port", ":8080")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}
