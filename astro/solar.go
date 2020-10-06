package astro

import (
	"go-swe/swe"
)

var SolarTermsString = [...]string{
	"春分", "清明", "谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋", "处暑", "白露", "秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至", "小寒", "大寒", "立春", "雨水", "惊蛰",
}

/**
 * 指定太阳黄经, 反推出当时的时间
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
		if FloatEqual(_eclipticLongitudeDelta, _lastLongDelta, 9) {
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
 * 指定太阳黄经(数组), 反推出当时的时间(数组)
 * year 年
 * eclipticLongitudes 黄经弧度数组
 * 不论eclipticLongitudes有多少项，都只计算year年的数据，这个和 LunarSolarEclipticLongitudeDeltaToTimes 有区别
 */
func (astro *Astronomy) SolarEclipticLongitudesToTimes(year int, eclipticLongitudes []float64) ([]JulianDay, error) {
	jd := DateToJulianDay(year, 1, 1, 0, 0, 0)
	jdET := NewEphemerisTime(jd)

	// 计算当年1月1日的黄经
	planet, err := astro.PlanetProperties(swe.Sun, jdET)
	if err != nil {
		return nil, err
	}

	times := make([]JulianDay, len(eclipticLongitudes))

	for i, angle := range eclipticLongitudes {
		jd, err := calcJulianDayBySolarEclipticLongitude(astro, jdET, planet, angle)
		if err != nil {
			return nil, err
		}
		times[i] = jd
	}

	return times, nil
}

/**
 * 该年的24节气的时间
 * year 年
 */
func (astro *Astronomy) SolarTerms(year int) ([]*JulianDayWithIndex, error) {
	jd := DateToJulianDay(year, 1, 1, 0, 0, 0)
	jdET := NewEphemerisTime(jd)
	solarTerms := make([]*JulianDayWithIndex, 24)

	// 15° 每个节气
	degreePerSolarTerm := 15
	// 1年24个节气
	solarTermCount := 24

	// 计算当年1月1日的黄经
	planet, err := astro.PlanetProperties(swe.Sun, jdET)
	if err != nil {
		return nil, err
	}

	// 第一个有效的角度
	firstValidLongDegree := NextMultiples(ToDegrees(planet.Ecliptic.Longitude), float64(degreePerSolarTerm))

	// 24节气分别对应的黄经, 春分0开始
	eclipticLongitudes := make([]float64, solarTermCount)
	for i := 0; i < 24; i++ {
		long := 15*float64(i) + firstValidLongDegree
		solarTerms[i] = NewJulianDayWithIndex(0, int(long)/degreePerSolarTerm%solarTermCount)
		eclipticLongitudes[i] = ToRadians(long)
	}

	times, err := astro.SolarEclipticLongitudesToTimes(year, eclipticLongitudes)

	if err != nil {
		return nil, err
	}

	for i, jd := range times {
		// i 转化为节气的索引，春分是0
		solarTerms[i].JdUT = jd
	}

	return solarTerms, nil

}
