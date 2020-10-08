package astro

import (
	"fmt"
	"go-swe/swe"
	"math"
)

var LunarPhaseStrings = [...]string{"朔", "上弦", "望", "下弦"}
var LunarMonthStrings = [...]string{"正月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "腊月", "闰"}
var LunarDayStrings = [...]string{
	"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九",
	"初十", "十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九",
	"二十", "廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九",
	"三十", "卅一",
}

func GetLunarMonthString(index int, leap bool) string {
	if leap {
		return LunarMonthStrings[12] + LunarMonthStrings[index%12]
	}
	return LunarMonthStrings[index%12]
}

/**
 * LunarSolarEclipticLongitudeDelta 月亮和太阳的黄经之差
 */
func (astro *Astronomy) LunarSolarEclipticLongitudeDelta(jdET *EphemerisTime) (float64, float64, error) {
	sun, err := astro.PlanetProperties(swe.Sun, jdET)
	if err != nil {
		return 0, 0, fmt.Errorf("LunarSolarEclipticLongitudeDelta of Sun: %w", err)
	}

	moon, err := astro.PlanetProperties(swe.Moon, jdET)
	if err != nil {
		return 0, 0, fmt.Errorf("LunarSolarEclipticLongitudeDelta of Moon: %w", err)
	}

	return moon.Ecliptic.Longitude - sun.Ecliptic.Longitude, moon.SpeedInLongitude - sun.SpeedInLongitude, nil
}

/**
 * LunarSolarEclipticLongitudeDeltaToTime 从jdET开始，以月球和太阳的黄经之差求出具体时间，比如 朔（0），上弦（90°），望（180°），下弦（270°）
 * jdET 以此时间开始
 * eclipticLongitudeDelta 弧差
 * precision 是否更加精确的计算，精确计算会比较耗时
 */
func (astro *Astronomy) LunarSolarEclipticLongitudeDeltaToTime(jdET *EphemerisTime, eclipticLongitudeDelta float64) (JulianDay, int, error) {
	eclipticLongitudeDelta = math.Abs(eclipticLongitudeDelta)

	const step = 27. * 0.125 // 4/1000,每次递进一小步
	lastJdUT := jdET.JdUT

	lastDelta, lastSpeedDelta, err := astro.LunarSolarEclipticLongitudeDelta(jdET)
	if err != nil {
		return 0, 0, fmt.Errorf("LunarSolarEclipticLongitudeDeltaToTime Calc 1: %w", err)
	}

	calcCount := 0

	if true {
		var firstDelta = lastDelta
		var eclipticLongitudeDeltaDelta = RadiansMod360(eclipticLongitudeDelta - firstDelta)
		var lastDeltaDelta = 0.
		for {
			if FloatEqual(lastDeltaDelta, eclipticLongitudeDeltaDelta, 9) {
				return lastJdUT, calcCount, nil
			}

			// 黄经速度 单位是天
			dayDelta := (eclipticLongitudeDeltaDelta - lastDeltaDelta) / lastSpeedDelta

			lastJdUT = lastJdUT.Add(dayDelta)

			lastDelta, lastSpeedDelta, err = astro.LunarSolarEclipticLongitudeDelta(NewEphemerisTime(lastJdUT))
			if err != nil {
				return 0, 0, fmt.Errorf("LunarSolarEclipticLongitudeDeltaToTime Calc 2: %w", err)
			}

			lastDeltaDelta = RadiansMod360(lastDelta - firstDelta)
			calcCount++
		}
	} else {
		// 已经相等了，说明传入的时间即当时的弧差
		if FloatEqual(lastDelta, eclipticLongitudeDelta, 9) {
			return lastJdUT, calcCount, nil
		}

		// 取当前 黄经差 与 所需黄经差 的 差值
		lastDeltaDelta := RadiansMod180(lastDelta - eclipticLongitudeDelta)

		var savedJdUT, meanJdUT JulianDay
		var savedDelta, savedDeltaDelta float64
		for {
			savedJdUT = lastJdUT
			savedDelta = lastDelta
			savedDeltaDelta = lastDeltaDelta

			// 渐步递增
			lastJdUT = lastJdUT.Add(step)
			lastDelta, lastSpeedDelta, err = astro.LunarSolarEclipticLongitudeDelta(NewEphemerisTime(lastJdUT))
			if err != nil {
				return 0, calcCount, fmt.Errorf("LunarSolarEclipticLongitudeDeltaToTime Calc 2: %w", err)
			}
			lastDeltaDelta = RadiansMod180(lastDelta - eclipticLongitudeDelta)

			calcCount++

			// Keep searching while error is large.
			if math.Abs(lastDeltaDelta) > Radian90 {
				continue
			}

			// 这两个的正负不同，说明最后一次计算超过了正确值，获取上一次saved的结果
			if !SameSign(savedDeltaDelta, lastDeltaDelta) {
				break
			}
		}

		// 利用JD平均数来缩小差距，二分法
		for !FloatEqual(float64(lastJdUT), float64(savedJdUT), 6) {

			// 两日期的平均值
			meanJdUT = (savedJdUT + lastJdUT) / 2.
			var lastDelta2 float64
			lastDelta2, lastSpeedDelta, err = astro.LunarSolarEclipticLongitudeDelta(NewEphemerisTime(meanJdUT))
			if err != nil {
				return 0, calcCount, fmt.Errorf("LunarSolarEclipticLongitudeDeltaToTime Calc 3: %w", err)
			}
			// 求两"差"之差
			lastDeltaDelta2 := RadiansMod180(lastDelta2 - eclipticLongitudeDelta)

			// 符号相同
			if SameSign(lastDeltaDelta2, lastDeltaDelta) {
				lastDelta = lastDelta2
				lastJdUT = meanJdUT
				lastDeltaDelta = lastDeltaDelta2
			} else {
				savedDelta = lastDelta2
				savedJdUT = meanJdUT
				savedDeltaDelta = lastDeltaDelta2
			}

			calcCount++

		}

		meanJdUT = savedJdUT.Add(float64(lastJdUT-savedJdUT) * -savedDeltaDelta / (lastDelta - savedDelta))

		return meanJdUT, calcCount, nil
	}

}

/**
 * 根据传入的月日黄经差值(数组), 返回对应的时间(数组)
 * startJdUT 起始时间
 * eclipticLongitudeDeltas 月日黄经差值 弧度，
 * 比如(这里以度表示，实际应该是弧度)：{90, 180, 270, 0, 90} 表示startJdUT之后的{上弦，望，下弦，下一个朔，下一个上弦}
 * 也就是说可以计算多个月的，下面的 LunarPhases，使用本函数一次计算了1年
 */
func (astro *Astronomy) LunarSolarEclipticLongitudeDeltaToTimes(startJdUT JulianDay, eclipticLongitudeDeltas []float64) ([]JulianDay, error) {
	times := make([]JulianDay, len(eclipticLongitudeDeltas))

	var jd = startJdUT
	var err error
	lastDelta := math.Inf(-1)
	for i, delta := range eclipticLongitudeDeltas {
		// 如果上一个弧度差和现在的相等，为了避免LunarSolarEclipticLongitudeDeltaToTime不往后面推进，jd累加半个月
		if FloatEqual(delta, lastDelta, 9) {
			jd += MeanLunarDays / 2
		}

		jd, _, err = astro.LunarSolarEclipticLongitudeDeltaToTime(NewEphemerisTime(jd), ToRadians(delta))
		if err != nil {
			return nil, fmt.Errorf("LunarSolarEclipticLongitudeDeltaToTimes: %w", err)
		}
		times[i] = jd

		lastDelta = delta
	}

	return times, nil
}

/**
 * 2时间之间的月相的时间，注意：只会返回如下月相：朔、上弦、望、下弦
 * startJdUT 起始时间
 * startJdUT 结束时间
 */
func (astro *Astronomy) LunarPhasesRange(startJdUT, endJdUT JulianDay) ([]*JulianDayExtra, error) {

	// 90° 每个节气
	const degreePerLunarPhases = 90
	firstLongDelta, _, err := astro.LunarSolarEclipticLongitudeDelta(NewEphemerisTime(startJdUT))
	if err != nil {
		return nil, fmt.Errorf("LunarPhases Calc 1: %w", err)
	}
	// 第一个有效的角度
	firstValidDelta := NextMultiples(ToDegrees(firstLongDelta), float64(degreePerLunarPhases))

	// 计算1年大致有多少个朔望上下弦，多给1个
	count := int(float64(endJdUT-startJdUT)/MeanLunarDays)*4 + 1
	eclipticLongitudeDelta := make([]float64, count)
	lunarTimes := make([]*JulianDayExtra, count)

	// 得到他们对应的角度
	for i, _ := range lunarTimes {
		delta := float64(degreePerLunarPhases*i) + firstValidDelta
		eclipticLongitudeDelta[i] = delta
		lunarTimes[i] = NewJulianDayExtra(0, int(delta)/degreePerLunarPhases%4)
	}

	// 运算
	times, err := astro.LunarSolarEclipticLongitudeDeltaToTimes(startJdUT, eclipticLongitudeDelta)
	if err != nil {
		return nil, fmt.Errorf("LunarPhases Calc 2: %w", err)
	}

	for i, jd := range times {
		lunarTimes[i].JdUT = jd
	}

	// 最后一个月的超过了结束时间，不返回
	if times[count-1] > endJdUT {
		return lunarTimes[0 : count-2], nil
	}

	return lunarTimes, nil
}

/**
 * 某年的月相
 * year 年
 */
func (astro *Astronomy) LunarPhases(year int) ([]*JulianDayExtra, error) {
	jd := DateToJulianDay(year, 1, 1, 0, 0, 0)
	return astro.LunarPhasesRange(jd, jd.AddYears(1))
}

/**
 * 某时间之后的朔日（数组）
 * startJdUT 起始时间
 * count 计算多少个朔日
 */
func (astro *Astronomy) NewMoons(startJdUT JulianDay, count int) ([]JulianDay, error) {

	eclipticLongitudeDeltas := make([]float64, count)
	for i := range eclipticLongitudeDeltas {
		eclipticLongitudeDeltas[i] = 0.
	}

	times, err := astro.LunarSolarEclipticLongitudeDeltaToTimes(startJdUT, eclipticLongitudeDeltas)
	if err != nil {
		return nil, fmt.Errorf("NewMoons: %w", err)
	}

	return times, nil
}

/**
 * 某时间之前的第一个朔日
 * startJdUT 起始时间
 */
func (astro *Astronomy) LastNewMoons(startJdUT JulianDay) (JulianDay, error) {
	times, err := astro.NewMoons(startJdUT-MeanLunarDays*1.5, 3)
	if err != nil {
		return 0, fmt.Errorf("LastNewMoons: %w", err)
	}
	// 第一个比startJdUT小的朔日
	for i := range times {
		n := len(times) - 1 - i
		if times[n] <= startJdUT {
			return times[n], nil
		}
	}

	// 原则上不可能出现这种错误
	return 0, fmt.Errorf("LastNewMoons: Unknown error")
}
