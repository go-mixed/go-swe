package astro

import (
	"go-swe/swe"
	"math"
)

type Astronomy struct {
	// 地理位置
	Geo *GeographicCoordinates
	// 儒略日(世界时)
	JdUt JulianDay
	// JD 和 ET 的 deltaT
	DeltaT float64
	// Swe的实例
	Swe swe.SweInterface
}

type EclipticProperties struct {
	// 真黄道真倾角(含章动)，即黄赤交角，转轴倾角
	// true obliquity of the Ecliptic (includes nutation)
	TrueObliquity float64
	// 黄道平均倾角
	// 比如地球是：23°26′20.512″
	// mean obliquity of the Ecliptic
	MeanObliquity float64
	// 黄经章动
	// nutation in longitude
	NutationInLongitude float64
	// 倾角章动
	// nutation in obliquity
	NutationInObliquity float64
}

type PlanetProperties struct {
	// 黄道坐标
	Ecliptic *EclipticCoordinates
	// 距离
	Distance float64
	// 经度的速度
	SpeedInLongitude float64
	// 纬度的速度
	SpeedInLatitude float64
	// 距离的速度
	SpeedInDistance float64
	// 时角
	HourAngle HourAngle
	// 赤道坐标
	Equatorial *EquatorialCoordinates
	// 地平坐标
	Horizontal *HorizontalCoordinates
}

func NewAstronomy(geo *GeographicCoordinates, jdUt JulianDay) *Astronomy {

	_swe := swe.NewSwe()
	return &Astronomy{
		Geo:    geo,
		JdUt:   jdUt,
		DeltaT: DeltaT(jdUt),
		Swe:    _swe,
	}
}

// 儒略日(天文时)
func (astro *Astronomy) JdEt() JulianDay {
	return astro.JdUt.ToEphemerisTime(astro.DeltaT)
}

func (astro *Astronomy) simpleCalcFlags() *swe.CalcFlags {
	var iFlag int32 = swe.FlagEphSwiss | swe.FlagRadians | swe.FlagSpeed
	return &swe.CalcFlags{
		Flags: iFlag,
	}
}

/**
 * 黄道属性，包含倾角、章动
 */
func (astro *Astronomy) EclipticProperties() (*EclipticProperties, error) {

	// 黄道章动
	res, _, err := astro.Swe.Calc(float64(astro.JdEt()), swe.EclNut, astro.simpleCalcFlags())

	if err != nil {
		return nil, err
	}

	return &EclipticProperties{
		TrueObliquity:       res[0],
		MeanObliquity:       res[1],
		NutationInLongitude: res[2],
		NutationInObliquity: res[3],
	}, nil
}

/**
 * 天体的基础属性，包括黄道经纬，距离，黄道经纬速度，距离速度
 */
func (astro *Astronomy) PlanetBaseProperties(planetId swe.Planet) (*PlanetProperties, error) {
	// 天体的基本属性：黄道坐标、距离等参数
	res, _, err := astro.Swe.Calc(float64(astro.JdEt()), planetId, astro.simpleCalcFlags())

	if err != nil {
		return nil, err
	}

	return &PlanetProperties{
		Ecliptic: &EclipticCoordinates{
			Longitude: res[0],
			Latitude:  res[1],
		},
		Distance:         0,
		SpeedInLongitude: res[3],
		SpeedInLatitude:  res[4],
		SpeedInDistance:  res[5],
	}, nil
}

/**
 * 返回天体的所有属性，除了基本属性外，还包含时角、赤道坐标、地平坐标
 * 关于HA的计算，因为章动同时影响恒星时和赤道坐标，所以不计算章动。
 * withRevise 是否修正，包含使用真黄道倾角、修正大气折射、修正地平坐标中视差
 */
func (astro *Astronomy) PlanetProperties(planetId swe.Planet, withRevise bool) (planet *PlanetProperties, err error) {
	// 当前黄道倾角、章动等参数
	ecliptic, err := astro.EclipticProperties()
	if err != nil {
		return
	}

	// 天体的基本属性：黄道坐标、距离等参数
	planet, err = astro.PlanetBaseProperties(planetId)
	if err != nil {
		return
	}

	// 修正光行差 20.5″
	if withRevise {
		planet.Ecliptic.Longitude -= 20.5 / DegreeSecondsPerRadian
	}

	// 黄道坐标 -> 赤道坐标
	equatorial := EclipticToEquatorial(planet.Ecliptic, IfThenElse(withRevise, ecliptic.TrueObliquity, ecliptic.MeanObliquity).(float64))

	// 不太精确的恒星时
	sidTime := GreenwichMeridianSiderealTime(astro.JdUt, astro.DeltaT)
	// 修正恒星时
	if withRevise {
		sidTime += ecliptic.NutationInLongitude * math.Cos(ecliptic.TrueObliquity)
	}

	/*
		https://en.wikipedia.org/wiki/Hour_angle
		如果θ是本地恒星时，θo是格林尼治恒星时，λ是观者站经度（从格林尼治向西为负，东为正），α是赤经，那么本地时角H计算如下：
		H = θ - α 或 H =θo + λ - α
		如果α含章动效果，那么H也含章动（见11章）。
	*/

	// 时角，转换到-180°~180°内
	hourAngle := HourAngle(RadianMod180(sidTime + astro.Geo.Longitude - equatorial.RightAscension))

	var horizontal *HorizontalCoordinates

	if withRevise {
		_horizontal := EclipticEquatorialConverter(&GeographicCoordinates{
			Longitude: Radian90 - float64(hourAngle),
			Latitude:  equatorial.Declination,
		}, Radian90-astro.Geo.Latitude)
		_horizontal.Longitude = RadianMod360(Radian90 - _horizontal.Longitude)

		//修正大气折射
		if _horizontal.Latitude > 0 {
			_horizontal.Latitude += AstronomicalRefraction2(_horizontal.Latitude)
		}
		// 直接在地平坐标中视差修正(这里把地球看为球形,精度比 Parallax 秒差一些)
		_horizontal.Latitude -= 8.794 / DegreeSecondsPerRadian / planet.Distance * math.Cos(_horizontal.Latitude)
		horizontal = _horizontal.AsHorizontalCoordinates()
	} else {
		horizontal = EquatorialToHorizontal(hourAngle, equatorial.Declination, astro.Geo.Latitude)
	}

	planet.HourAngle = hourAngle
	planet.Equatorial = equatorial
	planet.Horizontal = horizontal

	return
}
