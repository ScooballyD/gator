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

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("addfeed takes two arguments: name, url")
	}

	usr, err := s.dbq.GetUser(context.Background(), s.point.Current_user_name)
	if err != nil {
		return fmt.Errorf("unable to find current user: %v", err)
	}

	feed, err := s.dbq.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Arguments[0],
			Url:       cmd.Arguments[1],
			UserID:    usr.ID,
		})
	if err != nil {
		return fmt.Errorf("unable to create feed: %v", err)
	}
	fmt.Printf(
		"Feed successfully created\nID: %v\nCreatedAt: %v\nName: %v\nUrl: %v\nUserID: %v",
		feed.ID, feed.CreatedAt, feed.Name, feed.Url, feed.UserID)
	return nil
}

func HandlerAggregate(s *State, cmd Command) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("the agg handler takes no arguments")
	}

	feed, err := s.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("fetching error: %v", err)
	}

	fmt.Println(feed)
	return nil
}

func HandlerGetFeeds(s *State, cmd Command) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("the feeds handler takes no arguments")
	}

	feeds, err := s.dbq.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to retrieve feeds: %v", err)
	}

	for _, feed := range feeds {
		usrName, err := s.dbq.MatchUser(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("unable to match user to feed: %v", err)
		}
		fmt.Printf("Feed: %v\n	-URL: %v\n	-User Name: %v\n", feed.Name, feed.Url, usrName)
	}
	return nil
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
