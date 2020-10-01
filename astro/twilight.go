package astro

import (
	"go-swe/swe"
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
type SunTwilightAngles struct {
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
type DawnDuskTimes struct {
	// 晨
	Dawn JulianDay
	// 暮
	Dusk JulianDay
}

// 太阳的 升/降/晨/暮 时间
type SunTwilightTimes struct {
	TwilightTimes

	Astronomical DawnDuskTimes
	Nautical     DawnDuskTimes
	Civil        DawnDuskTimes
}

func NewSunTwilightAngles() *SunTwilightAngles {
	return &SunTwilightAngles{
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

/**
 * 指定高度角，反推出当时的儒略日，因为1天内任意高度角有2个相同的，所以会分别返回东、西两个方位的时间
 * 如果高度角=90°或-90°是中天，上中天返回相同的2个时间，下中天返回当日和次日的时间
 * lastPlanet 初步计算的天体属性
 * jdET 初步计算的时间
 * geo 观察者地理位置
 * altitude 指定需要计算的高度角，每个天体在不同纬度观察都有不同的高度角，并且并不是每天能达到90°
 * withRevise 是否修正
 */
func calcJulianDayByAltitude(
	astro *Astronomy,
	lastPlanet *PlanetProperties,
	jdET *EphemerisTime,
	geo *GeographicCoordinates,
	altitude float64,
	withRevise bool,
) (*[2]JulianDay, error) {

	var _ha HourAngle
	var err error
	var _delta float64
	var times = &[2]JulianDay{}

	// 计算第一次
	_ha = AltitudeToHourAngle(lastPlanet.Equatorial.Declination, geo.Latitude, altitude)

	fixValue := 1.
	switch lastPlanet.PlanetId {
	case swe.Moon:
		fixValue = 0.966
	}

	// 东边
	_delta = float64(-_ha-lastPlanet.HourAngle) / (Radian360 * fixValue)
	times[0] = jdET.JdUT.Add(_delta)

	// 西边
	_delta = float64(_ha-lastPlanet.HourAngle) / (Radian360 * fixValue)
	times[1] = jdET.JdUT.Add(_delta)

	// 计算第二次, 修正东
	lastPlanet, err = astro.PlanetPropertiesWithObserver(lastPlanet.PlanetId, NewEphemerisTime(times[0]), geo, withRevise)
	if err != nil {
		return nil, err
	}
	_ha = AltitudeToHourAngle(lastPlanet.Equatorial.Declination, geo.Latitude, altitude)
	_delta = RadiansMod180(float64(-_ha-lastPlanet.HourAngle)) / (Radian360 * fixValue)
	times[0] = times[0].Add(_delta)

	// 计算第二次, 修正西
	lastPlanet, err = astro.PlanetPropertiesWithObserver(lastPlanet.PlanetId, NewEphemerisTime(times[1]), geo, withRevise)
	if err != nil {
		return nil, err
	}
	_ha = AltitudeToHourAngle(lastPlanet.Equatorial.Declination, geo.Latitude, altitude)
	_delta = RadiansMod180(float64(_ha-lastPlanet.HourAngle)) / (Radian360 * fixValue)
	times[1] = times[1].Add(_delta)

	return times, nil
}

/**
 * 根据天体高度角，求时间，一个角度有2个时间
 * jdUT UT的儒略日，传入该天体中天的近似儒略日（需要转为UT），比如太阳是本地12点的JdUT(比如东八区是当天4点: TimeToJulianDay(2020-09-30 04:00:00))
 * geo 观察者地理位置
 * planetId 天体ID
 * altitude 带求值的高度角
 * withRevise 是否修正一些日光差，或者黄道章动
 */
func (astro *Astronomy) AltitudeToTimes(jdUT JulianDay, geo *GeographicCoordinates, planetId swe.Planet, altitude float64, withRevise bool) (*[2]JulianDay, error) {
	jdET := NewEphemerisTime(jdUT)

	// 天体属性
	planet, err := astro.PlanetPropertiesWithObserver(planetId, jdET, geo, withRevise)
	if err != nil {
		return nil, err
	}

	return calcJulianDayByAltitude(astro, planet, jdET, geo, altitude, withRevise)
}

/**
 * 太阳升/降/中天/晨/暮时间
 * jdUT UT的儒略日，传入本地12点的JdUT(比如东八区是当天4点: TimeToJulianDay(2020-09-30 04:00:00))
 * 关于中天，根据纬度的不同，太阳只有在春分（赤道）、秋分（赤道）、夏至（北回归线）、冬至（南回归线）才能达到90°，超过当日最大的高度角，一律按照90°计算
 * geo 观察者地理位置
 * withRevise: 是否修正一些日光差，或者黄道章动
 */
func (astro *Astronomy) SunTwilight(jdUT JulianDay, geo *GeographicCoordinates, withRevise bool) (*SunTwilightTimes, error) {
	// 查找最靠近当日中午的日上中天, mod2的第1参数为本地时角近似值
	noonJdUT := jdUT.Add(-Mod2(float64(jdUT.AsJD2000())+geo.Longitude/Radian360, 1))
	jdET := NewEphemerisTime(noonJdUT)

	angle := NewSunTwilightAngles()
	sunTimes := &SunTwilightTimes{}

	var times *[2]JulianDay
	var err error

	// 天体属性
	planet, err := astro.PlanetPropertiesWithObserver(swe.Sun, jdET, geo, withRevise)
	if err != nil {
		return nil, err
	}

	// 上中天
	times, err = calcJulianDayByAltitude(astro, planet, jdET, geo, angle.Culmination, withRevise)
	if err != nil {
		return nil, err
	}
	sunTimes.Culmination = times[0]

	// 下中天
	times, err = calcJulianDayByAltitude(astro, planet, jdET, geo, -angle.Culmination, withRevise)
	if err != nil {
		return nil, err
	}
	sunTimes.LowerCulmination = times[0]

	// 升/降
	times, err = calcJulianDayByAltitude(astro, planet, jdET, geo, angle.RiseSet, withRevise)
	if err != nil {
		return nil, err
	}

	sunTimes.Rise = times[0]
	sunTimes.Set = times[1]

	// 民用晨/暮
	times, err = calcJulianDayByAltitude(astro, planet, jdET, geo, angle.Civil, withRevise)
	if err != nil {
		return nil, err
	}

	sunTimes.Civil.Dawn = times[0]
	sunTimes.Civil.Dusk = times[1]

	// 航海晨/暮
	times, err = calcJulianDayByAltitude(astro, planet, jdET, geo, angle.Nautical, withRevise)
	if err != nil {
		return nil, err
	}

	sunTimes.Nautical.Dawn = times[0]
	sunTimes.Nautical.Dusk = times[1]

	// 天文晨/暮
	times, err = calcJulianDayByAltitude(astro, planet, jdET, geo, angle.Astronomical, withRevise)
	if err != nil {
		return nil, err
	}

	sunTimes.Astronomical.Dawn = times[0]
	sunTimes.Astronomical.Dusk = times[1]

	return sunTimes, nil
}

/**
 * 月亮升/降/中天时间
 * jdUT UT的儒略日，传入本地12点的JdUT(比如东八区是当天4点: TimeToJulianDay(2020-09-30 04:00:00))
 * geo 观察者地理位置
 * withRevise: 是否修正一些日光差，或者黄道章动
 */

func (astro *Astronomy) MoonTwilight(jdUT JulianDay, geo *GeographicCoordinates, withRevise bool) (*TwilightTimes, error) {
	deltaT := astro.DeltaT(jdUT)

	// 查找最靠近当日中午的月上中天, mod2的第1参数为本地时角近似值
	moonJdUT := jdUT.Add(-Mod2(0.1726222+0.966136808032357*float64(jdUT.AsJD2000())-0.0366*deltaT+geo.Longitude/Radian360, 1))
	jdET := NewEphemerisTime(moonJdUT)

	angle := NewTwilightAngle()
	moonTimes := &TwilightTimes{}

	var times *[2]JulianDay
	var err error

	// 天体属性
	planet, err := astro.PlanetPropertiesWithObserver(swe.Moon, jdET, geo, withRevise)
	if err != nil {
		return nil, err
	}

	angle.RiseSet = 0.7275*EquatorialRadius/planet.DistanceAsKilometer() - 34*60/DegreeSecondsPerRadian

	// 上中天
	times, err = calcJulianDayByAltitude(astro, planet, jdET, geo, angle.Culmination, withRevise)
	if err != nil {
		return nil, err
	}
	moonTimes.Culmination = times[0]

	// 下中天
	times, err = calcJulianDayByAltitude(astro, planet, jdET, geo, -angle.Culmination, withRevise)
	if err != nil {
		return nil, err
	}
	moonTimes.LowerCulmination = times[0]

	// 升/降
	times, err = calcJulianDayByAltitude(astro, planet, jdET, geo, angle.RiseSet, withRevise)
	if err != nil {
		return nil, err
	}

	moonTimes.Rise = times[0]
	moonTimes.Set = times[1]

	return moonTimes, nil
}
