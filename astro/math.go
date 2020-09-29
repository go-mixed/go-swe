package astro

import (
	"fmt"
	"math"
)

// 每弧度的角秒数，此处的3600是因为 1° == 3600″
// 所以可读的表达式为 角° = 1 弧度 * 180 / Pi; 角″ = 角° * 3600
const DegreeSecondsPerRadian float64 = 1 * 180 * 3600. / math.Pi
const Radian90 float64 = math.Pi / 2.
const Radian180 float64 = math.Pi
const Radian360 float64 = math.Pi * 2

// ToDegree 弧度转换为角度
func ToDegree(r float64) float64 {
	return r * 180. / math.Pi
}

// 对超过0~2PI的角度转为0~2PI,也就是超过360度的角度转换为360度以内
func RadianMod(r float64) float64 {

	_r := math.Mod(r, Radian360)

	if _r < 0 {
		return _r + Radian360
	}

	return _r
}

// ToRadian 角度转换为弧度
func ToRadian(d float64) float64 {
	return d * math.Pi / 180.
}

// String 角度转换为书面表达的 ° ′ ″
func DegreeToString(d float64) string {
	var deg, min int
	var sec float64
	var absDegree = math.Abs(d)

	deg = int(absDegree)
	min = int(math.Mod(absDegree * 60, 60))
	sec = math.Mod(absDegree * 3600, 60)

	if deg >= 0 {
		return fmt.Sprintf("%d°%d′%.4f″", deg, min, sec)
	}

	return fmt.Sprintf("-%d°%d′%.4f″", deg, min, sec)
}
