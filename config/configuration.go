package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type MySQLConfig struct {
	DBUser   string `json:"db_user"`
	DBPass   string `json:"db_pass"`
	DBServer string `json:"db_server"`
	DBPort   string `json:"db_port"`
	DBName   string `json:"db_name"`
}

type URLConfig struct {
	URLBase string `json:"url_base"`
	URLPort string `json:"url_port"`
}

type Configuration struct {
	MySQL MySQLConfig `json:"mysql"`
	URL   URLConfig   `json:"url"`
}

var AppConfig Configuration

func LoadConfig() error {
	configFile, err := os.Open("production.json")
	if err != nil {
		return err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := LoadConfig()
	if err != nil {
		fmt.Println("Error al cargar la configuración:", err)
		return
	}

}
