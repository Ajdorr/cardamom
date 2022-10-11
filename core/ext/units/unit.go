package units

import (
	"cardamom/core/ext/log_ext"
	"database/sql/driver"
	"errors"
	"fmt"
)

var unitMap = make(map[string]Unit)

type System string

const (
	METRIC   = "Metric"
	IMPERIAL = "Imperial"
	US       = "Imperial"
)

type Kind string

const (
	VOLUME = "Volume"
	MASS   = "Mass"
)

type Unit struct {
	name   string
	symbol string
	kind   Kind
	system System
}

func New(name, symbol string, kind Kind, system System) Unit {
	if _, ok := unitMap[symbol]; ok {
		panic(errors.New("duplicate symbol name: " + name))
	}

	u := Unit{
		name:   name,
		symbol: symbol,
		kind:   kind,
		system: system,
	}
	unitMap[symbol] = u

	return u
}

func (u *Unit) Name() string   { return u.name }
func (u *Unit) Symbol() string { return u.symbol }
func (u *Unit) Kind() Kind     { return u.kind }
func (u *Unit) System() System { return u.system }

// Gorm compatibility
func (u *Unit) Scan(value interface{}) error {
	if name, ok := value.(string); !ok {
		return errors.New(fmt.Sprint("gorm Scan failure, invalid Unit type: ", value))
	} else if unit, ok := unitMap[name]; !ok {
		return errors.New(fmt.Sprint("gorm Scan failure, missing Unit name: ", name))
	} else {
		*u = unit
	}

	return nil
}

func (u Unit) Value() (driver.Value, error) {
	return u.symbol, nil
}

// JSON compatibility
func (u *Unit) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, u.symbol)), nil
}

func (u *Unit) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) <= 2 {
		return log_ext.Errorf("json cannot be empty")
	}

	if unit, ok := unitMap[string(data[1:len(data)-1])]; !ok {
		return errors.New(fmt.Sprint("invalid Unit type: ", string(data)))
	} else {
		*u = unit
		return nil
	}
}
