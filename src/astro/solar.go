package astro

import (
	"fmt"
	"go-swe/src/swe"
	"math"
)

var SolarTermsString = [...]string{
	"春分", "清明", "谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋", "处暑", "白露", "秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至", "小寒", "大寒", "立春", "雨水", "惊蛰",
}

/**
 * 指定太阳黄经, 反推出当时的时间
 * startJdET 时间起点jdET
 * startPlanet 该起点的天体属性
 * eclipticLongitude 目标黄经
 */
func calcJulianDayBySolarEclipticLongitude(
	astro *Astronomy,
	startJdET *EphemerisTime,
	startPlanet *PlanetProperties,
	eclipticLongitude float64) (JulianDay, *PlanetProperties, int, error) {

	_jd := startJdET.JdUT
	_lastPlanet := startPlanet
	// 将1月1日正午设置为相对0rad，这种算法可以简化当前黄经已经越过360°之后从0开始，而lastLong还在360°内，导致时间差异360多天的情况
	_startLong := startPlanet.Ecliptic.Longitude
	// 当前黄经相对首日黄经的差值
	_lastLongDelta := 0.
	// 目标黄经相对首日黄经的差值 并换算到360°内
	_eclipticLongitudeDelta := RadiansMod360(eclipticLongitude - _startLong)

	calcCount := 0
	for {
		// 已经到了目标黄经, 返回结果
		if FloatEqual(_eclipticLongitudeDelta, _lastLongDelta, 9) {
			return _jd, _lastPlanet, calcCount, nil
		}

		// 黄经速度 单位是天
		dayDelta := (_eclipticLongitudeDelta - _lastLongDelta) / _lastPlanet.SpeedInLongitude

		_jd = _jd.Add(dayDelta)

		var err error
		_lastPlanet, err = astro.PlanetProperties(_lastPlanet.PlanetId, NewEphemerisTime(_jd))
		if err != nil {
			return 0, nil, calcCount, fmt.Errorf("calcJulianDayBySolarEclipticLongitude: %w", err)
		}

		// 当前黄经相对首日黄经的差值 并换算到360°内
		_lastLongDelta = RadiansMod360(_lastPlanet.Ecliptic.Longitude - _startLong)
		calcCount++
	}

}

/**
 * 指定太阳黄经, 反推出当时的时间
 * startJdUT 时间起始
 * eclipticLongitude 黄经弧度
 */
func (astro *Astronomy) SolarEclipticLongitudesToTime(startJdUT JulianDay, eclipticLongitude float64) (JulianDay, int, error) {
	jdET := NewEphemerisTime(startJdUT)

	// startJd的黄经
	planet, err := astro.PlanetProperties(swe.Sun, jdET)
	if err != nil {
		return 0, 0, fmt.Errorf("SolarEclipticLongitudesToTime Calc 1: %w", err)
	}

	jd, planet, calcCount, err := calcJulianDayBySolarEclipticLongitude(astro, jdET, planet, eclipticLongitude)
	if err != nil {
		return 0, 0, fmt.Errorf("SolarEclipticLongitudesToTime Calc 2: %w", err)
	}
	return jd, calcCount, nil
}

/**
 * 指定太阳黄经(数组), 反推出当时的时间(数组)
 * startJd 时间起始
 * eclipticLongitudes 黄经弧度数组
 * 比如(这里以度表示，实际应该是弧度)：{270, 0, 270}，表示startJdUT之后的{冬至, 春分，下一个冬至}
 */
func (astro *Astronomy) SolarEclipticLongitudesToTimes(startJdUT JulianDay, eclipticLongitudes []float64) ([]JulianDay, error) {

	times := make([]JulianDay, len(eclipticLongitudes))
	var jd = startJdUT

	planet, err := astro.PlanetProperties(swe.Sun, NewEphemerisTime(jd))
	if err != nil {
		return nil, fmt.Errorf("SolarEclipticLongitudesToTimes Calc 1: %w", err)
	}

	lastAngle := math.Inf(-1)
	for i, angle := range eclipticLongitudes {
		// 如果上一个弧度差和现在的相等，为了避免SolarEclipticLongitudesToTime不往后面推进，jd累加半年、planet重新计算
		if FloatEqual(angle, lastAngle, 9) {
			jd += MeanSolarDays * .8
			planet, err = astro.PlanetProperties(swe.Sun, NewEphemerisTime(jd))
			if err != nil {
				return nil, fmt.Errorf("SolarEclipticLongitudesToTimes Calc 3: %w", err)
			}
		}
		jd, planet, _, err = calcJulianDayBySolarEclipticLongitude(astro, NewEphemerisTime(jd), planet, angle)
		if err != nil {
			return nil, fmt.Errorf("SolarEclipticLongitudesToTimes Calc 2: %w", err)
		}
		times[i] = jd
		lastAngle = angle
	}

	return times, nil
}

/**
 * 2时间之间的所有节气
 * startJdUT 起始时间
 * endJdUT 结束时间
 */
func (astro *Astronomy) SolarTermsRange(startJdUT, endJdUT JulianDay) ([]*JulianDayExtra, error) {
	// 15° 每个节气
	const degreePerSolarTerm = 15.
	// 1年24个节气
	const solarTermCount = 24.

	solarTerms := make([]*JulianDayExtra, 0, int(float64(endJdUT-startJdUT)/MeanSolarDays*solarTermCount))

	// 计算startJdUT的黄经
	planet, err := astro.PlanetProperties(swe.Sun, NewEphemerisTime(startJdUT))
	if err != nil {
		return nil, fmt.Errorf("SolarTerms Calc 1: %w", err)
	}

	// 第一个有效的角度
	long := NextMultiples(ToDegrees(planet.Ecliptic.Longitude), degreePerSolarTerm)

	var jd = startJdUT
	for ; jd <= endJdUT; long += degreePerSolarTerm {
		jd, _, err = astro.SolarEclipticLongitudesToTime(jd, ToRadians(long))
		if err != nil {
			return nil, fmt.Errorf("SolarTerms Calc 2: %w", err)
		}
		index := int(long/degreePerSolarTerm) % int(solarTermCount)
		solarTerms = append(solarTerms, NewJulianDayExtra(jd, index))
	}

	// 去掉最后一个不符合的项目
	if solarTerms[len(solarTerms)-1].JdUT > endJdUT {
		return solarTerms[0 : len(solarTerms)-1], nil
	}

	return solarTerms, nil
}

/**
 * 该年的24节气的时间
 * year 年
 */
func (astro *Astronomy) SolarTerms(year int) ([]*JulianDayExtra, error) {
	jd := DateToJulianDay(year, 1, 1, 0, 0, 0)
	return astro.SolarTermsRange(jd, jd.AddYears(1))
}
