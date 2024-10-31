package main

import (
	"fmt"
	"os"

	"github.com/ScooballyD/gator/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s, err := cfg.NewState()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmds := config.Commands{
		Library: make(map[string]func(*config.State, config.Command) error),
	}
	cmds.Register("login", config.HandlerLogin)
	cmds.Register("register", config.HandlerRegister)
	cmds.Register("reset", config.HandlerReset)
	cmds.Register("users", config.HandlerGetUsers)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("error: not enough arguments")
		os.Exit(1)
	}

	cmd := config.Command{
		Name:      args[1],
		Arguments: args[2:],
	}
	err = cmds.Run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
