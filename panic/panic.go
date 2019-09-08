package panic

type Wrap struct {
	Err error
}

func Panic(err error) {
	if err != nil {
		panic(Wrap{Err: err})
	}
}

func Recover(errp *error) {
	if err := RecoverCheck(recover()); err != nil {
		*errp = err
	}
}

func RecoverCheck(r interface{}) error {
	if r == nil {
		return nil
	}
	if w, ok := r.(Wrap); ok {
		return w.Err
	}
	panic(r)
}
