package commando

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type command struct {
	flags *flag.FlagSet
	desc  string
	f     func() error
}

// Output is the Writer where help and errors are printed with regards to command line parsing.
// It is os.Stderr by default.
// It should be set before Run.
var Output = io.Writer(os.Stderr)

type CommandFunc func() error

// ArgError is used to wrap an error returned from CommandFunc,
// in order to signal that the error should be regarded as a command line validation issue.
// If the error returned by CommandFunc is wrapped like this, help will be printed.
func ArgError(err error) error {
	return argError{err: err}
}

type argError struct {
	err error
}

func (ae argError) Error() string {
	return ae.err.Error()
}

func isArgError(err error) bool {
	_, ok := err.(argError)
	return ok
}

var commands = make(map[string]command)

// Register registers the given FlagSet as a command. The command name is fs.Name().
// fs should preferably be made by commando.NewFlagSet,
// but otherwise it should be created with flag.ContinueOnError.
func Register(fs *flag.FlagSet, desc string, f CommandFunc) {
	commands[fs.Name()] = command{
		flags: fs,
		desc:  desc,
		f:     f,
	}
}

// NewFlagSet creates a new FlagSet with ContinueOnError,
// and Usage set directly to PrintDefaults.
func NewFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = fs.PrintDefaults
	return fs
}

func commandNames() ([]string, int) {
	if len(commands) == 0 {
		return nil, 0
	}
	longestName := 0
	names := make([]string, 0, len(commands))
	for name := range commands {
		if len(name) > longestName {
			longestName = len(name)
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names, longestName
}

func usageExit(err error) {
	if err != nil {
		fmt.Fprintf(Output, "%s: %v\n", os.Args[0], err)
	}
	fmt.Fprintf(Output, "usage: %s COMMAND\n", os.Args[0])
	fmt.Fprintf(Output, "commands:\n")
	names, longestName := commandNames()
	for _, name := range names {
		cmd := commands[name]
		fmt.Fprintf(Output, "  %-*s: %s\n", longestName, name, cmd.desc)
	}
	// This may not be correct for -help, but some commands do this.
	os.Exit(2)
}

func cmdUsageExit(fs *flag.FlagSet, err error) {
	if err != nil {
		fmt.Fprintf(Output, "%s: %v\n", os.Args[0], err)
	}
	fmt.Fprintf(Output, "usage: %s %s\n", os.Args[0], fs.Name())
	fs.Usage()
	// This may not be correct for -help, but some commands do this.
	os.Exit(2)
}

// Run will parse the command line flags and run the appropriate command.
// It will exit the application if a problem occurs, and print help and errors to Output if needed.
// The usage line is printed by Run, not Usage; keep this in mind when setting Usage.
func Run() {
	if len(os.Args) < 2 {
		usageExit(errors.New("missing command name"))
	}
	selected := os.Args[1]
	if selected == "-h" || selected == "--help" || selected == "-help" {
		usageExit(nil)
	}
	if strings.HasPrefix(selected, "-") {
		usageExit(errors.New("command name required as first argument"))
	}
	cmd, ok := commands[selected]
	if !ok {
		usageExit(errors.New("unknown command '" + selected + "'"))
	}

	// flags.Parse prints help in case of -help. I'm not a fan of this.
	cmd.flags.SetOutput(ioutil.Discard)
	err := cmd.flags.Parse(os.Args[2:])
	cmd.flags.SetOutput(Output)
	if err == flag.ErrHelp {
		cmdUsageExit(cmd.flags, nil)
	}
	if err != nil {
		cmdUsageExit(cmd.flags, err)
	}
	err = cmd.f()
	if err != nil {
		if isArgError(err) {
			cmdUsageExit(cmd.flags, err)
		}
		fmt.Fprintf(Output, "Error running %s command: %+v", cmd.flags.Name(), err)
		os.Exit(1)
	}
	os.Exit(0)
}
