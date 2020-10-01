package astro

import (
	"go-swe/swe"
	"math"
)

var LunarPhaseStrings = [...]string{"朔", "上弦", "望", "下弦"}
var LunarMonthStrings = [...]string{"正月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "腊月"}
var LunarDayStrings = [...]string{
	"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九",
	"初十", "十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九",
	"二十", "廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九",
	"三十", "卅一",
}

/**
 * 太阳和月亮的黄经之差
 */
func (astro *Astronomy) SolarLunarEclipticLongitudeDelta(jdET *EphemerisTime) (float64, error) {
	sun, err := astro.PlanetProperties(swe.Sun, jdET)
	if err != nil {
		return 0, err
	}

	moon, err := astro.PlanetProperties(swe.Moon, jdET)
	if err != nil {
		return 0, err
	}

	return moon.Ecliptic.Longitude - sun.Ecliptic.Longitude, nil
}

/**
 * eclipticLongitudeDelta 太阳和月球的黄经之差，比如 朔（0），上弦（90°），望（180°），下弦（270°）
 */
func (astro *Astronomy) SolarLunarEclipticLongitudeDeltaToTime(jdET *EphemerisTime, eclipticLongitudeDelta float64) (JulianDay, error) {
	eclipticLongitudeDelta = math.Abs(eclipticLongitudeDelta)

	step := 27. * 0.125 // 4/1000,每次递进一小步
	lastJdUT := jdET.JdUT

	lastDelta, err := astro.SolarLunarEclipticLongitudeDelta(jdET)
	if err != nil {
		return 0, err
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
		lastDelta, err = astro.SolarLunarEclipticLongitudeDelta(NewEphemerisTime(lastJdUT))
		if err != nil {
			return 0, err
		}
		lastDeltaDelta = RadiansMod180(lastDelta - eclipticLongitudeDelta)

		// Keep searching while error is large.
		if math.Abs(lastDeltaDelta) > Radian90 {
			continue
		}

		// 这两个的正负不同，说明最后一次计算超过了正确值，获取上一次saved的结果
		if (savedDeltaDelta * lastDeltaDelta) <= 0. {
			break
		}
	}

	// 利用JD平均数来缩小差距，二分法
	for !FloatEqual(float64(lastJdUT), float64(savedJdUT), 9) {

		// 两日期的平均值
		meanJdUT = (savedJdUT + lastJdUT) / 2.
		lastDelta2, err := astro.SolarLunarEclipticLongitudeDelta(NewEphemerisTime(meanJdUT))
		if err != nil {
			return 0, err
		}
		// 求两"差"之差
		lastDeltaDelta2 := RadiansMod180(lastDelta2 - eclipticLongitudeDelta)

		if lastDeltaDelta2*lastDeltaDelta > 0 {
			lastDelta = lastDelta2
			lastJdUT = meanJdUT
			lastDeltaDelta = lastDeltaDelta2
		} else {
			savedDelta = lastDelta2
			savedJdUT = meanJdUT
			savedDeltaDelta = lastDeltaDelta2
		}

	}

	meanJdUT = savedJdUT.Add(float64(lastJdUT-savedJdUT) * -savedDeltaDelta / (lastDelta - savedDelta))

	return meanJdUT, nil
}

/**
 * 根据传入的月日黄经差值(数组), 返回对应的时间(数组)
 * startJdUT 起始时间
 * eclipticLongitudeDelta 月日黄经差值 弧度，
 * 比如(这里以度表示，实际应该是弧度)：{90, 180, 270, 0, 90} 表示startJdUT之后的{上弦，望，下弦，朔，上弦}。前面3个{上弦，望，下弦}是当前阴历月的, 最后2个{朔，上弦}是下个阴历月的
 * 也就是说可以计算多个月的，下面的 LunarPhases，使用本函数一次计算了1年
 */
func (astro *Astronomy) SolarLunarEclipticLongitudeDeltaToTimes(startJdUT JulianDay, eclipticLongitudeDelta []float64) ([]JulianDay, error) {
	times := make([]JulianDay, len(eclipticLongitudeDelta))

	var jd = startJdUT
	var err error
	for i, delta := range eclipticLongitudeDelta {
		jd, err = astro.SolarLunarEclipticLongitudeDeltaToTime(NewEphemerisTime(jd), ToRadians(delta))
		if err != nil {
			return nil, err
		}
		times[i] = jd
	}

	return times, nil
}

/**
 * 当年的所有月相的时间
 * year 年
 * 可能在数组末尾有次年的数据
 */
func (astro *Astronomy) LunarPhases(year int) ([]*JulianDayWithIndex, error) {

	// 90° 每个节气
	degreePerLunarPhases := 90
	//
	jdUT := DateToJulianDay(year, 1, 1, 0, 0, 0)
	nextYear := jdUT.AddYears(1)
	firstLongDelta, err := astro.SolarLunarEclipticLongitudeDelta(NewEphemerisTime(jdUT))
	if err != nil {
		return nil, err
	}
	// 第一个有效的角度
	firstValidDelta := NextMultiples(ToDegrees(firstLongDelta), float64(degreePerLunarPhases))

	// 计算1年大致有多少个朔望上下弦，多给1个
	eclipticLongitudeDelta := make([]float64, int(math.Ceil(float64(4*(nextYear-jdUT)/MeanLunarDays)+1)))
	lunarTimes := make([]*JulianDayWithIndex, len(eclipticLongitudeDelta))

	// 得到他们对应的角度
	for i, _ := range lunarTimes {
		delta := float64(degreePerLunarPhases*i) + firstValidDelta
		eclipticLongitudeDelta[i] = delta
		lunarTimes[i] = NewJulianDayWithIndex(0, int(delta)/degreePerLunarPhases%4)
	}

	// 运算
	times, err := astro.SolarLunarEclipticLongitudeDeltaToTimes(jdUT, eclipticLongitudeDelta)
	if err != nil {
		return nil, err
	}

	for i, jd := range times {
		lunarTimes[i].JdUT = jd
	}

	return lunarTimes, nil
}
