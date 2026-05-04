package main

import (
	"fmt"

	"github.com/etuhoha/gator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Printf("read err: %v\n", err)
	}
	err = conf.SetUser("hoha")
	if err != nil {
		fmt.Printf("write err: %v\n", err)
	}
	conf, err = config.Read()
	if err != nil {
		fmt.Printf("read err: %v\n", err)
	}
	fmt.Printf("%v", conf)
}
