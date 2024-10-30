package config

import "fmt"

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	Library map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	err := s.point.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unable to set user via pointer: %v", err)
	}

	fmt.Println("User has been set")
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
		return fmt.Errorf("%v not found in commands", cmd.Arguments)
	}

	return c(s, cmd)
}
