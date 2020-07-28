package domain

// A set of sample voyages.
var (
	V100 = New("V100", Schedule{
		[]CarrierMovement{
			{DepartureLocation: CNHKG, ArrivalLocation: JNTKO},
			{DepartureLocation: JNTKO, ArrivalLocation: USNYC},
		},
	})

	V300 = New("V300", Schedule{
		[]CarrierMovement{
			{DepartureLocation: JNTKO, ArrivalLocation: NLRTM},
			{DepartureLocation: NLRTM, ArrivalLocation: DEHAM},
			{DepartureLocation: DEHAM, ArrivalLocation: AUMEL},
			{DepartureLocation: AUMEL, ArrivalLocation: JNTKO},
		},
	})

	V400 = New("V400", Schedule{
		[]CarrierMovement{
			{DepartureLocation: DEHAM, ArrivalLocation: SESTO},
			{DepartureLocation: SESTO, ArrivalLocation: FIHEL},
			{DepartureLocation: FIHEL, ArrivalLocation: DEHAM},
		},
	})
)

// These voyages are hard-coded into the current pathfinder. Make sure
// they exist.
var (
	V0100S = New("0100S", Schedule{[]CarrierMovement{}})
	V0200T = New("0200T", Schedule{[]CarrierMovement{}})
	V0300A = New("0300A", Schedule{[]CarrierMovement{}})
	V0301S = New("0301S", Schedule{[]CarrierMovement{}})
	V0400S = New("0400S", Schedule{[]CarrierMovement{}})
)
