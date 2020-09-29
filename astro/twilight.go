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
	Noon, Midnight float64
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
type TwilightTime struct {
	// 升天, 降天
	Rise, Set JulianDay
	// 上中天，下中天
	Noon, Midnight JulianDay
}

// 天体的 晨/暮 时间
type DawnDuskTime struct {
	// 晨
	Dawn JulianDay
	// 暮
	Dusk JulianDay
}

// 太阳的 升/降/晨/暮 时间
type SunTwilightTime struct {
	TwilightTime

	Astronomical DawnDuskTime
	Nautical     DawnDuskTime
	Civil        DawnDuskTime
}

func NewSunTwilightAngle() *SunTwilightAngle {
	return &SunTwilightAngle{
		TwilightAngle: TwilightAngle{
			RiseSet:  ToRadian(-50. / 60.), // 地平线以下 50′
			Noon:     ToRadian(90.),
			Midnight: ToRadian(-90.),
		},
		Astronomical: ToRadian(-18.),
		Nautical:     ToRadian(-12.),
		Civil:        ToRadian(-6.),
	}
}

/**
 * 日照时长
 * 日落 - 升日
 */
func (stt *SunTwilightTime) Daylight() float64 {
	return float64(stt.Set - stt.Rise)
}

/**
 * 夜晚时长，因为跨午夜，所以取今日00:00 ~ 23:59的时段
 * (0:00 ~ 天文晨光) + (天文暮光 ~ 23:59)
 */
func (stt *SunTwilightTime) Night() float64 {
	return float64(
		stt.Astronomical.Dawn - stt.Astronomical.Dawn.Midnight() +
			(stt.Astronomical.Dusk.AddDays(1).Midnight() - stt.Astronomical.Dusk),
	)
}

/**
 * 太阳
 * 传入本地12点的JdUT(比如东八区需要4点: TimeToJulianDay(2020-09-30 04:00:00))，地理坐标
 * withRevise: 是否修正一些日光差，或者黄道章动
 */
func SunTwilight(jdUT JulianDay, geo *GeographicCoordinates, withRevise bool) *SunTwilightTime {

	// 查找最靠近当日中午的日上中天, mod2的第1参数为本地时角近似值
	noonJdUT := jdUT.Add(-Mod2(float64(jdUT)+geo.Longitude/Radian360, 1))

	angle := NewSunTwilightAngle()
	sunTime := &SunTwilightTime{}

	astro := NewAstronomy(geo, noonJdUT)

	// 天体属性
	planet, err := astro.PlanetProperties(swe.Sun, withRevise)

	if err != nil {
		return nil
	}

	type haCallback func(hourAngle HourAngle) HourAngle
	var negativeCallback = func(hourAngle HourAngle) HourAngle {
		return -hourAngle
	}
	var positiveCallback = func(hourAngle HourAngle) HourAngle {
		return hourAngle
	}
	var zeroCallback = func(hourAngle HourAngle) HourAngle {
		return 0
	}
	var piCallback = func(hourAngle HourAngle) HourAngle {
		return HourAngle(Radian180)
	}

	var revise = func(planet *PlanetProperties, jd JulianDay, angle float64, callback haCallback) (_jd JulianDay) {
		_jd = jd
		var ha, _ha HourAngle

		// 计算第一次
		if math.Abs(angle) < Radian90 {
			ha = AltitudeToHourAngle(planet.Equatorial.Declination, geo.Latitude, angle)
		}
		ha = callback(ha)

		_jd = _jd.Add(float64(ha-planet.HourAngle) / Radian360)

		// 多修正几次
		for i := 0; i < 3; i++ {
			_planet, _err := NewAstronomy(geo, _jd).PlanetProperties(swe.Sun, withRevise)
			if _err != nil {
				return
			}

			if math.Abs(angle) < Radian90 {
				_ha = AltitudeToHourAngle(_planet.Equatorial.Declination, geo.Latitude, angle)
			}
			_ha = callback(_ha)

			_jd = _jd.Add(RadianMod180(float64(_ha-_planet.HourAngle)) / Radian360)
		}

		return
	}

	sunTime.Rise = revise(planet, noonJdUT, angle.RiseSet, negativeCallback)
	sunTime.Set = revise(planet, noonJdUT, angle.RiseSet, positiveCallback)
	sunTime.Civil.Dawn = revise(planet, noonJdUT, angle.Civil, negativeCallback)
	sunTime.Civil.Dusk = revise(planet, noonJdUT, angle.Civil, positiveCallback)
	sunTime.Nautical.Dawn = revise(planet, noonJdUT, angle.Nautical, negativeCallback)
	sunTime.Nautical.Dusk = revise(planet, noonJdUT, angle.Nautical, positiveCallback)
	sunTime.Astronomical.Dawn = revise(planet, noonJdUT, angle.Astronomical, negativeCallback)
	sunTime.Astronomical.Dusk = revise(planet, noonJdUT, angle.Astronomical, positiveCallback)
	sunTime.Noon = revise(planet, noonJdUT, angle.Noon, zeroCallback)
	sunTime.Midnight = revise(planet, noonJdUT, angle.Midnight, piCallback)

	return sunTime
}
