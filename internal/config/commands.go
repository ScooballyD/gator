package config

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("addfeed takes two arguments: name, url")
	}

	//

	feed, err := s.dbq.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Arguments[0],
			Url:       cmd.Arguments[1],
			UserID:    user.ID,
		})
	if err != nil {
		return fmt.Errorf("unable to create feed: %v", err)
	}

	_, err = s.dbq.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		return fmt.Errorf("unable to follow feed: %v", err)
	}

	fmt.Printf(
		"Feed successfully created\nID: %v\nCreatedAt: %v\nName: %v\nUrl: %v\nUser: %v",
		feed.ID, feed.CreatedAt, feed.Name, feed.Url, user.Name)
	return nil
}

func HandlerAggregate(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("the agg handler takes 1 argument: time between reqs\nEx: '1m'")
	}

	dur, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unable to parse duration: %v", err)
	}

	ticker := time.NewTicker(dur)
	fmt.Printf("Collecting feed every %v\n", cmd.Arguments[0])
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	lim := 2
	var err error

	if len(cmd.Arguments) > 0 {
		lim, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return fmt.Errorf("unable to process limit %v: %v", cmd.Arguments[0], err)
		}
	}

	posts, err := s.dbq.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(lim),
		})
	if err != nil {
		return fmt.Errorf("unable to get posts: %v", err)
	}

	for _, pst := range posts {
		fmt.Printf("\ntitle: %v\n", pst.Title)
		fmt.Printf("--published at: %v\n", pst.PublishedAt)
		fmt.Printf("--description: %v\n", pst.Description)
		fmt.Printf("--url: %v\n", pst.Url)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("the follow handler takes 1 argumane: url")
	}

	//usr

	feed, err := s.dbq.GetFeed(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unable to find feed: %v", err)
	}

	_, err = s.dbq.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		return fmt.Errorf("unable to follow feed: %v", err)
	}

	fmt.Printf("Feed: %v, followed by %v\n", feed.Name, user.Name)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("the following handler takes no arguments")
	}

	//usr

	feeds, err := s.dbq.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to retrieve followed feeds: %v", err)
	}

	fmt.Printf("Feeds followed by %v:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf(" -%v", feed.FeedName)
	}
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
	fmt.Printf("user has been created: %v\n", usr.Name)
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

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("unfollow handler takes 1 argument: feed URL")
	}

	_, err := s.dbq.Unfollow(context.Background(), database.UnfollowParams{
		UserID: user.ID,
		Url:    cmd.Arguments[0],
	})
	if err != nil {
		return fmt.Errorf("unable to unfollow feed: %v", err)
	}

	return nil
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		usr, err := s.dbq.GetUser(context.Background(), s.point.Current_user_name)
		if err != nil {
			return fmt.Errorf("unable to find current user: %v", err)
		}
		handler(s, cmd, usr)
		return nil
	}
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

func scrapeFeeds(s *State) error {
	feed, err := s.dbq.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to fetch next feed: %v", err)
	}

	feed, err = s.dbq.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("unable to mark next feed: %v", err)
	}

	items, err := s.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("unable to list feed: %v", err)
	}
	if items == nil {
		return fmt.Errorf("items is %v", err)
	}

	for _, itm := range items.Channel.Item {
		t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", itm.PubDate)
		if err != nil {
			return fmt.Errorf("unable to parse date: %v", err)
		}
		_, err = s.dbq.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       itm.Title,
				Url:         itm.Link,
				Description: itm.Description,
				PublishedAt: t,
				FeedID:      feed.ID,
			})
		if err != nil {
			return fmt.Errorf("unable to save post: %v", err)
		}
	}
	return nil
}
