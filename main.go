package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/etuhoha/gator/internal/config"
	"github.com/etuhoha/gator/internal/database"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	all map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.all[cmd.name]
	if !ok {
		return fmt.Errorf("command not found: %v", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.all[name] = f
}

func main() {
	cmds := commands{make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	conf, err := config.Read()
	if err != nil {
		fmt.Printf("read err: %v\n", err)
	}

	db, err := sql.Open("postgres", conf.DbUrl)
	if err != nil {
		fmt.Printf("can't open DB connection: %v\n", err)
		os.Exit(1)
	}

	state := state{&conf, database.New(db)}

	args := os.Args
	if len(args) < 2 {
		fmt.Printf("not enough arguments\n")
		os.Exit(1)
	}

	cmd := command{args[1], args[2:]}
	err = cmds.run(&state, cmd)
	if err != nil {
		fmt.Printf("command error: %v\n", err)
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing username for login")
	}

	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user '%v' not found: [%v]", name, err)
	}

	err = s.config.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("user set to '%v'\n", s.config.User)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing username for register")
	}

	name := cmd.args[0]

	params := database.CreateUserParams{}
	params.Name = name
	params.ID = uuid.New()
	params.CreatedAt = time.Now()
	params.UpdatedAt = params.CreatedAt

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	err = s.config.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("user created: %v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("users reset\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range users {
		name := u.Name
		if name == s.config.User {
			name += " (current)"
		}
		fmt.Printf("* %v\n", name)
	}

	return nil
}
