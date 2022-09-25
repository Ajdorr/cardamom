package gorm_ext

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Rational struct {
	// big.Rat
	a, b uint
	// isNegative bool
}

var whitespaceRe *regexp.Regexp

// FIXME Implement
// func (a *Rational) Add(b *Rational) { a.a += b.a }
// func (a *Rational) Subtract(b *Rational) { a.a += b.a }
// func (a *Rational) Multply(b *Rational) { a.a += b.a }
// func (a *Rational) Divide(b *Rational) { a.a += b.a }

func (r *Rational) String() string {
	if r.b == 1 || r.a == 0 {
		return strconv.Itoa(int(r.a))
	} else if r.a >= r.b {
		whole := r.a / r.b
		num := r.a % r.b
		return fmt.Sprintf("%d %d/%d", whole, num, r.b)
	} else {
		return fmt.Sprintf("%d/%d", r.a, r.b)
	}
}

func (r *Rational) Parse(s string) error {
	var err error
	whole := 0

	if arr := whitespaceRe.Split(s, -1); len(arr) == 2 {
		whole, err = strconv.Atoi(arr[0])
		if err != nil {
			return err
		}
		s = arr[1]
	}

	if arr := strings.Split(s, "/"); len(arr) == 2 {
		den, err := strconv.Atoi(arr[1])
		if err != nil {
			return errors.WithStack(err)
		}
		r.b = uint(den)

		num, err := strconv.Atoi(arr[0])
		if err != nil {
			return errors.WithStack(err)
		}
		r.a = uint(num + (whole * den))
	} else {
		num, err := strconv.Atoi(arr[0])
		if err != nil {
			return errors.WithStack(err)
		}
		r.a = uint(num)
		r.b = 1
	}

	return nil
}

// Gorm compatibility
func (r *Rational) Scan(value interface{}) error {
	if str, ok := value.(string); !ok {
		return errors.New(fmt.Sprint("gorm Scan failure, invalid type: ", value))
	} else if err := r.Parse(str); err != nil {
		return fmt.Errorf("gorm Scan failure (%s) -- %w ", str, err)
	}
	return nil
}

func (r Rational) Value() (driver.Value, error) {
	return r.String(), nil
}

// JSON marshalling
func (r *Rational) MarshalJSON() ([]byte, error) {
	str := r.String()
	if strings.Contains(str, "/") {
		return []byte(fmt.Sprintf("\"%s\"", str)), nil
	} else {
		return []byte(str), nil
	}
}

func (r *Rational) UnmarshalJSON(data []byte) error {
	str := strings.ReplaceAll(string(data), "\"", "")
	if err := r.Parse(str); err != nil {
		return err
	}

	return nil
}

func init() {
	var err error
	whitespaceRe, err = regexp.Compile(`\s+`)
	if err != nil {
		panic(err)
	}
}
