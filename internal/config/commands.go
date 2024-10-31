package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ScooballyD/gator/internal/database"
	"github.com/google/uuid"
)

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	Library map[string]func(*State, Command) error
}

func HandlerGetUsers(s *State, cmd Command) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("the users handler takes no arguments")
	}

	usrs, err := s.dbq.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to retrieve users: %v", err)
	}

	for _, usr := range usrs {
		if usr == s.point.Current_user_name {
			fmt.Printf("%v (current)\n", usr)
		} else {
			fmt.Println(usr)
		}
	}
	return nil
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	_, err := s.dbq.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unknown user: %v", err)
	}

	err = s.point.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unable to set user via pointer: %v", err)
	}
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("the register handler expects a single argument, a name")
	}

	usr, err := s.dbq.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Arguments[0],
		})
	if err != nil {
		return fmt.Errorf("unable to register %v: %v", cmd.Arguments[0], err)
	}

	err = s.point.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unable to set user via pointer: %v", err)
	}
	fmt.Printf("user has been created: %v\n", usr)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Arguments) > 0 {
		return errors.New("too many arguments, reset takes none")
	}

	err := s.dbq.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to reset users table: %v", err)
	}
	return nil
}

// Registers new command into library
func (cmds Commands) Register(name string, f func(*State, Command) error) {
	cmds.Library[name] = f

	_, exist := cmds.Library[name]
	if !exist {
		fmt.Printf("failed to add %v to commands\n", name)
	}
}

func (cmds Commands) Run(s *State, cmd Command) error {
	c, exist := cmds.Library[cmd.Name]
	if !exist {
		return fmt.Errorf("%v not found in commands", cmd.Name)
	}

	return c(s, cmd)
}
