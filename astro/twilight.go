package astro

import (
	"go-swe/swe"
	"math"
)

// 星星的 升/降 角度
type TwilightAngle struct {
	// 升天, 降天
	Rise, Set,
	// 上中天，下中天
	Noon, Midnight float64
}

// 星星的 晨/暮 角度
type DawnDuskAngle struct {
	// 晨光 暮光
	Dawn, Dusk float64
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
	// 各个专业领域的角度
	Astronomical DawnDuskAngle
	Nautical DawnDuskAngle
	Civil DawnDuskAngle
}

// 星星的 升/降 时间
type TwilightTime struct {
	// 升天, 降天
	Rise, Set JulianDay
	// 上中天，下中天
	Noon, Midnight JulianDay
}

// 星星的 晨/暮 时间
type DawnDuskTime struct {
	// 晨光 暮光
	Dawn, Dusk JulianDay
}

// 太阳的 升/降/晨/暮 时间
type SunTwilightTime struct {
	TwilightTime

	Astronomical DawnDuskTime
	Nautical DawnDuskTime
	Civil DawnDuskTime
}

func NewSunTwilightAngle() *SunTwilightAngle {
	return &SunTwilightAngle{
		TwilightAngle: TwilightAngle{
			Rise: ToRadian(-50. / 60.), // 地平线以下 50′
			Set:  ToRadian(-50. / 60.),
			Noon: ToRadian(90.),
			Midnight: ToRadian(-90.),
		},
		Astronomical:  DawnDuskAngle{
			Dawn: ToRadian(-18.),
			Dusk: ToRadian(-18.),
		},
		Nautical:      DawnDuskAngle{
			Dawn: ToRadian(-12.),
			Dusk: ToRadian(-12.),
		},
		Civil:         DawnDuskAngle{
			Dawn: ToRadian(-6.),
			Dusk: ToRadian(-6.),
		},
	}
}

/**
 * 日照时长
 * 日落 - 升日
 */
func (stt *SunTwilightTime) Daylight() JulianDayDelta  {
	return JulianDayDelta(stt.Set - stt.Rise)
}

/**
 * 夜晚时长，因为跨午夜，所以取今日00:00 ~ 23:59的时段
 * (0:00 ~ 天文晨光) + (天文暮光 ~ 23:59)
 */
func (stt *SunTwilightTime) Night() JulianDayDelta {
	return JulianDayDelta(
		stt.Astronomical.Dawn - stt.Astronomical.Dawn.Midnight() +
			(stt.Astronomical.Dusk.AddDays(1).Midnight() - stt.Astronomical.Dusk),
		)
}

/**
 * 太阳
 * 传入本地12点的JdUT(需要是ut),经度(东为负),纬度(南为负)
 */
func SunTwilight(noonJdUT JulianDay, geo *GeographicCoordinates, withRevise bool) *SunTwilightTime  {
	angle := NewSunTwilightAngle()
	time := &SunTwilightTime{
		TwilightTime: TwilightTime{
			Rise:     noonJdUT,
			Set:      noonJdUT,
			Noon:     noonJdUT,
			Midnight: noonJdUT,
		},
		Astronomical: DawnDuskTime{
			Dawn: noonJdUT,
			Dusk: noonJdUT,
		},
		Nautical:     DawnDuskTime{
			Dawn: noonJdUT,
			Dusk: noonJdUT,
		},
		Civil:        DawnDuskTime{
			Dawn: noonJdUT,
			Dusk: noonJdUT,
		},
	}

	astro := NewAstronomy(geo.Longitude, geo.Latitude, noonJdUT)

	var calcHa = func(withRevise bool) (res map[string]float64) {
		// 时角，星星属性，赤道坐标，错误
		ha, planet, equator, err := astro.PlanetHourAngle(swe.Sun, withRevise)

		if err != nil {
			return nil
		}

		haRise := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Rise)
		haSet := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Set)
		haNoon := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Noon)
		haMidnight := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Midnight)
		haAstronomicalDawn  := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Astronomical.Dawn)
		haAstronomicalDusk  := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Astronomical.Dusk)
		haNauticalDawn  := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Nautical.Dawn)
		haNauticalDusk  := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Nautical.Dusk)
		haCivilDawn  := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Civil.Dawn)
		haCivilDusk  := EquatorToHourAngle(equator.Latitude, geo.Latitude, angle.Civil.Dusk)

		if withRevise {
			// 方位角(地平经度)
			azimuth :=  Radian90 - ha
			// 高度角(地平纬度)
			altitude := equator.Latitude

			coord := EclipticEquatorConventor(&Coordinates{
				Longitude: azimuth,
				Latitude:  altitude,
			}, Radian90 - geo.Latitude)

			azimuth = coord.Longitude
			altitude = coord.Latitude

			azimuth = RadianMod(Radian90 - azimuth)

			ha = azimuth
			haRise = altitude

			if haRise > 0 {
				// 修正大气折射
				haRise += AstronomicalRefraction2(haRise)
			}

			// 直接在地平坐标中视差修正(这里把地球看为球形，精度比 Parallax 秒差一些)
			haRise -= 8.794 / DegreeSecondsPerRadian / planet.Distance * math.Cos(haRise)


		}

		res = map[string]float64{}

		res["ha"] = ha
	}





}