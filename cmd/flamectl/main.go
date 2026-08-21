package main

import (
	"fmt"
	"gestureflame/cli"
	"gestureflame/clock"
	"gestureflame/service"
	"gestureflame/store"
	"os"
)

func main() {
	path := "flame.db"
	db, err := store.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	app := service.New(db, clock.New())
	if err := cli.Run(app, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
