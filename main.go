package main

import (
	"fmt"
	"os"

	"github.com/etuhoha/gator/internal/config"
)

type state struct {
	config *config.Config
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

	conf, err := config.Read()
	if err != nil {
		fmt.Printf("read err: %v\n", err)
	}
	state := state{&conf}

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
		return fmt.Errorf("missing name for login")
	}

	user := cmd.args[0]
	err := s.config.SetUser(user)
	if err != nil {
		return err
	}

	fmt.Printf("user set to '%v'\n", s.config.User)
	return nil
}
