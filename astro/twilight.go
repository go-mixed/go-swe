package astro

import (
	"go-swe/swe"
	"math"
)

// 天体的 升/降 角度
type TwilightAngle struct {
	// 升天, 降天
	RiseSet,
	// 上中天，下中天
	Culmination, LowerCulmination float64
}

/**
 * 太阳 晨/暮 角度
 *
 *  地平线  -----  0°
 *           |
 *  升/落   -----  -0.833°
 *           |
 *  民用      |
 *           |
 *         -----  -6°
 *           |
 *  航海      |
 *           |
 *         ----- -12°
 *           |
 *  天文      |
 *           |
 *         ----- -18°
 *           |
 *   夜晚     |
 *           |
 *           |
 */
type SunTwilightAngle struct {
	TwilightAngle
	// 各个专业领域的晨/暮角度
	Astronomical, Nautical, Civil float64
}

// 天体的 升/降 时间
type TwilightTimes struct {
	// 升天, 降天
	Rise, Set JulianDay
	// 上中天，下中天
	Culmination, LowerCulmination JulianDay
}

// 天体的 晨/暮 时间
type DawnDuskTime struct {
	// 晨
	Dawn JulianDay
	// 暮
	Dusk JulianDay
}

// 太阳的 升/降/晨/暮 时间
type SunTwilightTimes struct {
	TwilightTimes

	Astronomical DawnDuskTime
	Nautical     DawnDuskTime
	Civil        DawnDuskTime
}

func NewSunTwilightAngle() *SunTwilightAngle {
	return &SunTwilightAngle{
		TwilightAngle: TwilightAngle{
			RiseSet:          ToRadians(-50. / 60.), // 地平线以下 50′
			Culmination:      ToRadians(90.),
			LowerCulmination: ToRadians(-90.),
		},
		Astronomical: ToRadians(-18.),
		Nautical:     ToRadians(-12.),
		Civil:        ToRadians(-6.),
	}
}

func NewTwilightAngle() *TwilightAngle {
	return &TwilightAngle{
		RiseSet:          0,
		Culmination:      ToRadians(90.),
		LowerCulmination: ToRadians(-90.),
	}
}

/**
 * 日照时长
 * 日落 - 升日
 */
func (stt *SunTwilightTimes) Daylight() float64 {
	return float64(stt.Set - stt.Rise)
}

/**
 * 夜晚时长，因为跨午夜，所以取今日00:00 ~ 23:59的时段
 * (0:00 ~ 天文晨光) + (天文暮光 ~ 23:59)
 */
func (stt *SunTwilightTimes) Night() float64 {
	return float64(
		stt.Astronomical.Dawn - stt.Astronomical.Dawn.Midnight() +
			(stt.Astronomical.Dusk.AddDays(1).Midnight() - stt.Astronomical.Dusk),
	)
}

type haToDeltaCallback func(hourAngle, baseHourAngle HourAngle, microStep bool) float64

var riseCallback = func(hourAngle HourAngle, baseHourAngle HourAngle, microStep bool) float64 {
	if microStep {
		return RadiansMod180(float64(-hourAngle-baseHourAngle)) / Radian360
	}
	return float64(-hourAngle-baseHourAngle) / Radian360
}
var setCallback = func(hourAngle HourAngle, baseHourAngle HourAngle, microStep bool) float64 {
	if microStep {
		return RadiansMod180(float64(hourAngle-baseHourAngle)) / Radian360
	}
	return float64(hourAngle-baseHourAngle) / Radian360
}
var culminationCallback = func(hourAngle HourAngle, baseHourAngle HourAngle, microStep bool) float64 {
	return float64(0-baseHourAngle) / Radian360
}
var lowerCulminationCallback = func(hourAngle HourAngle, baseHourAngle HourAngle, microStep bool) float64 {
	if microStep {
		return RadiansMod180(Radian180-float64(baseHourAngle)) / Radian360
	}
	return Radian180 - float64(baseHourAngle)/Radian360
}
var moonRiseCallback = func(hourAngle HourAngle, baseHourAngle HourAngle, microStep bool) float64 {
	if microStep {
		return RadiansMod180(float64(-hourAngle-baseHourAngle)) / (Radian360 * 0.966)
	}
	return float64(-hourAngle-baseHourAngle) / (Radian360 * 0.966)
}
var moonSetCallback = func(hourAngle HourAngle, baseHourAngle HourAngle, microStep bool) float64 {
	if microStep {
		return RadiansMod180(float64(hourAngle-baseHourAngle)) / (Radian360 * 0.966)
	}
	return float64(hourAngle-baseHourAngle) / (Radian360 * 0.966)
}
var moonCulminationCallback = func(hourAngle HourAngle, baseHourAngle HourAngle, microStep bool) float64 {
	if microStep {
		return RadiansMod180(0-float64(baseHourAngle)) / (Radian360 * 0.966)
	}
	return float64(0-baseHourAngle) / (Radian360 * 0.966)
}
var moonLowerCulminationCallback = func(hourAngle HourAngle, baseHourAngle HourAngle, microStep bool) float64 {
	if microStep {
		return RadiansMod180(Radian180-float64(baseHourAngle)) / (Radian360 * 0.966)
	}
	return Radian180 - float64(baseHourAngle)/(Radian360*0.966)
}

/**
 * 指定高度角，反推出当时的儒略日
 * lastPlanet 初步计算的天体属性
 * jdET 初步计算的时间
 * geo 观察者地理位置
 * angle 指定需要计算的高度角
 * withRevise 是否修正
 * callback ha的计算函数
 */
func calcJulianDayByAltitude(
	astro *Astronomy,
	lastPlanet *PlanetProperties,
	jdET *EphemerisTime,
	geo *GeographicCoordinates,
	angle float64,
	withRevise bool,
	callback haToDeltaCallback,
) (_jdUT JulianDay) {

	_jdUT = jdET.JdUT
	_jdET := jdET
	var _ha HourAngle
	var _err error

	// 计算第一次
	if math.Abs(angle) < Radian90 {
		_ha = AltitudeToHourAngle(lastPlanet.Equatorial.Declination, geo.Latitude, angle)
	}

	_jdUT = _jdUT.Add(callback(_ha, lastPlanet.HourAngle, false))

	// 多修正几次
	for i := 0; i < 3; i++ {
		_jdET = NewEphemerisTime(_jdUT)
		lastPlanet, _err = astro.PlanetPropertiesWithObserver(lastPlanet.PlanetId, _jdET, geo, withRevise)
		if _err != nil {
			return
		}

		if math.Abs(angle) < Radian90 {
			_ha = AltitudeToHourAngle(lastPlanet.Equatorial.Declination, geo.Latitude, angle)
		}

		_jdUT = _jdUT.Add(callback(_ha, lastPlanet.HourAngle, true))
	}

	return
}

/**
 * 天体升/降/中天时间
 * jdUT UT的儒略日
 * geo 观察者地理位置
 * planetId 天体ID
 * angle 升/降/中天的角度
 * withRevise 是否修正一些日光差，或者黄道章动
 */
func (astro *Astronomy) PlanetTwilight(jdUT JulianDay, geo *GeographicCoordinates, planetId swe.Planet, angle *TwilightAngle, withRevise bool) *TwilightTimes {
	jdET := NewEphemerisTime(jdUT)

	if angle == nil {
		angle = NewTwilightAngle()
	}
	times := &TwilightTimes{}

	// 天体属性
	planet, err := astro.PlanetPropertiesWithObserver(planetId, jdET, geo, withRevise)
	if err != nil {
		return nil
	}

	var _reviseHourAngle = func(angle float64, callback haToDeltaCallback) JulianDay {
		return calcJulianDayByAltitude(astro, planet, jdET, geo, angle, withRevise, callback)
	}

	times.Rise = _reviseHourAngle(angle.RiseSet, riseCallback)
	times.Set = _reviseHourAngle(angle.RiseSet, setCallback)
	times.Culmination = _reviseHourAngle(angle.Culmination, culminationCallback)
	times.LowerCulmination = _reviseHourAngle(angle.LowerCulmination, lowerCulminationCallback)

	return times
}

/**
 * 太阳升/降/中天/晨/暮时间
 * jdUT UT的儒略日，传入本地12点的JdUT(比如东八区是当天4点: TimeToJulianDay(2020-09-30 04:00:00))
 * geo 观察者地理位置
 * withRevise: 是否修正一些日光差，或者黄道章动
 */
func (astro *Astronomy) SunTwilight(jdUT JulianDay, geo *GeographicCoordinates, withRevise bool) *SunTwilightTimes {
	// 查找最靠近当日中午的日上中天, mod2的第1参数为本地时角近似值
	noonJdUT := jdUT.Add(-Mod2(float64(jdUT.AsJD2000())+geo.Longitude/Radian360, 1))
	jdET := NewEphemerisTime(noonJdUT)

	angle := NewSunTwilightAngle()
	sunTimes := &SunTwilightTimes{}

	// 天体属性
	planet, err := astro.PlanetPropertiesWithObserver(swe.Sun, jdET, geo, withRevise)
	if err != nil {
		return nil
	}

	var _reviseHourAngle = func(angle float64, callback haToDeltaCallback) JulianDay {
		return calcJulianDayByAltitude(astro, planet, jdET, geo, angle, withRevise, callback)
	}

	sunTimes.Rise = _reviseHourAngle(angle.RiseSet, riseCallback)
	sunTimes.Set = _reviseHourAngle(angle.RiseSet, setCallback)
	sunTimes.Civil.Dawn = _reviseHourAngle(angle.Civil, riseCallback)
	sunTimes.Civil.Dusk = _reviseHourAngle(angle.Civil, setCallback)
	sunTimes.Nautical.Dawn = _reviseHourAngle(angle.Nautical, riseCallback)
	sunTimes.Nautical.Dusk = _reviseHourAngle(angle.Nautical, setCallback)
	sunTimes.Astronomical.Dawn = _reviseHourAngle(angle.Astronomical, riseCallback)
	sunTimes.Astronomical.Dusk = _reviseHourAngle(angle.Astronomical, setCallback)
	sunTimes.Culmination = _reviseHourAngle(angle.Culmination, culminationCallback)
	sunTimes.LowerCulmination = _reviseHourAngle(angle.LowerCulmination, lowerCulminationCallback)

	return sunTimes
}

/**
 * 月亮升/降/中天时间
 * jdUT UT的儒略日，传入本地12点的JdUT(比如东八区是当天4点: TimeToJulianDay(2020-09-30 04:00:00))
 * geo 观察者地理位置
 * withRevise: 是否修正一些日光差，或者黄道章动
 */
func (astro *Astronomy) MoonTwilight(jdUT JulianDay, geo *GeographicCoordinates, withRevise bool) *TwilightTimes {
	deltaT := astro.DeltaT(jdUT)

	// 查找最靠近当日中午的月上中天, mod2的第1参数为本地时角近似值
	moonJdUT := jdUT.Add(-Mod2(0.1726222+0.966136808032357*float64(jdUT.AsJD2000())-0.0366*deltaT+geo.Longitude/Radian360, 1))
	jdET := NewEphemerisTime(moonJdUT)

	angle := NewTwilightAngle()
	moonTimes := &TwilightTimes{}

	// 天体属性
	planet, err := astro.PlanetPropertiesWithObserver(swe.Moon, jdET, geo, withRevise)
	if err != nil {
		return nil
	}

	angle.RiseSet = 0.7275*EquatorialRadius/planet.DistanceAsKilometer() - 34*60/DegreeSecondsPerRadian

	var _reviseHourAngle = func(angle float64, callback haToDeltaCallback) JulianDay {
		return calcJulianDayByAltitude(astro, planet, jdET, geo, angle, withRevise, callback)
	}

	moonTimes.Rise = _reviseHourAngle(angle.RiseSet, moonRiseCallback)
	moonTimes.Set = _reviseHourAngle(angle.RiseSet, moonSetCallback)
	moonTimes.Culmination = _reviseHourAngle(angle.Culmination, moonCulminationCallback)
	moonTimes.LowerCulmination = _reviseHourAngle(angle.Culmination, moonLowerCulminationCallback)

	return moonTimes
}
