package astro

import (
	"go-swe/swe"
	"math"
)

const (
	EquatorialRadius           = 6378.1366                                                                // 地球赤道半径(千米)
	MeanEquatorialRadius       = 0.99834 * EquatorialRadius                                               // 平均半径
	PolarEquatorialRatio       = 0.99664719                                                               // 地球极赤半径比
	PolarEquatorialRatioSquare = PolarEquatorialRatio * PolarEquatorialRatio                              // 地球极赤半径比的平方
	AU                         = 1.49597870691e8                                                          // 天文单位长度(千米)
	SinSolarParallax           = EquatorialRadius / AU                                                    // sin(太阳视差)
	LightVelocity              = 299792.458                                                               // 光速(行米/秒)
	LightTimePerAU             = AU / LightVelocity / 86400 / 36525                                       // 每天文单位的光行时间(儒略世纪)
	LunarEarthRatio            = 0.2725076                                                                // 月亮与地球的半径比(用于半影计算)
	LunarEarthRatio2           = 0.2722810                                                                // 月亮与地球的半径比(用于本影计算)
	SolarEarthRatio            = 109.1222                                                                 // 太阳与地球的半径比(对应959.64)
	ApparentLunarRadius        = LunarEarthRatio * EquatorialRadius * 1.0000036 * DegreeSecondsPerRadian  // 用于月亮视半径计算
	ApparentLunarRadius2       = LunarEarthRatio2 * EquatorialRadius * 1.0000036 * DegreeSecondsPerRadian // 用于月亮视半径计算
	ApparentSolarRadius        = 959.64                                                                   // 用于太阳视半径计算

	MeanLunarDays = 29.530587981 // 平均农历月的日数
	MeanSolarDays = 365.2425     // 平均太阳年的日数
)

var SolarParallax = math.Asin(SinSolarParallax)                                      // 太阳视差
var PlanetaryRendezvousPeriod = [...]float64{116, 584, 780, 399, 378, 370, 367, 367} //行星会合周期

type Astronomy struct {
	// Swe的实例
	Swe swe.SweInterface
}

type EclipticProperties struct {
	// 真黄赤交角(含章动)，即转轴倾角
	// true obliquity of the Ecliptic (includes nutation)
	TrueObliquity float64
	// 平黄赤交角
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
	PlanetId swe.Planet
	// 黄道坐标
	Ecliptic *EclipticCoordinates
	// 距离 单位是 AU, 转化为千米 则是 Distance * AU
	Distance float64
	// 黄经的速度 单位是 弧度/天
	SpeedInLongitude float64
	// 黄纬的速度 单位是 弧度/天
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

// 返回千米为单位的距离
func (p *PlanetProperties) DistanceAsKilometer() float64 {
	return p.Distance * AU
}

func NewAstronomy() *Astronomy {
	_swe := swe.NewSwe()
	return &Astronomy{
		Swe: _swe,
	}
}

func (astro *Astronomy) DeltaTEx(jdUT JulianDay) float64 {
	deltaT, _ := astro.Swe.DeltaTEx(float64(jdUT), swe.FlagEphSwiss)
	return deltaT
}

func (astro *Astronomy) DeltaT(jdUT JulianDay) float64 {
	return astro.Swe.DeltaT(float64(jdUT))
}

func (astro *Astronomy) simpleCalcFlags(deltaT float64) *swe.CalcFlags {
	var iFlag int32 = swe.FlagEphSwiss | swe.FlagRadians | swe.FlagSpeed
	return &swe.CalcFlags{
		Flags:  iFlag,
		DeltaT: &deltaT,
	}
}

/**
 * 黄道属性，包含倾角、章动
 */
func (astro *Astronomy) EclipticProperties(jdET *EphemerisTime) (*EclipticProperties, error) {

	// 黄道章动
	res, _, err := astro.Swe.Calc(jdET.Value(), swe.EclNut, astro.simpleCalcFlags(jdET.DeltaT))

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
 * 天体属性，包括黄道经纬，距离，黄道经纬速度，距离速度
 */
func (astro *Astronomy) PlanetProperties(planetId swe.Planet, jdET *EphemerisTime) (*PlanetProperties, error) {
	// 天体的基本属性：黄道坐标、距离等参数
	res, _, err := astro.Swe.Calc(jdET.Value(), planetId, astro.simpleCalcFlags(jdET.DeltaT))

	if err != nil {
		return nil, err
	}

	return &PlanetProperties{
		PlanetId: planetId,
		Ecliptic: &EclipticCoordinates{
			Longitude: res[0],
			Latitude:  res[1],
		},
		Distance:         res[2],
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
func (astro *Astronomy) PlanetPropertiesWithObserver(
	planetId swe.Planet,
	jdET *EphemerisTime,
	geo *GeographicCoordinates,
	withRevise bool) (planet *PlanetProperties, err error) {
	// 当前黄道倾角、章动等参数
	ecliptic, err := astro.EclipticProperties(jdET)
	if err != nil {
		return
	}

	// 天体的基本属性：黄道坐标、距离等参数
	planet, err = astro.PlanetProperties(planetId, jdET)
	if err != nil {
		return
	}

	// 修正光行差 20.5″
	if withRevise {
		planet.Ecliptic.Longitude -= 20.5 / DegreeSecondsPerRadian
	}

	// 黄道坐标 -> 赤道坐标
	equatorial := EclipticToEquatorial(planet.Ecliptic, IfThenElse(withRevise, ecliptic.TrueObliquity, ecliptic.MeanObliquity).(float64))

	var sidTime float64
	/* 快速计算sidreal time
	// 不太精确的恒星时
	sidTime = GreenwichMeridianSiderealTime(astro.JdUt, astro.DeltaT)
	// 修正恒星时
	if withRevise {
		sidTime += ecliptic.NutationInLongitude * math.Cos(ecliptic.TrueObliquity)
	}
	*/

	// 使用swe计算恒星时
	if withRevise {
		sidTime, _ = astro.Swe.SidTime0(float64(jdET.JdUT), ToDegrees(ecliptic.TrueObliquity), ToDegrees(ecliptic.NutationInLongitude), &swe.SidTimeFlags{DeltaT: &jdET.DeltaT})
	} else {
		sidTime, _ = astro.Swe.SidTime(float64(jdET.JdUT), &swe.SidTimeFlags{DeltaT: &jdET.DeltaT})
	}
	// swe 返回的单位是角度的时间表示方式. * 15 则为度
	sidTime = ToRadians(HoursToDegrees(sidTime))

	/*
		https://en.wikipedia.org/wiki/Hour_angle
		如果θ是本地恒星时，θo是格林尼治恒星时，λ是观者站经度（从格林尼治向西为负，东为正），α是赤经，那么本地时角H计算如下：
		H = θ - α 或 H =θo + λ - α
		如果α含章动效果，那么H也含章动（见11章）。
	*/

	// 时角，转换到-180°~180°内
	hourAngle := HourAngle(RadiansMod180(sidTime + geo.Longitude - equatorial.RightAscension))

	var horizontal *HorizontalCoordinates

	if withRevise {
		_horizontal := EclipticEquatorialConverter(&GeographicCoordinates{
			Longitude: Radian90 - float64(hourAngle),
			Latitude:  equatorial.Declination,
		}, Radian90-geo.Latitude)
		_horizontal.Longitude = RadiansMod360(Radian90 - _horizontal.Longitude)

		//修正大气折射
		if _horizontal.Latitude > 0 {
			_horizontal.Latitude += AstronomicalRefraction2(_horizontal.Latitude)
		}
		// 直接在地平坐标中视差修正(这里把地球看为球形,精度比 Parallax 秒差一些)
		_horizontal.Latitude -= 8.794 / DegreeSecondsPerRadian / planet.Distance * math.Cos(_horizontal.Latitude)
		horizontal = _horizontal.AsHorizontalCoordinates()
	} else {
		horizontal = EquatorialToHorizontal(hourAngle, equatorial.Declination, geo.Latitude)
	}

	planet.HourAngle = hourAngle
	planet.Equatorial = equatorial
	planet.Horizontal = horizontal

	return
}
