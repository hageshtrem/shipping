package domain

// Sample UN locodes.
var (
	SESTO UNLocode = "SESTO"
	SEGOT UNLocode = "SEGOT"
	AUMEL UNLocode = "AUMEL"
	CNHKG UNLocode = "CNHKG"
	CNSHA UNLocode = "CNSHA"
	CNHGH UNLocode = "CNHGH"
	USNYC UNLocode = "USNYC"
	USCHI UNLocode = "USCHI"
	USDAL UNLocode = "USDAL"
	JNTKO UNLocode = "JNTKO"
	DEHAM UNLocode = "DEHAM"
	NLRTM UNLocode = "NLRTM"
	FIHEL UNLocode = "FIHEL"
)

// Sample locations.
var (
	Stockholm = &Location{SESTO, "Stockholm"}
	Goteborg  = &Location{SEGOT, "Goteborg"}
	Melbourne = &Location{AUMEL, "Melbourne"}
	Hongkong  = &Location{CNHKG, "Hongkong"}
	Shanghai  = &Location{CNSHA, "Shanghai"}
	Hangzhou  = &Location{CNHGH, "Hangzhou"}
	NewYork   = &Location{USNYC, "New York"}
	Chicago   = &Location{USCHI, "Chicago"}
	Dallas    = &Location{USDAL, "Dallas"}
	Tokyo     = &Location{JNTKO, "Tokyo"}
	Hamburg   = &Location{DEHAM, "Hamburg"}
	Rotterdam = &Location{NLRTM, "Rotterdam"}
	Helsinki  = &Location{FIHEL, "Helsinki"}
)
