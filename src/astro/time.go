package astro

import (
	"go-swe/src/swe"
	"math"
	"time"
)

// JulianDay 儒略日 儒略日是无时区意义的, 即世界时(UT)
type JulianDay float64
type MJD float64
type JD2000 float64
type JulianDayWithLocation float64

const JD2000_OFFSET float64 = 2451545.
const JD_CST_OFFSET = 8. / 24.

type EphemerisTime struct {
	JdUT   JulianDay `json:"jd_ut"`
	DeltaT float64   `json:"delta_t"`
}

type JulianDayExtra struct {
	JdUT  JulianDay `json:"jd_ut"`
	Index int       `json:"index"`
}

// TimeToJulianDay time 转为 儒略日 JulianDay
// 带时区的时间会被自动转为UTC
func TimeToJulianDay(_time time.Time) JulianDay {
	_time = _time.UTC()
	_hour := float64(_time.Hour()) + float64(_time.Minute())/60. + float64(_time.Second()+_time.Nanosecond()/1e9)/3600.
	jd, _ := swe.NewSwe().JulDay(_time.Year(), int(_time.Month()), _time.Day(), _hour, swe.Gregorian)

	return JulianDay(jd)
}

// MakeJulianDayHours 将 时,分,秒 (只能为UTC) 转化为 儒略日的小数部分 JulianDay
func MakeJulianDayHours(hour, minute int, second float64) float64 {
	return float64(hour) + float64(minute)/60. + second/3600.
}

// ExtractJulianDayHours 将 儒略日的小数部分 转化为 时,分,秒 (UTC)
func ExtractJulianDayHours(hours float64) (hour, minute int, second float64) {
	hour = int(hours)
	minute = int(math.Mod(hours*60, 60))
	second = math.Mod(hours*3600, 60)

	return
}

// DateToJulianDay 年,月,日,时,分,秒 (只能为UTC) 转换为 儒略日 JulianDay
func DateToJulianDay(year, month, day, hour, minute int, second float64) JulianDay {
	_hour := MakeJulianDayHours(hour, minute, second)

	jd, _ := swe.NewSwe().JulDay(year, month, day, _hour, swe.Gregorian)

	return JulianDay(jd)
}

// ExtractJulianDay 儒略日 JulianDay 转化为 年,月,日,时,分,秒 (UTC)
func ExtractJulianDay(jd JulianDay) (year, month, day, hour, minute int, second float64) {
	year, month, day, hours, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	hour, minute, second = ExtractJulianDayHours(hours)
	return
}

func NewEphemerisTime(jdUT JulianDay) *EphemerisTime {
	return (&EphemerisTime{}).Update(jdUT)
}

func (et *EphemerisTime) Value() float64 {
	return float64(et.JdUT.Add(et.DeltaT))
}

func (et *EphemerisTime) Update(jdUT JulianDay) *EphemerisTime {
	et.JdUT = jdUT
	et.DeltaT = DeltaT(jdUT) //.DeltaT(float64(jdUT))
	return et
}

func NewJulianDayExtra(jdUT JulianDay, index int) *JulianDayExtra {
	return &JulianDayExtra{JdUT: jdUT, Index: index}
}

// ToTime Ephemeris time 天文历时 jd转为local的 time.Time, 此时jd的时区
func (jd JulianDay) ToTime(local *time.Location) time.Time {
	if local == nil {
		local = time.UTC
	}
	year, month, day, hour, minute, second := ExtractJulianDay(jd)
	date := time.Date(year, time.Month(month), day, hour, minute, int(second), int((second-float64(int(second)))*1e9), time.UTC)
	return date.In(local)
}

// ToCST JD UT 转 CST
func (jd JulianDay) ToCST() JulianDayWithLocation {
	return jd.ToLocation(JD_CST_OFFSET)

}

func (jd JulianDay) ToLocation(offset float64) JulianDayWithLocation {
	return JulianDayWithLocation(float64(jd) + offset)
}

// ToJD2000 转换为JD2000的表示方式
func (jd JulianDay) ToJD2000() JD2000 {
	return JD2000(jd - JulianDay(JD2000_OFFSET))
}

// ToJulianDay JD2000转化为JD的表示方式
func (jd2000 JD2000) ToJulianDay() JulianDay {
	return JulianDay(jd2000 + JD2000(JD2000_OFFSET))
}

// ToMJD 转换为MJD的表示方式
func (jd JulianDay) ToMJD() MJD {
	return MJD(jd - 2400000.5)
}

// ToJulianDay MJD转化为JD的表示方式
func (mjd MJD) ToJulianDay() JulianDay {
	return JulianDay(mjd + 2400000.5)
}

// ToJulianDay 带时区的JD时间转为JD UT
func (jdz JulianDayWithLocation) ToJulianDay(offset float64) JulianDay {
	return JulianDay(float64(jdz) - offset)
}

// Add 增加日期，float64 的生成规则和儒略日一致
// 减少请传入负数
func (jd JulianDay) Add(delta float64) JulianDay {
	return jd + JulianDay(delta)
}

// AddYears 增加N年，减少用负数
// 注意：大部分日历操作类都存在这个问题 2000-02-29 + 1 year -> 2001-03-01
func (jd JulianDay) AddYears(years int) JulianDay {
	y, m, d, t, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y+years, m, d, t, swe.Gregorian)
	return JulianDay(_jd)
}

// AddMonths 增加N月，减少用负数
// 注意：大部分日历操作类都存在这个问题 2000-03-31 - 1 month -> 2000-03-02
func (jd JulianDay) AddMonths(months int) JulianDay {
	y, m, d, t, _ := swe.NewSwe().RevJul(float64(jd), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m+months, d, t, swe.Gregorian)
	return JulianDay(_jd)
}

// AddDays 增加N日，减少用负数
func (jd JulianDay) AddDays(days int) JulianDay {
	return jd + JulianDay(days)
}

// AddHours 增加N小时，减少用负数
func (jd JulianDay) AddHours(hours int) JulianDay {
	return jd + JulianDay(MakeJulianDayHours(hours, 0, 0))
}

// AddMinutes 增加N分钟，减少用负数
func (jd JulianDay) AddMinutes(minutes int) JulianDay {
	return jd + JulianDay(MakeJulianDayHours(0, minutes, 0))
}

// AddSeconds 增加N秒，减少用负数
func (jd JulianDay) AddSeconds(seconds float64) JulianDay {
	return jd + JulianDay(MakeJulianDayHours(0, 0, seconds))
}

// Midnight 午夜零点
// the nature day's 00:00
func (jdz JulianDayWithLocation) Midnight() JulianDayWithLocation {
	return jdz.StartOfDay()
}

// Noon 正午
// the nature day's 12:00
func (jdz JulianDayWithLocation) Noon() JulianDayWithLocation {
	y, m, d, _, _ := swe.NewSwe().RevJul(float64(jdz), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m, d, 12, swe.Gregorian)
	return JulianDayWithLocation(_jd)
}

// StartOfDay 今天的零点，午夜 00:00:00
func (jdz JulianDayWithLocation) StartOfDay() JulianDayWithLocation {
	y, m, d, _, _ := swe.NewSwe().RevJul(float64(jdz), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m, d, 0, swe.Gregorian)
	return JulianDayWithLocation(_jd)
}

// EndOfDay 今天的 23:59:59
func (jdz JulianDayWithLocation) EndOfDay() JulianDayWithLocation {
	y, m, d, _, _ := swe.NewSwe().RevJul(float64(jdz), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m, d, MakeJulianDayHours(23, 59, 59), swe.Gregorian)
	return JulianDayWithLocation(_jd)
}

// StartOfMonth 月初 XXXX-XX-01 00:00:00
func (jdz JulianDayWithLocation) StartOfMonth() JulianDayWithLocation {
	y, m, _, _, _ := swe.NewSwe().RevJul(float64(jdz), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m, 1, 0, swe.Gregorian)
	return JulianDayWithLocation(_jd)
}

// EndOfMonth 月尾 XXXX-XX-日 23:59:59 其中，日可能为：28,29,30,31
func (jdz JulianDayWithLocation) EndOfMonth() JulianDayWithLocation {
	return jdz.StartOfMonth().AddMonths(1).AddSeconds(-1)
}

// StartOfYear 年初 XXXX-01-01 00:00:00
func (jdz JulianDayWithLocation) StartOfYear() JulianDayWithLocation {
	y, _, _, _, _ := swe.NewSwe().RevJul(float64(jdz), swe.Gregorian)
	_jdz, _ := swe.NewSwe().JulDay(y, 1, 1, 0, swe.Gregorian)
	return JulianDayWithLocation(_jdz)
}

// EndOfYear 年尾 XXXX-12-31 23:59:59
func (jdz JulianDayWithLocation) EndOfYear() JulianDayWithLocation {
	y, _, _, _, _ := swe.NewSwe().RevJul(float64(jdz), swe.Gregorian)
	_jdz, _ := swe.NewSwe().JulDay(y, 12, 31, MakeJulianDayHours(23, 59, 59), swe.Gregorian)
	return JulianDayWithLocation(_jdz)
}

// Add 增加日期，float64 的生成规则和儒略日一致
// 减少请传入负数
func (jdz JulianDayWithLocation) Add(delta float64) JulianDayWithLocation {
	return jdz + JulianDayWithLocation(delta)
}

// AddYears 增加N年，减少用负数
// 注意：大部分日历操作类都存在这个问题 2000-02-29 + 1 year -> 2001-03-01
func (jdz JulianDayWithLocation) AddYears(years int) JulianDayWithLocation {
	y, m, d, t, _ := swe.NewSwe().RevJul(float64(jdz), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y+years, m, d, t, swe.Gregorian)
	return JulianDayWithLocation(_jd)
}

// AddMonths 增加N月，减少用负数
// 注意：大部分日历操作类都存在这个问题 2000-03-31 - 1 month -> 2000-03-02
func (jdz JulianDayWithLocation) AddMonths(months int) JulianDayWithLocation {
	y, m, d, t, _ := swe.NewSwe().RevJul(float64(jdz), swe.Gregorian)
	_jd, _ := swe.NewSwe().JulDay(y, m+months, d, t, swe.Gregorian)
	return JulianDayWithLocation(_jd)
}

// AddDays 增加N日，减少用负数
func (jdz JulianDayWithLocation) AddDays(days int) JulianDayWithLocation {
	return jdz + JulianDayWithLocation(days)
}

// AddHours 增加N小时，减少用负数
func (jdz JulianDayWithLocation) AddHours(hours int) JulianDayWithLocation {
	return jdz + JulianDayWithLocation(MakeJulianDayHours(hours, 0, 0))
}

// AddMinutes 增加N分钟，减少用负数
func (jdz JulianDayWithLocation) AddMinutes(minutes int) JulianDayWithLocation {
	return jdz + JulianDayWithLocation(MakeJulianDayHours(0, minutes, 0))
}

// AddSeconds 增加N秒，减少用负数
func (jdz JulianDayWithLocation) AddSeconds(seconds float64) JulianDayWithLocation {
	return jdz + JulianDayWithLocation(MakeJulianDayHours(0, 0, seconds))
}

// GreenwichMeridianSiderealTime 格林尼治平恒星时(不含赤经章动及非多项式部分),即格林尼治子午圈的平春分点起算的赤经
// 可以使用 astro.Swe.SidTime 精确计算
func GreenwichMeridianSiderealTime(jdUT JulianDay, deltaT float64) float64 {
	//t是力学时(世纪数)
	t := float64(jdUT.Add(deltaT).ToJD2000()) / 36525.
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t

	return Radian360*
		(0.7790572732640+1.00273781191135448*(float64(jdUT)-2451545.)) +
		(0.014506+4612.15739966*t+1.39667721*t2-0.00009344*t3+0.00001882*t4)/DegreeSecondsPerRadian
}

// TD - UT1 计算表
var dtAt = [...]float64{
	-4000, 108371.7, -13036.80, 392.000, 0.0000,
	-500, 17201.0, -627.82, 16.170, -0.3413,
	-150, 12200.6, -346.41, 5.403, -0.1593,
	150, 9113.8, -328.13, -1.647, 0.0377,
	500, 5707.5, -391.41, 0.915, 0.3145,
	900, 2203.4, -283.45, 13.034, -0.1778,
	1300, 490.1, -57.35, 2.085, -0.0072,
	1600, 120.0, -9.81, -1.532, 0.1403,
	1700, 10.2, -0.91, 0.510, -0.0370,
	1800, 13.4, -0.72, 0.202, -0.0193,
	1830, 7.8, -1.81, 0.416, -0.0247,
	1860, 8.3, -0.13, -0.406, 0.0292,
	1880, -5.4, 0.32, -0.183, 0.0173,
	1900, -2.3, 2.06, 0.169, -0.0135,
	1920, 21.2, 1.69, -0.304, 0.0167,
	1940, 24.2, 1.22, -0.064, 0.0031,
	1960, 33.2, 0.51, 0.231, -0.0109,
	1980, 51.0, 1.29, -0.026, 0.0032,
	2000, 63.87, 0.1, 0, 0,
	2005, 64.7, 0.4, 0, 0, // 一次项记为x,则 10x=0.4秒/年*(2015-2005),解得x=0.4
	2015, 69,
}

// 二次曲线外推
func dtExt(y, jsd float64) float64 {
	var dy = (y - 1820) / 100
	return -20 + jsd*dy*dy
}

// 计算世界时与原子时之差,传入年
func dtCalc(y float64) float64 {
	dtAtLen := len(dtAt)

	var y0 = dtAt[dtAtLen-2] // 表中最后一年
	var t0 = dtAt[dtAtLen-1] // 表中最后一年的deltaT
	if y >= y0 {
		var jsd = 31. // sjd是y1年之后的加速度估计。瑞士星历表jsd=31,NASA网站jsd=32,skmap的jsd=29
		if y > y0+100 {
			return dtExt(y, jsd)
		}
		var v = dtExt(y, jsd)        // 二次曲线外推
		var dv = dtExt(y0, jsd) - t0 // ye年的二次外推与te的差
		return v - dv*(y0+100-y)/100
	}
	var i int
	for i = 0; i < dtAtLen; i += 5 {
		if y < dtAt[i+5] {
			break
		}
	}
	t1 := (y - dtAt[i]) / (dtAt[i+5] - dtAt[i]) * 10
	t2 := t1 * t1
	t3 := t2 * t1
	return dtAt[i+1] + dtAt[i+2]*t1 + dtAt[i+3]*t2 + dtAt[i+4]*t3
}

// DeltaT 计算UT 和 ET 的 DeltaT
// 可以使用 astro.Swe.DeltaT 精确计算
func DeltaT(jdUT JulianDay) float64 {
	jd := jdUT.ToJD2000()
	return dtCalc(float64(jd)/365.2425+2000) / 86400.0
}
