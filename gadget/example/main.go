package main

import (
	"fmt"
	"os"

	"github.com/PieterD/gadget"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	info, err := gadget.Generate()
	if err != nil {
		return errors.Wrapf(err, "failed to fetch go generate info")
	}
	fmt.Printf("%#v\n", info)
	file, err := info.Open()
	if err != nil {
		return errors.Wrapf(err, "failed to open file to generate")
	}
	name, typ, err := info.GetType()
	if err != nil {
		return errors.Wrapf(err, "failed to get type to generate")
	}
	fmt.Printf("name: %s, type: %#v\n", name, typ)
	fmt.Printf("concise: %s\n", typ)

	methods := file.GetMethods(name)

	fmt.Printf("methods: %v\n", methods)

	return nil
}
