package log_ext

import (
	"fmt"

	"github.com/pkg/errors"
)

func Errorf(format string, a ...any) error {
	return errors.WithStack(fmt.Errorf(format, a...))
}

func ReturnBoth(err string) (string, error) {
	return err, fmt.Errorf(err)
}

func ReturnBothErr(err error) (string, error) {
	return err.Error(), err
}
