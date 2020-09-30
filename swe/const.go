package swe

import "fmt"

// CalType represents the calendar type used in julian date conversions.
type CalType int

// Calendar types defined in swephexp.h.
const (
	Julian    CalType = 0
	Gregorian CalType = 1
)

// Planet is the type of planet constants.
type Planet int

// Planet, fictional body and asteroid constants defined in swephexp.h.
const (
	Sun          Planet = 0
	Moon         Planet = 1
	Mercury      Planet = 2
	Venus        Planet = 3
	Mars         Planet = 4
	Jupiter      Planet = 5
	Saturn       Planet = 6
	Uranus       Planet = 7
	Neptune      Planet = 8
	Pluto        Planet = 9
	MeanNode     Planet = 10
	TrueNode     Planet = 11
	MeanApogee   Planet = 12
	OscuApogee   Planet = 13
	Earth        Planet = 14
	Chiron       Planet = 15
	Pholus       Planet = 16
	Ceres        Planet = 17
	Pallas       Planet = 18
	Juno         Planet = 19
	Vesta        Planet = 20
	InterApogee  Planet = 21
	InterPerigee Planet = 22

	Varuna Planet = AstOffset + 20000
	Nessus Planet = AstOffset + 7066

	Cupido   Planet = 40
	Hades    Planet = 41
	Zeus     Planet = 42
	Kronos   Planet = 43
	Apollon  Planet = 44
	Admetos  Planet = 45
	Vulkanus Planet = 46
	Poseidon Planet = 47

	Isis             Planet = 48
	Nibiru           Planet = 49
	Harrington       Planet = 50
	NeptuneLeverrier Planet = 51
	NeptuneAdams     Planet = 52
	PlutoLowell      Planet = 53
	PlutoPickering   Planet = 54
	Vulcan           Planet = 55
	WhiteMoon        Planet = 56
	Proserpina       Planet = 57
	Waldemath        Planet = 58

	EclNut Planet = -1

	AstOffset = 10000
)

const (
	_Planet_name_0 = "SunMoonMercuryVenusMarsJupiterSaturnUranusNeptunePlutoMeanNodeTrueNodeMeanApogeeOscuApogeeEarthChironPholusCeresPallasJunoVestaInterApogeeInterPerigee"
	_Planet_name_1 = "CupidoHadesZeusKronosApollonAdmetosVulkanusPoseidonIsisNibiruHarringtonNeptuneLeverrierNeptuneAdamsPlutoLowellPlutoPickeringVulcanWhiteMoonProserpinaWaldemath"
	_Planet_name_2 = "Nessus"
	_Planet_name_3 = "Varuna"
)

var (
	_Planet_index_0 = [...]uint8{0, 3, 7, 14, 19, 23, 30, 36, 42, 49, 54, 62, 70, 80, 90, 95, 101, 107, 112, 118, 122, 127, 138, 150}
	_Planet_index_1 = [...]uint8{0, 6, 11, 15, 21, 28, 35, 43, 51, 55, 61, 71, 87, 99, 110, 124, 130, 139, 149, 158}
	_Planet_index_2 = [...]uint8{0, 6}
	_Planet_index_3 = [...]uint8{0, 6}
)

func (i Planet) String() string {
	switch {
	case 0 <= i && i <= 22:
		return _Planet_name_0[_Planet_index_0[i]:_Planet_index_0[i+1]]
	case 40 <= i && i <= 58:
		i -= 40
		return _Planet_name_1[_Planet_index_1[i]:_Planet_index_1[i+1]]
	case i == 17066:
		return _Planet_name_2
	case i == 30000:
		return _Planet_name_3
	default:
		return fmt.Sprintf("Planet(%d)", i)
	}
}

//go:generate stringer -type=Planet

// Indexes of related house positions defined in swephexp.h.
const (
	Asc    = 0
	MC     = 1
	ARMC   = 2
	Vertex = 3
	EquAsc = 4 // "equatorial ascendant"
	CoAsc1 = 5 // "co-ascendant" (W. Koch)
	CoAsc2 = 6 // "co-ascendant" (M. Munkasey)
	PolAsc = 7 // "polar ascendant" (M. Munkasey)
)

// Ephemeris represents an ephemeris implemented in the C library.
type Ephemeris int32

// Ephemerides that are implemented in the C library.
const (
	JPL        Ephemeris = FlagEphJPL
	Swiss      Ephemeris = FlagEphSwiss
	Moshier    Ephemeris = FlagEphMoshier
	DefaultEph Ephemeris = FlagEphDefault
)

// Calculation flags defined in swephexp.h.
const (
	FlagEphJPL       = 1 << 0
	FlagEphSwiss     = 1 << 1
	FlagEphMoshier   = 1 << 2
	FlagHelio        = 1 << 3
	FlagTruePos      = 1 << 4
	FlagJ2000        = 1 << 5
	FlagNoNut        = 1 << 6
	FlagSpeed        = 1 << 8
	FlagNoGDefl      = 1 << 9
	FlagNoAbber      = 1 << 10
	FlagAstrometric  = FlagNoAbber | FlagNoGDefl
	FlagEquatorial   = 1 << 11
	FlagXYZ          = 1 << 12
	FlagRadians      = 1 << 13
	FlagBary         = 1 << 14
	FlagTopo         = 1 << 15
	FlagSidereal     = 1 << 16
	FlagICRS         = 1 << 17
	FlagJPLHor       = 1 << 18
	FlagJPLHorApprox = 1 << 19
	FlagEphDefault   = FlagEphSwiss
)

// Ayanamsa is the type of sidereal mode constants.
type Ayanamsa int32

// Sidereal modes (ayanamsas) implemented in the C library.
const (
	SidmFaganBradley       Ayanamsa = 0
	SidmLahiri             Ayanamsa = 1
	SidmDeluce             Ayanamsa = 2
	SidmRaman              Ayanamsa = 3
	SidmUshashashi         Ayanamsa = 4
	SidmKrishnamurti       Ayanamsa = 5
	SidmDjwhalKhul         Ayanamsa = 6
	SidmYukteshwar         Ayanamsa = 7
	SidmJNBhasin           Ayanamsa = 8
	SidmBabylKruger1       Ayanamsa = 9
	SidmBabylKruger2       Ayanamsa = 10
	SidmBabylKruger3       Ayanamsa = 11
	SidmBabylHuber         Ayanamsa = 12
	SidmBabylEtaPiscium    Ayanamsa = 13
	SidmAldebaran15Tau     Ayanamsa = 14
	SidmHipparchos         Ayanamsa = 15
	SidmSassanian          Ayanamsa = 16
	SidmGalCent0Sag        Ayanamsa = 17
	SidmJ2000              Ayanamsa = 18
	SidmJ1900              Ayanamsa = 19
	SidmB1950              Ayanamsa = 20
	SidmSuryasiddhanta     Ayanamsa = 21
	SidmSuryasiddhantaMSun Ayanamsa = 22
	SidmAryabhata          Ayanamsa = 23
	SidmAryabhataMSun      Ayanamsa = 24
	SidmSSRevati           Ayanamsa = 25
	SidmSSCitra            Ayanamsa = 26
	SidmTrueCitra          Ayanamsa = 27
	SidmTrueRevati         Ayanamsa = 28
	SidmTruePushya         Ayanamsa = 29
	SidmGalCentGilBrand    Ayanamsa = 30
	SidmGalAlignMardyks    Ayanamsa = 31
	SidmGalEquIAU1958      Ayanamsa = 32
	SidmGalEquTrue         Ayanamsa = 33
	SidmGalEquMula         Ayanamsa = 34
	SidmGalTrueMula        Ayanamsa = 35
	SidmGalCentMulaWilhelm Ayanamsa = 36
	SidmAryabhata522       Ayanamsa = 37
	SidmBabylBritton       Ayanamsa = 38
	SidmUser               Ayanamsa = 255
)

// Options that augment a sidereal mode (ayanamsa).
const (
	SidbitEclT0    Ayanamsa = 256
	SidbitSSYPlane Ayanamsa = 512
	SidbitUserUT   Ayanamsa = 1024
)

// NodApsMethod is the type of Nodbit constants.
type NodApsMethod int32

// Nodes and apsides calculation bits defined in swephexp.h.
const (
	NodbitMean     NodApsMethod = 1
	NodbitOscu     NodApsMethod = 2
	NodbitOscuBary NodApsMethod = 4
	NodbitFoPoint  NodApsMethod = 256
)

// File name of JPL data files defined in swephexp.h.
const (
	FnameDE200 = "de200.eph"
	FnameDE406 = "de406.eph"
	FnameDE431 = "de431.eph"
	FnameDft   = FnameDE431
	FnameDft2  = FnameDE406
)

// HSys represents house system identifiers used in the C library.
type HSys byte

// House systems implemented in the C library.
const (
	Alcabitius           HSys = 'B'
	Campanus             HSys = 'C'
	EqualMC              HSys = 'D' // Equal houses, where cusp 10 = MC
	Equal                HSys = 'E' // also 'A'
	CarterPoliEquatorial HSys = 'F'
	Gauquelin            HSys = 'G'
	Azimuthal            HSys = 'H' // a.k.a Horizontal
	Sunshine             HSys = 'I' // Makransky, solution Treindl
	SunshineAlt          HSys = 'i' // Makransky, solution Makransky
	Koch                 HSys = 'K'
	PullenSD             HSys = 'L'
	Morinus              HSys = 'M'
	EqualAsc             HSys = 'N' // Equal houses, where cusp 1 = 0° Aries
	Porphyrius           HSys = 'O' // a.k.a Porphyry
	Placidus             HSys = 'P'
	PullenSR             HSys = 'Q'
	Regiomontanus        HSys = 'R'
	Sripati              HSys = 'S'
	PolichPage           HSys = 'T' // a.k.a. Topocentric
	KrusinskiPisaGoelzer HSys = 'U'
	VehlowEqual          HSys = 'V' // Equal Vehlow (Asc in middle of house 1)
	WholeSign            HSys = 'W'
	AxialRotation        HSys = 'X' // a.k.a. Meridian
	APCHouses            HSys = 'Y'
)

// SidMode represents library state changed by swe_set_sid_mode.
type SidMode struct {
	Mode   Ayanamsa
	T0     float64
	AyanT0 float64
}

// GeoLoc represents a geographic location.
type GeoLoc struct {
	Long float64
	Lat  float64
	Alt  float64
}

// CalcFlags represents the library state of swe_calc and swe_calc_ut.
type CalcFlags struct {
	Flags   int32
	TopoLoc *GeoLoc  // Arguments to swe_set_topo
	SidMode *SidMode // Arguments to swe_set_sid_mode
	JPLFile string   // Argument to swe_set_jpl_file
	DeltaT  *float64 // Argument to swe_set_delta_t_userdef, nil resets it.
}

// Copy returns a copy of the calculation flags fl.
func (cf *CalcFlags) Copy() *CalcFlags {
	copy := new(CalcFlags)
	*copy = *cf
	return copy
}

// SetEphemeris sets the ephemeris flag in fl.
func (cf *CalcFlags) SetEphemeris(eph Ephemeris) {
	cf.Flags |= int32(eph)
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (cf *CalcFlags) SetDeltaT(f float64) {
	cf.DeltaT = &f
}

// AyanamsaExFlags represents the library state of swe_get_ayanamsa_ex and
// swe_get_ayanamsa_ex_ut.
type AyanamsaExFlags struct {
	Flags   int32
	SidMode *SidMode // Argument to swe_set_sid_mode
	DeltaT  *float64 // Argument to swe_set_delta_t_userdef, nil resets it.
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (af *AyanamsaExFlags) SetDeltaT(f float64) {
	af.DeltaT = &f
}

// DateConvertFlags represents the library state of swe_utc_to_jd,
// swe_jdet_to_utc and swe_jdut1_to_utc.
type DateConvertFlags struct {
	Calendar CalType // clearifies the input year, Julian or Gregorian
	DeltaT   *float64
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (df *DateConvertFlags) SetDeltaT(f float64) {
	df.DeltaT = &f
}

// HousesExFlags represents library state of swe_houses_ex in a stateless way.
type HousesExFlags struct {
	Flags   int32
	SidMode *SidMode // Argument to swe_set_sid_mode
	DeltaT  *float64 // Argument to swe_set_delta_t_userdef, nil resets it.
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (hf *HousesExFlags) SetDeltaT(f float64) {
	hf.DeltaT = &f
}

// NewHSys validates the input and returns a HSys value if valid.
func NewHSys(c byte) (hsys HSys, ok bool) {
	if c == 'i' {
		return HSys(c), true
	}

	// It's trivial to convert lower case to upper case in ASCII.
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}

	switch c {
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y':
		return HSys(c), true
	default:
		return 0, false
	}
}

// TimeEquFlags represents the library state of swe_time_equ, swe_lmt_to_lat
// and swe_lat_to_lmt.
type TimeEquFlags struct {
	DeltaT *float64
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (tf *TimeEquFlags) SetDeltaT(f float64) {
	tf.DeltaT = &f
}

// SidTimeFlags represents the library state of swe_sidtime0 and swe_sidtime.
type SidTimeFlags struct {
	DeltaT *float64
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (sf *SidTimeFlags) SetDeltaT(f float64) {
	sf.DeltaT = &f
}
