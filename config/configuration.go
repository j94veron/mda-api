package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	MySQLConfig MySQLConfig `json:"mysql_config"`
	MDAConfig   MDAConfig   `json:"mda_config"`
}

type MySQLConfig struct {
	DBUser   string `json:"db_user"`
	DBPass   string `json:"db_pass"`
	DBServer string `json:"db_server"`
	DBPort   string `json:"db_port"`
	DBName   string `json:"db_name"`
}

type MDAConfig struct {
	URLBase string `json:"url_base"`
	Path    string `json:"url_path"`
}

func LoadConfig(config *Configuration) error {
	configFile, err := os.Open("config/production.json")
	if err != nil {
		return err
	}

	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}

	return nil
}
