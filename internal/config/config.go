package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

// Creates Config struct from "~/.gatorconfig.json"
func Read() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to find home dir : ", err)
		return Config{}
	}

	path := home + "/" + configFileName

	jfile, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Failed to access 'gatorconfig.json' : ", err)
		return Config{}
	}

	newConfig := Config{}
	err = json.Unmarshal(jfile, &newConfig)
	if err != nil {
		fmt.Println("Unmarshal error : ", err)
		return Config{}
	}

	return newConfig
}

func (cfg Config) SetUser(user string) {
	cfg.Current_user_name = user

	jData, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println("Marshal error : ", err)
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to find home dir : ", err)
		return
	}

	path := home + "/" + configFileName
	err = os.WriteFile(path, jData, 0666)
	if err != nil {
		fmt.Println("Write error : ", err)
		return
	}
}
