package units

var (
	// Metric
	Gram      = metricMass("gram", "g")
	Yottagram = Yotta(Gram)
	Zettagram = Zetta(Gram)
	Exagram   = Exa(Gram)
	Petagram  = Peta(Gram)
	Teragram  = Tera(Gram)
	Gigagram  = Giga(Gram)
	Megagram  = Mega(Gram)
	Kilogram  = Kilo(Gram)
	Hectogram = Hecto(Gram)
	Decagram  = Deca(Gram)
	Decigram  = Deci(Gram)
	Centigram = Centi(Gram)
	Milligram = Milli(Gram)
	Microgram = Micro(Gram)
	Nanogram  = Nano(Gram)
	Picogram  = Pico(Gram)
	Femtogram = Femto(Gram)
	Attogram  = Atto(Gram)
	Zeptogram = Zepto(Gram)
	Yoctogram = Yocto(Gram)

	// Imperial
	Grain  = imperialMass("grain", "gr")
	Drachm = imperialMass("drachm", "dr")
	Ounce  = imperialMass("ounce", "oz")
	Pound  = imperialMass("pound", "lbs")
	Stone  = imperialMass("stone", "st")
	Ton    = imperialMass("ton", "t")
	Slug   = imperialMass("slug", "slug")
)

func metricMass(name, symbol string) Unit {
	return New(name, symbol, MASS, METRIC)
}

func imperialMass(name, symbol string) Unit {
	return New(name, symbol, MASS, IMPERIAL)
}
