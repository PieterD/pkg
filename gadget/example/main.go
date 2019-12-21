package main

import (
	"fmt"
	"os"

	"github.com/PieterD/pkg/gadget"
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
		return fmt.Errorf("failed to fetch go generate info: %w", err)
	}
	fmt.Printf("%#v\n", info)
	file, err := info.Open()
	if err != nil {
		return fmt.Errorf("failed to open file to generate: %w", err)
	}
	name, typ, err := info.GetType()
	if err != nil {
		return fmt.Errorf("failed to get type to generate: %w", err)
	}
	fmt.Printf("name: %s, type: %#v\n", name, typ)
	fmt.Printf("concise: %s\n", typ)

	methods := file.GetMethods(name)

	fmt.Printf("methods: %v\n", methods)

	return nil
}
