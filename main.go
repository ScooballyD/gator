package main

import (
	"fmt"

	"github.com/ScooballyD/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("Casey")
	cfg = config.Read()

	fmt.Printf("db url= %v\ncurrent user= %v\n", cfg.Db_url, cfg.Current_user_name)
}
