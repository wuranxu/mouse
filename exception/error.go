package exception

import "fmt"

type Err string

func (e Err) New(s interface{}) Err {
	return Err(fmt.Sprintf("%v: %v", e, s))
}

func (e Err) Error() string {
	return string(e)
}

func Msg(err error, s interface{}) string {
	return fmt.Sprintf("%v: %v", err, s)
}

func Error(err error, s interface{}) error {
	return fmt.Errorf("%v: %v", err, s)
}
