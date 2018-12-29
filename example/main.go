package main

import (
	"errors"
	"fmt"

	"github.com/PieterD/commando"
)

func main() {
	lc := &listCommand{}
	lcfs := commando.NewFlagSet("list")
	lcfs.StringVar(&lc.dir, "dir", "", "Directory to list")
	commando.Register(lcfs, "List all the thingies", lc.run)



	commando.Run()
}

type listCommand struct {
	dir string
}

func (lc *listCommand) run() error {
	if lc.dir == "" {
		return commando.ArgError(errors.New("missing dir flag"))
	}
	fmt.Printf("Listing %s", lc.dir)
	return nil
}