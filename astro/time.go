package astro

import (
	"go-swe/swe"
	"math"
	"time"
)

type JulianDay float64
type JulianDayDelta float64

// time 转为 儒略日 JulianDay
// 如果需要UT的儒略日，需要传递_time.UTC()
func TimeToJulianDay(_time time.Time) JulianDay {
	_hour := float64(_time.Hour()) + float64(_time.Minute()) / 60. + float64(_time.Nanosecond()) / 3600_000_000_000.
	jd, _ := swe.NewSwe().JulDay(_time.Year(), int(_time.Month()), _time.Day(), _hour, swe.Gregorian)

	return JulianDay(jd)
}

// 将 时,分,秒 转化为 儒略日的小数部分 JulianDay
func MakeJulianDayTime(hour, minute int, second float64) float64 {
	return float64(hour) + float64(minute) / 60. + second / 3600.
}

// 将 儒略日的小数部分 转化为 时,分,秒
func ExtractJulianDayTime(hours float64) (hour, minute int, second float64) {
	hour = int(hours)
	minute = int( math.Mod(hours * 60, 60) )
	second = math.Mod(hours * 3600, 60)

	return
}

// 年,月,日,时,分,秒 转换为 儒略日 JulianDay
func DateToJulianDay(year, month, day, hour, minute int, second float64) JulianDay {
	_hour := MakeJulianDayTime(hour, minute, second)

	jd, _ := swe.NewSwe().JulDay(year, month, day, _hour, swe.Gregorian)

	return JulianDay(jd)
}

// 儒略日 JulianDay 转化为 年,月,日,时,分,秒
func ExtractJulianDay(jd JulianDay) (year, month, day, hour, minute int, second float64) {
	year, month, day, hours, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	hour, minute, second = ExtractJulianDayTime(hours)
	return
}

// Ephemeris time 天文历时
func (jd JulianDay) ToEphemerisTime(deltaT JulianDayDelta) JulianDay {
	return jd.Add(deltaT)
}

// 增加日期，JulianDayDelta 的生成规则和儒略日一致
// 减少请传入负数
func (jd JulianDay) Add(delta JulianDayDelta) JulianDay {
	return jd + JulianDay(delta)
}

// 增加N年，减少用负数
// 注意：大部分日历操作类都存在这个问题 2000-02-29 + 1 year -> 2001-03-01
func (jd JulianDay) AddYears(years int) JulianDay {
	y, m, d, t, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y + years, m, d, t, swe.Gregorian)
	return JulianDay(_jd)
}

// 增加N年，减少用负数
// 注意：大部分日历操作类都存在这个问题 2000-03-31 - 1 month -> 2000-03-02
func (jd JulianDay) AddMonths(months int) JulianDay {
	y, m, d, t, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m + months, d, t, swe.Gregorian)
	return JulianDay(_jd)
}

// 增加N日，减少用负数
func (jd JulianDay) AddDays(days int) JulianDay {
	return jd + JulianDay(days)
}

// 增加N小时，减少用负数
func (jd JulianDay) AddHours(hours int) JulianDay {
	return jd + JulianDay(MakeJulianDayTime(hours, 0, 0))
}

// 增加N分钟，减少用负数
func (jd JulianDay) AddMinutes(minutes int) JulianDay {
	return jd + JulianDay(MakeJulianDayTime(0, minutes, 0))
}

// 增加N秒，减少用负数
func (jd JulianDay) AddSeconds(seconds float64) JulianDay {
	return jd + JulianDay(MakeJulianDayTime(0, 0, seconds))
}

// 午夜零点
// the nature day's 00:00
func (jd JulianDay) Midnight() JulianDay  {
	return jd.StartOfDay()
}

// 正午
// the nature day's 12:00
func (jd JulianDay) Noon() JulianDay  {
	y, m, d, _, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m, d, 12, swe.Gregorian)
	return JulianDay(_jd)
}

// 今天的零点，午夜 00:00:00
func (jd JulianDay) StartOfDay() JulianDay  {
	y, m, d, _, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m, d, 0, swe.Gregorian)
	return JulianDay(_jd)
}

// 今天的 23:59:59
func (jd JulianDay) EndOfDay() JulianDay  {
	y, m, d, _, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m, d, MakeJulianDayTime(23, 59, 59), swe.Gregorian)
	return JulianDay(_jd)
}

// 月初 XXXX-XX-01 00:00:00
func (jd JulianDay) StartOfMonth() JulianDay  {
	y, m, _, _, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m, 1, 0, swe.Gregorian)
	return JulianDay(_jd)
}

// 月尾 XXXX-XX-日 23:59:59 其中，日可能为：28,29,30,31
func (jd JulianDay) EndOfMonth() JulianDay  {
	return jd.StartOfMonth().AddMonths(1).AddSeconds(-1)
}

// 年初 XXXX-01-01 00:00:00
func (jd JulianDay) StartOfYear() JulianDay  {
	y, _, _, _, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, 1, 1, 0, swe.Gregorian)
	return JulianDay(_jd)
}

// 年尾 XXXX-12-31 23:59:59
func (jd JulianDay) EndOfYear() JulianDay  {
	y, _, _, _, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, 12, 31, MakeJulianDayTime(23, 59, 59), swe.Gregorian)
	return JulianDay(_jd)
}

// 格林尼治恒星时(不含赤经章动及非多项式部分),即格林尼治子午圈的平春分点起算的赤经
func GreenwichMeridianSiderealTime(jdET JulianDay) float64 {
	//t是力学时(世纪数)
	t := (float64(jdET) - 2451545.) / 36525.
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t

	return Radian360*
		(0.7790572732640 + 1.00273781191135448 * (float64(jdET) - 2451545.)) +
		(0.014506 + 4612.15739966 * t + 1.39667721 * t2 - 0.00009344 * t3 + 0.00001882 * t4) /
			DegreeSecondsPerRadian
}