package units

import "strings"

func Yotta(u Unit) Unit { return prefix(u, "yotta", "Y") }
func Zetta(u Unit) Unit { return prefix(u, "zetta", "Z") }
func Exa(u Unit) Unit   { return prefix(u, "exa", "E") }
func Peta(u Unit) Unit  { return prefix(u, "peta", "P") }
func Tera(u Unit) Unit  { return prefix(u, "tera", "T") }
func Giga(u Unit) Unit  { return prefix(u, "giga", "G") }
func Mega(u Unit) Unit  { return prefix(u, "mega", "M") }
func Kilo(u Unit) Unit  { return prefix(u, "kilo", "k") }
func Hecto(u Unit) Unit { return prefix(u, "hecto", "h") }
func Deca(u Unit) Unit  { return prefix(u, "deca", "da") }
func Deci(u Unit) Unit  { return prefix(u, "deci", "d") }
func Centi(u Unit) Unit { return prefix(u, "centi", "c") }
func Milli(u Unit) Unit { return prefix(u, "milli", "m") }
func Micro(u Unit) Unit { return prefix(u, "micro", "μ") }
func Nano(u Unit) Unit  { return prefix(u, "nano", "n") }
func Pico(u Unit) Unit  { return prefix(u, "pico", "p") }
func Femto(u Unit) Unit { return prefix(u, "femto", "f") }
func Atto(u Unit) Unit  { return prefix(u, "atto", "a") }
func Zepto(u Unit) Unit { return prefix(u, "zepto", "z") }
func Yocto(u Unit) Unit { return prefix(u, "yocto", "y") }

func prefix(u Unit, namePrefix, symbolPrefix string) Unit {
	return New(namePrefix+strings.ToLower(u.name),
		symbolPrefix+u.symbol, u.kind, u.system,
	)
}
