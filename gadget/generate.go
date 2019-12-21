package gadget

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Info contains the information about a go generate call.
type Info struct {
	Arch    string
	OS      string
	Package string
	File    string
	Line    int

	file *File
}

// Generate returns the information about the current go generate call.
// It will return an error if this execution was not apparently initiated by go generate.
func Generate() (*Info, error) {
	i := &Info{}
	if err := i.generate(); err != nil {
		return nil, errors.Wrapf(err, "failed to get go:generate info")
	}

	return i, nil
}

func (i *Info) generate() error {
	i.Arch = os.Getenv("GOARCH")
	if i.Arch == "" {
		return errors.Errorf("missing GOARCH environment variable")
	}
	i.OS = os.Getenv("GOOS")
	if i.OS == "" {
		return errors.Errorf("missing GOOS environment variable")
	}
	i.Package = os.Getenv("GOPACKAGE")
	if i.Package == "" {
		return errors.Errorf("missing GOPACKAGE environment variable")
	}
	i.File = os.Getenv("GOFILE")
	if i.File == "" {
		return errors.Errorf("missing GOFILE environment variable")
	}
	lineString := os.Getenv("GOLINE")
	if lineString == "" {
		return errors.Errorf("missing GOLINE environment variable")
	}
	dollar := os.Getenv("DOLLAR")
	if dollar != "$" {
		return errors.Errorf("missing DOLLAR environment variable")
	}
	line, err := strconv.Atoi(lineString)
	if err != nil {
		return errors.Wrapf(err, "invalid GOLINE environment variable '%s': %v", lineString, err)
	}
	i.Line = line
	return nil
}

// Open will open the file the Info refers to.
func (i *Info) Open() (*File, error) {
	file, err := NewFile(i.File, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse file")
	}
	i.file = file
	return file, nil
}

// GetType returns the name and type of the type selected to generate for.
// If we could not find the type we are generating for, GetType returns an error.
func (i *Info) GetType() (string, Type, error) {
	if i.file == nil {
		return "", nil, errors.Errorf("Info.GetType called before Info.Open")
	}
	for _, typ := range i.file.Types {
		if typ.Line == i.Line+1 {
			return typ.Name, typ.Type, nil
		}
	}
	return "", nil, errors.Errorf("unable to find type we are generating for")
}
