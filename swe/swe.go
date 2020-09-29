package swe

import (
	"sync"
)

// BaseInterface defines a standardized way for interfacing with the Swiss
// Ephemeris library from Go.
type BaseInterface interface {
	// Version returns the version of the Swiss Ephemeris.
	Version() (string, error)

	// PlanetName returns the name of planet pl.
	PlanetName(pl Planet) (string, error)

	// Calc computes the position and optionally the speed of planet pl at Julian
	// Date (in Ephemeris Time) et with calculation flags fl.
	Calc(et float64, pl Planet, fl *CalcFlags) (xx []float64, cfl int, err error)
	// CalcUT computes the position and optionally the speed of planet pl at
	// Julian Date (in Universal Time) ut with calculation flags fl. Within the C
	// library swe_deltat is called to convert Universal Time to Ephemeris Time.
	CalcUT(ut float64, pl Planet, fl *CalcFlags) (xx []float64, cfl int, err error)

	// NodAps computes the positions of planetary nodes and apsides (perihelia,
	// aphelia, second focal points of the orbital ellipses) for planet pl at
	// Julian Date (in Ephemeris Time) et with calculation flags fl using method
	// m.
	NodAps(et float64, pl Planet, fl *CalcFlags, m NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error)
	// NodApsUT computes the positions of planetary nodes and apsides (perihelia,
	// aphelia, second focal points of the orbital ellipses) for planet pl at
	// Julian Date (in Ephemeris Time) et with calculation flags fl using method
	// m. Within the C library swe_deltat is called to convert Universal Time to
	// Ephemeris Time.
	NodApsUT(ut float64, pl Planet, fl *CalcFlags, m NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error)

	// GetAyanamsaEx returns the ayanamsa for Julian Date (in Ephemeris Time) et.
	// It is equal to GetAyanamsa but uses the ΔT consistent with the ephemeris
	// passed in fl.Flags.
	GetAyanamsaEx(et float64, fl *AyanamsaExFlags) (float64, error)
	// GetAyanamsaExUT returns the ayanamsa for Julian Date (in Universal Time) ut.
	// It is equal to GetAyanamsaUT but uses the ΔT consistent with the ephemeris
	// passed in fl.Flags.
	GetAyanamsaExUT(ut float64, fl *AyanamsaExFlags) (float64, error)
	// GetAyanamsaName returns the name of sidmode.
	GetAyanamsaName(ayan Ayanamsa) (string, error)

	// JulDay returns the corresponding Julian Date for the given date.
	// Calendar type ct is used to clearify the year y, Julian or Gregorian.
	JulDay(y, m, d int, h float64, ct CalType) (float64, error)
	// RevJul returns the corresponding calendar date for the given Julian Date.
	// Calendar type ct is used to clearify the year y, Julian or Gregorian.
	RevJul(jd float64, ct CalType) (y, m, d int, h float64, err error)
	// UTCToJD returns the corresponding Julian Date in Ephemeris and Universal
	// Time for the given date and accounts for leap seconds in the conversion.
	UTCToJD(y, m, d, h, i int, s float64, fl *DateConvertFlags) (et, ut float64, err error)
	// JdETToUTC returns the corresponding calendar date for the given Julian
	// Date in Ephemeris Time and accounts for leap seconds in the conversion.
	JdETToUTC(et float64, fl *DateConvertFlags) (y, m, d, h, i int, s float64, err error)
	// JdETToUTC returns the corresponding calendar date for the given Julian
	// Date in Universal Time and accounts for leap seconds in the conversion.
	JdUT1ToUTC(ut1 float64, fl *DateConvertFlags) (y, m, d, h, i int, s float64, err error)

	// HousesEx returns the house cusps and related positions for the given
	// geographic location using the given house system and the provided flags
	// (reference frame). The return values may contain data in case of an error.
	// Geolat and geolon are in degrees.
	HousesEx(ut float64, fl *HousesExFlags, geolat, geolon float64, hsys HSys) ([]float64, []float64, error)
	// HousesArmc returns the house cusps and related positions for the given
	// geographic location using the given house system, ecliptic obliquity and
	// ARMC (also known as RAMC). The return values may contain data in case of
	// an error. ARMC, geolat, geolon and eps are in degrees.
	HousesARMC(armc, geolat, eps float64, hsys HSys) ([]float64, []float64, error)
	// HousePos returns the house position for the ecliptic longitude and
	// latitude of a planet for a given ARMC (also known as RAMC) and geocentric
	// latitude using the given house system. ARMC, geolat, eps, pllng and pllat
	// are in degrees.
	// Before calling HousePos either Houses, HousesEx or HousesARMC should be
	// called first.
	HousePos(armc, geolat, eps float64, hsys HSys, pllng, pllat float64) (float64, error)
	// HouseName returns the name of the house system.
	HouseName(hsys HSys) (string, error)

	// DeltaTEx returns the ΔT for the Julian Date jd.
	// 通过星历表精确的计算ET - UT
	DeltaTEx(jd float64, eph Ephemeris) (float64, error)
	// 粗略的计算ET - UT
	DeltaT(jd float64) (float64)

	// TimeEqu returns the difference between local apparent and local mean time
	// in days for the given Julian Date (in Universal Time).
	TimeEqu(jd float64, fl *TimeEquFlags) (float64, error)
	// LMTToLAT returns the local apparent time for the given Julian Date (in
	// Local Mean Time) and the geographic longitude.
	LMTToLAT(jdLMT, geolon float64, fl *TimeEquFlags) (float64, error)
	// LATToLMT returns the local mean time for the given Julian Date (in Local
	// Apparent Time) and the geographic longitude.
	LATToLMT(jdLAT, geolon float64, fl *TimeEquFlags) (float64, error)

	// SidTime0 returns the sidereal time for Julian Date jd, ecliptic obliquity
	// eps and nutation nut at the Greenwich medidian, measured in hours.
	SidTime0(ut, eps, nut float64, fl *SidTimeFlags) (float64, error)
	// SidTime returns the sidereal time for Julian Date jd at the Greenwich
	// medidian, measured in hours.
	SidTime(ut float64, fl *SidTimeFlags) (float64, error)
}

// SweInterface extends the main library interface by exposing C library
// life-cycle methods.
type SweInterface interface {
	// The following methods will always return nil as error:
	//  Version
	//  PlanetName
	//  GetAyanamsaName
	//  JulDay
	//  RevJul
	//  JdETToUTC
	//  JdUT1ToUTC
	//  HouseName
	//  SidTime
	//  SidTime0
	BaseInterface

	// SetPath opens the ephemeris and sets the data path.
	SetPath(path string)

	// Close closes the Swiss Ephemeris library.
	// The ephemeris can be reopened by calling SetPath.
	Close()

	// used for locking and prevent other interface implementations
	acquire()
	release()
}

var once sync.Once
var sweInstance SweInterface

// NewSwe returns an object that calls the Swiss Ephemeris C library.
// The returned object is safe for concurrent use.
func NewSwe() SweInterface {
	once.Do(func() {
		checkLibrary()
		sweInstance = &swe{locker: new(sync.Mutex)}
	})

	return sweInstance
}

// Open initializes the Swiss Ephemeris C library with DefaultPath as
// ephemeris path. The returned object is safe for concurrent use.
func Open() SweInterface {
	return OpenWithPath(DefaultPath)
}

// OpenWithPath initializes the Swiss Ephemeris C library and calls
// swe_set_ephe_path with ephePath as argument afterwards.
// The returned object is safe for concurrent use.
func OpenWithPath(ephePath string) SweInterface {
	swe := NewSwe()
	swe.SetPath(ephePath)
	return swe
}

// swe interfaces between swego.BaseInterface and the library functions.
// It protect stateful library functions with a mutex. When the swe is
// exclusively locked, the mutex is temporary replaced by a no-op lock.
type swe struct {
	locker sync.Locker
}

func (s *swe) acquire() { s.locker.Lock() }
func (s *swe) release() { s.locker.Unlock() }

// acquire locks the swe for exclusive library access.
// release unlocks the swe from exclusive library access.

var _ SweInterface = (*swe)(nil) // assert interface

func (s *swe) Version() (string, error) {
	return Version, nil
}

func (s *swe) SetPath(ephePath string) {
	s.acquire()
	setEphePath(ephePath)
	s.release()
}

func (s *swe) Close() {
	s.acquire()
	closeEphemeris()
	s.release()
}

const resetDeltaT = -1e-10

func setDeltaT(dt *float64) {
	var f float64
	if dt == nil {
		f = resetDeltaT
	} else {
		f = *dt
	}

	setDeltaTUserDef(f)
}

func setCalcFlagsState(cf *CalcFlags) int32 {
	if cf == nil {
		setDeltaT(nil)
		return 0
	}

	if (cf.Flags & flgTopo) == flgTopo {
		var lng, lat, alt float64

		if cf.TopoLoc != nil {
			lng = cf.TopoLoc.Long
			lat = cf.TopoLoc.Lat
			alt = cf.TopoLoc.Alt
		}

		setTopo(lng, lat, alt)
	}

	if (cf.Flags & flgSidereal) == flgSidereal {
		var mode Ayanamsa
		var t0, ayanT0 float64

		if cf.SidMode != nil {
			mode = cf.SidMode.Mode
			t0 = cf.SidMode.T0
			ayanT0 = cf.SidMode.AyanT0
		}

		setSidMode(mode, t0, ayanT0)
	}

	if cf.JPLFile != "" {
		cf.JPLFile = FnameDft
	}

	setJPLFile(cf.JPLFile)
	setDeltaT(cf.DeltaT)
	return cf.Flags
}

func (s *swe) PlanetName(pl Planet) (string, error) {
	s.acquire()
	name := planetName(pl)
	s.release()
	return name, nil
}

func (s *swe) Calc(et float64, pl Planet, cf *CalcFlags) ([]float64, int, error) {
	s.acquire()
	flags := setCalcFlagsState(cf)
	xx, cfl, err := calc(et, pl, flags)
	s.release()
	return xx, cfl, err
}

func (s *swe) CalcUT(ut float64, pl Planet, cf *CalcFlags) ([]float64, int, error) {
	s.acquire()
	flags := setCalcFlagsState(cf)
	xx, cfl, err := calcUT(ut, pl, flags)
	s.release()
	return xx, cfl, err
}

func (s *swe) NodAps(et float64, pl Planet, cf *CalcFlags, m NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	s.acquire()
	flags := setCalcFlagsState(cf)
	nasc, ndsc, peri, aphe, err = nodAps(et, pl, flags, m)
	s.release()
	return
}

func (s *swe) NodApsUT(ut float64, pl Planet, cf *CalcFlags, m NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	s.acquire()
	flags := setCalcFlagsState(cf)
	nasc, ndsc, peri, aphe, err = nodApsUT(ut, pl, flags, m)
	s.release()
	return
}

func (s *swe) GetAyanamsaEx(et float64, af *AyanamsaExFlags) (float64, error) {
	s.acquire()
	setSidMode(af.SidMode.Mode, af.SidMode.T0, af.SidMode.AyanT0)
	f, err := getAyanamsaEx(et, af.Flags)
	s.release()
	return f, err
}

func (s *swe) GetAyanamsaExUT(ut float64, af *AyanamsaExFlags) (float64, error) {
	s.acquire()
	setSidMode(af.SidMode.Mode, af.SidMode.T0, af.SidMode.AyanT0)
	f, err := getAyanamsaExUT(ut, af.Flags)
	s.release()
	return f, err
}

func (s *swe) GetAyanamsaName(ayan Ayanamsa) (string, error) {
	s.acquire()
	name := getAyanamsaName(ayan)
	s.release()
	return name, nil
}

func (s *swe) JulDay(year, month, day int, hour float64, ct CalType) (float64, error) {
	jd := julDay(year, month, day, hour, int(ct))
	return jd, nil
}

func (s *swe) RevJul(jd float64, ct CalType) (year, month, day int, hour float64, err error) {
	year, month, day, hour = revJul(jd, int(ct))
	return year, month, day, hour, nil
}

func (s *swe) UTCToJD(year, month, day, hour, minute int, second float64, dcf *DateConvertFlags) (et, ut float64, err error) {
	s.acquire()
	setDeltaT(dcf.DeltaT)
	et, ut, err = utcToJD(year, month, day, hour, minute, second, int(dcf.Calendar))
	s.release()
	return
}

func (s *swe) JdETToUTC(et float64, dcf *DateConvertFlags) (year, month, day, hour, minute int, second float64, err error) {
	s.acquire()
	setDeltaT(dcf.DeltaT)
	year, month, day, hour, minute, second = jdETToUTC(et, int(dcf.Calendar))
	s.release()
	return year, month, day, hour, minute, second, nil
}

func (s *swe) JdUT1ToUTC(ut1 float64, dcf *DateConvertFlags) (year, month, day, hour, minute int, second float64, err error) {
	s.acquire()
	setDeltaT(dcf.DeltaT)
	year, month, day, hour, minute, second = jdUT1ToUTC(ut1, int(dcf.Calendar))
	s.release()
	return year, month, day, hour, minute, second, nil
}

func (s *swe) HousesEx(ut float64, hf *HousesExFlags, geolat, geolon float64, hsys HSys) ([]float64, []float64, error) {
	s.acquire()
	var flags int32
	if hf != nil {
		flags = hf.Flags
		if (flags & flgSidereal) == flgSidereal {
			setSidMode(hf.SidMode.Mode, hf.SidMode.T0, hf.SidMode.AyanT0)
		}

		setDeltaT(hf.DeltaT)
	} else {
		setDeltaT(nil)
	}

	cusps, ascmc, err := housesEx(ut, flags, geolat, geolon, hsys)
	s.release()
	return cusps, ascmc, err
}

func (s *swe) HousesARMC(armc, geolat, eps float64, hsys HSys) ([]float64, []float64, error) {
	s.acquire()
	cusps, ascmc, err := housesARMC(armc, geolat, eps, hsys)
	s.release()
	return cusps, ascmc, err
}

func (s *swe) HousePos(armc, geolat, eps float64, hsys HSys, pllng, pllat float64) (float64, error) {
	s.acquire()
	pos, err := housePos(armc, geolat, eps, hsys, pllng, pllat)
	s.release()
	return pos, err
}

func (s *swe) HouseName(hsys HSys) (string, error) {
	s.acquire()
	name := houseName(hsys)
	s.release()
	return name, nil
}

func (s *swe) DeltaTEx(jd float64, eph Ephemeris) (float64, error) {
	s.acquire()
	dt, err := deltaTEx(jd, int32(eph))
	s.release()
	return dt, err
}

func (s *swe) DeltaT(jd float64) float64 {
	return deltaT(jd)
}

func setTimeEquDeltaT(tf *TimeEquFlags) {
	if tf == nil {
		setDeltaT(nil)
	} else {
		setDeltaT(tf.DeltaT)
	}
}

func (s *swe) TimeEqu(jd float64, tf *TimeEquFlags) (float64, error) {
	s.acquire()
	setTimeEquDeltaT(tf)
	f, err := timeEqu(jd)
	s.release()
	return f, err
}

func (s *swe) LMTToLAT(lmt, geolon float64, tf *TimeEquFlags) (float64, error) {
	s.acquire()
	setTimeEquDeltaT(tf)
	lat, err := lmtToLAT(lmt, geolon)
	s.release()
	return lat, err
}

func (s *swe) LATToLMT(lat, geolon float64, tf *TimeEquFlags) (float64, error) {
	s.acquire()
	setTimeEquDeltaT(tf)
	lmt, err := latToLMT(lat, geolon)
	s.release()
	return lmt, err
}

func setSidTimeDeltaT(fl *SidTimeFlags) {
	if fl == nil {
		setDeltaT(nil)
	} else {
		setDeltaT(fl.DeltaT)
	}
}

func (s *swe) SidTime0(ut, eps, nut float64, stf *SidTimeFlags) (float64, error) {
	s.acquire()
	setSidTimeDeltaT(stf)
	f := sidTime0(ut, eps, nut)
	s.release()
	return f, nil
}

func (s *swe) SidTime(ut float64, stf *SidTimeFlags) (float64, error) {
	s.acquire()
	setSidTimeDeltaT(stf)
	f := sidTime(ut)
	s.release()
	return f, nil
}
