package main

import (
	"fmt"
	"io"
	. "strconv"
)

//go:generate go run ./
type ExaType []map[int]struct {
	Err error `tag`
}

type Alias = io.ReadWriter

func (et ExaType) String() string {
	return fmt.Sprintf("Hello, %s!", Itoa(65))
}

type Smoo int

const (
	A Smoo = iota
	B
	C
	D
	E
)

func hello() {
	fmt.Printf("hi!\n")
}
