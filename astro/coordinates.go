package astro

import "math"

const (
	NorthString = "N"
	SouthString = "S"
	EastString  = "E"
	WestString  = "W"
)

/**
 * 时角 τ tau
 * https://zh.wikipedia.org/wiki/File:HourAngle_Observer_en.png
 * 该天体的赤经与当地的恒星时的差值，从子午线沿天赤道，到目标所在过地球极点的球面大圆的角距离
 * 在子午线的东边则为负时角，在子午线的西边则为正时角，或者向西为正的360度，时角与经度的换算方法为24h = 360°
 */
type HourAngle float64

type Coordinates struct {
	x float64
	y float64
	z float64
}

// 地理坐标系
type GeographicCoordinates struct {
	/**
	 * 经度 λ lambda
	 * 西 -180° ~ 本初子午线 0 ~ 东 180°
	 */
	Longitude float64
	/**
	 * 纬度 φ phi
	 * 南 -90° ~ 赤道 0 ~ 北 90°
	 */
	Latitude float64
}

// 地平坐标系
// https://zh.wikipedia.org/wiki/File:Azimuth-Altitude_schematic.svg
type HorizontalCoordinates struct {
	/**
	 * 方位角，地平经度 Az A
	 * 0 ~ 180° 子午线东, 北 0° ~ 东 90° ~ 南 180°
	 * 180° ~ 360° 子午线西, 南 180° ~ 西 270° ~ 北 360°
	 * 沿着地平线测量的角度（由正北方为起点向东方测量）
	 */
	Azimuth float64
	/**
	 * 高度角，仰角，地平纬度 Alt a
	 * 有时就称为高度或海拔标高（elevation, geometric height）
	 * 天底 -90° ~ 升 0° ~ 天顶 90° ~ 降 0°
	 */
	Altitude float64
}

// 赤道坐标系
// https://en.wikipedia.org/wiki/File:Ra_and_dec_demo_animation_small.gif
type EquatorialCoordinates struct {
	/**
	 * 赤经 α alpha RA
	 * 春分点 0 ~ 秋分点 180° ~ 360°
	 * 西 -> 东 逆时针
	 * 春分点：太阳在3月下旬运行至北天球时所通过的点，也是地球的升交点
	 */
	RightAscension float64
	/**
	 * 赤纬 δ delta
	 * 南 -90° ~ 赤道 0 ~ 北 90°
	 */
	Declination float64
}

// 黄道坐标系
type EclipticCoordinates struct {
	/**
	 * 黄道经度 λ lambda
	 * 春分点 0 ~ 秋分点 180° ~ 360°
	 * 西 -> 东 逆时针
	 */
	Longitude float64
	/**
	 * 黄道纬度 β beta
	 * 南 -90° ~ 黄道面 0 ~ 北 90°
	 */
	Latitude float64
}

// 表示为地理坐标
func (h *HorizontalCoordinates) AsGeographicCoordinates() *GeographicCoordinates {
	return &GeographicCoordinates{
		Longitude: h.Azimuth,
		Latitude:  h.Altitude,
	}
}

// 表示为地理坐标
func (e *EquatorialCoordinates) AsGeographicCoordinates() *GeographicCoordinates {
	return &GeographicCoordinates{
		Longitude: e.RightAscension,
		Latitude:  e.Declination,
	}
}

// 表示为地理坐标
func (e *EclipticCoordinates) AsGeographicCoordinates() *GeographicCoordinates {
	return &GeographicCoordinates{
		Longitude: e.Longitude,
		Latitude:  e.Latitude,
	}
}

// 表示为地平坐标
func (g *GeographicCoordinates) AsHorizontalCoordinates() *HorizontalCoordinates {
	return &HorizontalCoordinates{
		Azimuth:  g.Longitude,
		Altitude: g.Latitude,
	}
}

// 表示为赤道坐标
func (g *GeographicCoordinates) AsEquatorialCoordinates() *EquatorialCoordinates {
	return &EquatorialCoordinates{
		RightAscension: g.Longitude,
		Declination:    g.Latitude,
	}
}

// 表示为黄道坐标
func (g *GeographicCoordinates) AsEclipticCoordinates() *EclipticCoordinates {
	return &EclipticCoordinates{
		Longitude: g.Longitude,
		Latitude:  g.Latitude,
	}
}

/**
 * 赤道坐标 -> 地平坐标
 * https://zh.wikipedia.org/wiki/%E5%9C%B0%E5%B9%B3%E5%9D%90%E6%A8%99%E7%B3%BB
 * HourAngle: 时角
 * declination: 赤纬
 * latitude: 观察者纬度
 */
func EquatorialToHorizontal(hourAngle HourAngle, declination, latitude float64) *HorizontalCoordinates {
	/*
		纬度用φ表示，赤纬用δ表示，地方时(时角)以H表示：
		sin(Alt) ＝ sin(φ) * sin(δ) ＋ cos(φ) * cos(δ) * cos(H)
		tan(Az) = cos(δ) * sin(H) / ( -cos(φ) * sin(δ) + sin(φ) * cos(δ) * cos(H) )

		完整表达式是：

		cos(Az) * cos(Alt) ＝ -cos(φ) * sin(δ) ＋ sin(φ) * cos(δ) * cos(H)
		sin(Az) * cos(Alt) = cos(δ) * sin(H)

	*/
	_hourAngle := float64(hourAngle)
	altitude := math.Asin(
		math.Sin(latitude)*math.Sin(declination) + math.Cos(latitude)*math.Cos(declination)*math.Cos(_hourAngle),
	)
	azimuth := math.Atan2(
		math.Cos(declination)*math.Sin(_hourAngle),
		-math.Cos(latitude)*math.Sin(declination)+math.Sin(latitude)+math.Cos(declination)+math.Cos(_hourAngle),
	)

	return &HorizontalCoordinates{
		Azimuth:  azimuth,
		Altitude: altitude,
	}
}

/**
 * 根据高度角和当前的赤纬、地理纬度 计算 时角
 * 一般用于计算星星的升天、中天等值。但是此时得到的时角是不准确的，因为高度角变化时，星星的赤纬也在变化，需要多次求解
 * declination: 赤纬
 * latitude: 观察者纬度
 * altitude: 高度角
 */
func AltitudeToHourAngle(declination, latitude, altitude float64) HourAngle {
	// 纬度用φ表示，赤纬用δ表示，地方时(时角)以H表示：
	// sin(Alt) = sin(φ) * sin(δ) + cos(φ) * cos(δ) * cos(H)
	ha := (math.Sin(altitude) - math.Sin(latitude)*math.Sin(declination)) / math.Cos(declination) / math.Cos(latitude)

	// > 180°
	if math.Abs(ha) > 1 {
		return HourAngle(Radian180)
	}
	return HourAngle(math.Acos(ha))
}

/**
 * 黄道坐标 <-> 赤道坐标 互转 也就是 球面坐标旋转
 * 如果是 赤->黄 eclipticObliquity 设置为负数
 * eclipticObliquity: 黄赤交角，比如地球是：23°26′20.512″
 * https://zh.wikipedia.org/wiki/%E9%BB%83%E9%81%93%E5%9D%90%E6%A8%99%E7%B3%BB
 */
func EclipticEquatorialConverter(coordinates *GeographicCoordinates, eclipticObliquity float64) *GeographicCoordinates {
	/*
		λ和β代表黄经和黄纬
		α和δ代表赤经和赤纬
		ε = eclipticObliquity ≈ 23°26′20.512″

		赤 -> 黄
		tan(λ) = ( sin(α)*cos(ε) + tan(δ)*sin(ε) ) / cos(α)
		sin(β) = sin(δ)*cos(ε) - cos(δ)*sin(ε)*sin(α)

		下面2个表达式可以简化为上面的tan
		cos(λ)*cos(β) = cos(α) * cos(δ)
		sin(λ)*cos(β) = sin(ε) * sin(δ) + cos(δ)*cos(ε)*sin(α)
		解：根据 tan = sin / cos
		左项：sin(λ) * cos(β) / cos(λ) * cos(β)
		   = tan(λ)
		右项：(sin(ε) * sin(δ) + cos(δ) * cos(ε) * sin(α)) / (cos(α) * cos(δ))
		   = (sin(ε) * tan(δ) +           cos(ε) * sin(α)) / cos(α)

		黄 -> 赤 算法
		tan(α) = ( sin(λ) * cos(ε) - tan(λ) * sin(ε) ) / cos(λ)
		sin(δ) = sin(ε) * sin(λ) * cos(β) + cos(ε) * sin(β)

		同理，tan的表达式由以下两个表达式简化而成
		cos(α) * cos(δ) = cos(λ) * cos(β)
		sin(α) * cos(δ) = cos(ε) * sin(λ) * cos(β) - sin(ε) * sin(β)


		如果将上面4个公式中的右边表达式经纬替换成x y，可见 这4个公式 仅仅只有 +、- 的区别
		下文是 黄 -> 赤 的算法，如果将 eclipticObliquity 设置为负，sinE就变成了负数，完成了表达式的+、-转换，就变成了 赤 -> 黄 的算法
	*/

	sinE := math.Sin(eclipticObliquity) // ε
	cosE := math.Cos(eclipticObliquity)
	sinJ := math.Sin(coordinates.Longitude) // 经度
	cosJ := math.Cos(coordinates.Longitude)
	sinW := math.Sin(coordinates.Latitude) // 纬度
	cosW := math.Cos(coordinates.Latitude)
	tanW := math.Tan(coordinates.Latitude)

	/*
		黄 -> 赤 公式
	*/
	longitude := math.Atan2(sinJ*cosE-tanW*sinE, cosJ)
	latitude := math.Asin(sinW*cosE + cosW*sinE*sinJ)

	longitude = RadianMod360(longitude) //保证在360度内

	return &GeographicCoordinates{
		Longitude: longitude,
		Latitude:  latitude,
	}

}

/**
 * 黄道坐标 -> 赤道坐标
 * eclipticObliquity: 黄赤交角，比如地球是：23°26′20.512″
 */
func EclipticToEquatorial(coordinates *EclipticCoordinates, eclipticObliquity float64) *EquatorialCoordinates {
	return EclipticEquatorialConverter(coordinates.AsGeographicCoordinates(), eclipticObliquity).AsEquatorialCoordinates()
}

/**
 * 赤道坐标 -> 黄道坐标
 * eclipticObliquity: 黄赤交角，比如地球是：23°26′20.512″
 */
func EquatorialToEcliptic(coordinates *EquatorialCoordinates, eclipticObliquity float64) *EclipticCoordinates {
	return EclipticEquatorialConverter(coordinates.AsGeographicCoordinates(), -eclipticObliquity).AsEclipticCoordinates()
}
