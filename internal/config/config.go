package config

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ScooballyD/gator/internal/database"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

type State struct {
	dbq   *database.Queries
	point *Config
}

func (cfg Config) NewState() (State, error) {
	s := State{
		point: &cfg,
	}

	if s.point == nil {
		return State{}, errors.New("failed to create new state")
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		return State{}, fmt.Errorf("failed to open database: %v", err)
	}
	dbQueries := database.New(db)
	s.dbq = dbQueries
	if s.dbq == nil {
		return State{}, errors.New("failed to assign dbQueries to state")
	}

	return s, nil
}

// Creates Config struct from "~/.gatorconfig.json"
func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		//fmt.Println("Unable to find home dir : ", err)
		return Config{}, fmt.Errorf("unable to find home dir: %v", err)
	}

	path := home + "/" + configFileName

	jfile, err := os.ReadFile(path)
	if err != nil {
		//fmt.Println("Failed to access 'gatorconfig.json' : ", err)
		return Config{}, fmt.Errorf("failed to access 'gatorconfig.json': %v", err)
	}

	newConfig := Config{}
	err = json.Unmarshal(jfile, &newConfig)
	if err != nil {
		//fmt.Println("Unmarshal error : ", err)
		return Config{}, fmt.Errorf("unmarshal error: %v", err)
	}

	return newConfig, nil
}

func (cfg Config) SetUser(user string) error {
	cfg.Current_user_name = user

	jData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal error: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to find home dir: %v", err)
	}

	path := home + "/" + configFileName
	err = os.WriteFile(path, jData, 0666)
	if err != nil {
		return fmt.Errorf("write error: %v", err)
	}

	return nil
}
