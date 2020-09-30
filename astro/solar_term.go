package astro

import (
	"go-swe/swe"
)

var SolarTermsString = [...]string{
	"春分", "清明", "谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋", "处暑", "白露", "秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至", "小寒", "大寒", "立春", "雨水", "惊蛰",
}

/**
 * 指定太阳黄经, 反推出当时的儒略日
 * firstDayJdET 当年首日的 jdET
 * firstDayPlanet 当年首日的天体属性
 * eclipticLongitude 目标黄经
 */
func calcJulianDayBySolarEclipticLongitude(
	astro *Astronomy,
	firstDayJdET *EphemerisTime,
	firstDayPlanet *PlanetProperties,
	eclipticLongitude float64) (_jd JulianDay, err error) {

	_jd = firstDayJdET.JdUT
	_lastPlanet := firstDayPlanet
	// 将1月1日正午设置为相对0rad，这种算法可以简化当前黄经已经越过360°之后从0开始，而lastLong还在360°内，导致时间差异360多天的情况
	_firstDayLong := firstDayPlanet.Ecliptic.Longitude
	// 当前黄经相对首日黄经的差值
	_lastLongDelta := 0.
	// 目标黄经相对首日黄经的差值 并换算到360°内
	_eclipticLongitudeDelta := RadiansMod360(eclipticLongitude - _firstDayLong)

	for {
		// 已经到了目标黄经, 返回结果
		if FloatEqual(_eclipticLongitudeDelta, _lastLongDelta) {
			return _jd, nil
		}

		// 黄经速度 单位是天
		delta := (_eclipticLongitudeDelta - _lastLongDelta) / _lastPlanet.SpeedInLongitude

		_jd = _jd.Add(delta)
		_lastJdET := NewEphemerisTime(_jd)

		_lastPlanet, err = astro.PlanetProperties(_lastPlanet.PlanetId, _lastJdET)
		if err != nil {
			return 0, err
		}

		// 当前黄经相对首日黄经的差值 并换算到360°内
		_lastLongDelta = RadiansMod360(_lastPlanet.Ecliptic.Longitude - _firstDayLong)
	}

}

/**
 * 该年的24节气的儒略日
 * year 年
 */
func (astro *Astronomy) SolarTerms(year int) (solarTerms *[24]JulianDay, err error) {
	jd := DateToJulianDay(year, 1, 1, 0, 0, 0)
	jdET := NewEphemerisTime(jd)

	// 计算当年1月1日的黄经
	planet, err := astro.PlanetProperties(swe.Sun, jdET)

	if err != nil {
		return nil, err
	}

	solarTerms = &[24]JulianDay{}

	firstLong := ToDegrees(planet.Ecliptic.Longitude)
	nextLong := NextMultiples(firstLong, 15) // 下一个整点的黄经

	for i := 0; i < 24; i++ {
		long := nextLong + 15*float64(i)
		jd, err := calcJulianDayBySolarEclipticLongitude(astro, jdET, planet, ToRadians(long))
		if err != nil {
			return nil, err
		}

		solarTerms[int(long/15)%24] = jd
	}

	return

}
