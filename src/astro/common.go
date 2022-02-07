package astro

import (
	"math"
)

func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

// ShadowLength 1m 的物体影子长度
func ShadowLength(altitude float64) float64 {
	//tan(90) 错误,大于90度也无效
	if altitude <= 0 || altitude > Radian90 {
		return 0
	}

	// 已知直角边a = 1，和直角边b与斜边的角度∠A（即太阳的高度角），求b;
	//        /|
	//       / |
	//      /  | a = 1 物体高1m
	// ∠A /___|
	//      b = 影子长度
	// a / b = tan(A);
	// 则 b = a / tan(A),或者b = a * tan(90-A);
	return math.Tan(Radian90 - altitude)
}

// AstronomicalRefraction 大气折射修正
//	ho: 视高度
func AstronomicalRefraction(ho float64) float64 {
	return -0.0002909 / math.Tan(ho+0.002227/(ho+0.07679))
}

// AstronomicalRefraction2 大气折射修正
//	h: 真高度
func AstronomicalRefraction2(h float64) float64 {
	return 0.0002967 / math.Tan(h+0.003138/(h+0.08919))
}

// Parallax 视差修正
//	equator: 赤道坐标
//	distance: 距离，以AU为单位
//	hourAngle: 时角
//	latitude: 观察者纬度
//	height: 观察者海拔
func Parallax(equatorial EquatorialCoordinates, distance, hourAngle, latitude, height float64) *GeographicCoordinates {
	//赤道地平视差
	sinP := 8.794 / DegreeSecondsPerRadian / distance
	ba := 0.99664719
	u := math.Atan(ba * math.Tan(latitude))
	sinD := -sinP * (math.Sin(u)*ba + height*math.Sin(latitude)/6378.14)
	cosD := -sinP * (math.Cos(u) + height*math.Cos(latitude)/6378.14)
	sinH := math.Sin(hourAngle)
	cosH := math.Cos(hourAngle)
	sinW := math.Sin(equatorial.Declination) // 赤道纬度
	cosW := math.Cos(equatorial.Declination)

	a := math.Atan2(cosD*sinH, cosW+cosD*cosH)

	return &GeographicCoordinates{
		Longitude: RadiansMod360(equatorial.RightAscension + a),
		Latitude:  math.Atan((sinW + sinD) / (cosW + cosD*cosH) * math.Cos(a)),
	}
}
