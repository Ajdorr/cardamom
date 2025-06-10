package units

var (
	// Metric
	Liter = metricVolume("liter", "L")

	Yottaliter = Yotta(Liter)
	Zettaliter = Zetta(Liter)
	Exaliter   = Exa(Liter)
	Petaliter  = Peta(Liter)
	Teraliter  = Tera(Liter)
	Gigaliter  = Giga(Liter)
	Megaliter  = Mega(Liter)
	Kiloliter  = Kilo(Liter)
	Hectoliter = Hecto(Liter)
	Decaliter  = Deca(Liter)
	Deciliter  = Deci(Liter)
	Centiliter = Centi(Liter)
	Milliliter = Milli(Liter)
	Microliter = Micro(Liter)
	Nanoliter  = Nano(Liter)
	Picoliter  = Pico(Liter)
	Femtoliter = Femto(Liter)
	Attoliter  = Atto(Liter)
	Zeptoliter = Zepto(Liter)
	Yoctoliter = Yocto(Liter)

	// Imperial
	FluidOunce = imperialVolume("fluid ounce", "fl oz")
	Pint       = imperialVolume("pint", "pt")
	Quart      = imperialVolume("quart", "qt")
	Gallon     = imperialVolume("gallon", "gal")

	// US
	Cup        = usVolume("cup", "cup")
	Teaspoon   = usVolume("teaspoon", "tsp")
	Tablespoon = usVolume("tablespoon", "Tbsp")
	// FluidQuart          = usVolume("fluid quart", "qt")
	// FluidPint           = usVolume("fluid pint", "pt")
	// FluidGallon         = usVolume("fluid gallon", "")
	// CustomaryFluidOunce = usVolume("customary fluid ounce", "")
)

func metricVolume(name, symbol string) Unit {
	return New(name, symbol, VOLUME, METRIC)
}

func imperialVolume(name, symbol string) Unit {
	return New(name, symbol, VOLUME, IMPERIAL)
}

func usVolume(name, symbol string) Unit {
	return New(name, symbol, VOLUME, US)
}
