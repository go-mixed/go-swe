package astro

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

const (
	// 每弧度的角秒数，此处的3600是因为 1° == 3600″
	// 所以可读的表达式为 角° = 1 弧度 * 180 / Pi; 角″ = 角° * 3600
	DegreeSecondsPerRadian float64 = 1 * 180 * 3600. / math.Pi
	// 90°的弧度
	Radian90 float64 = math.Pi / 2.
	// 180°的弧度
	Radian180 float64 = math.Pi
	// 360°的弧度
	Radian360 float64 = math.Pi * 2
)

const (
	DegreesRune = '\u00B0'
	MinutesRune = '\''
	SecondsRune = '"'
	PointRune   = '.'
)

// ToDegrees 弧度转换为角度
func ToDegrees(r float64) float64 {
	return r * 180. / math.Pi
}

// 对超过0~2PI的角度转为0~2PI,也就是超过360°转换为360°以内
func RadiansMod360(r float64) float64 {
	_r := math.Mod(r, Radian360)

	if _r < 0 {
		return _r + Radian360
	}
	return _r
}

// 对超过-PI到PI的角度转为-PI到PI, 即转化到-180°~180°
func RadiansMod180(r float64) float64 {
	_r := math.Mod(r, Radian360)
	if _r > Radian180 {
		_r -= Radian360
	}
	if _r < -Radian180 {
		_r += Radian360
	}
	return _r
}

// ToRadians 角度转换为弧度
func ToRadians(d float64) float64 {
	return d * math.Pi / 180.
}

// String 角度转换为书面表达的 ° ′ ″
func DegreesToString(d float64) string {
	var deg, min int
	var sec float64
	var absDegree = math.Abs(d)

	deg = int(absDegree)
	min = int(math.Mod(absDegree*60, 60))
	sec = math.Mod(absDegree*3600, 60)

	return fmt.Sprintf("%d%c%d%c%.4f%c", deg, DegreesRune, min, MinutesRune, sec, SecondsRune)
}

// ParseDMS parses a coordinate in degrees, minutes, seconds.
// - e.g. 33° 23' 22"ma
func StringToDegrees(s string) (float64, error) {
	degrees := 0
	minutes := 0
	seconds := 0.0
	// Whether a number has finished parsing (i.e whitespace after it)
	endNumber := false
	// Temporary parse buffer.
	var tmpBytes []byte
	var err error

	s = strings.ReplaceAll(s, "′", string(MinutesRune))
	s = strings.ReplaceAll(s, "″", string(SecondsRune))

	for i, r := range s {
		if unicode.IsNumber(r) || r == PointRune {
			if !endNumber {
				tmpBytes = append(tmpBytes, s[i])
			} else {
				return 0, errors.New("parse error (no delimiter)")
			}
		} else if unicode.IsSpace(r) && len(tmpBytes) > 0 {
			endNumber = true
		} else if r == DegreesRune {
			if degrees, err = strconv.Atoi(string(tmpBytes)); err != nil {
				return 0, errors.New("parse error (degrees)")
			}
			tmpBytes = tmpBytes[:0]
			endNumber = false
		} else if s[i] == MinutesRune {
			if minutes, err = strconv.Atoi(string(tmpBytes)); err != nil {
				return 0, errors.New("parse error (minutes)")
			}
			tmpBytes = tmpBytes[:0]
			endNumber = false
		} else if s[i] == SecondsRune {
			if seconds, err = strconv.ParseFloat(string(tmpBytes), 64); err != nil {
				return 0, errors.New("parse error (seconds)")
			}
			tmpBytes = tmpBytes[:0]
			endNumber = false
		} else if unicode.IsSpace(r) && len(tmpBytes) == 0 {
			continue
		} else {
			return 0, fmt.Errorf("parse error (unknown symbol [%d])", s[i])
		}
	}
	val := float64(degrees) + (float64(minutes) / 60.0) + (seconds / 60.0 / 60.0)
	return val, nil
}

/**
 * 时间 转 角度
 * 赤经、恒星时、时角 有一种小时表示的单位，也就是24h = 360°
 */
func HoursToDegrees(hours float64) float64 {
	return hours * 15
}

// 临界余数(v与最近的整倍数n相差的距离)
func Mod2(v, n float64) float64 {
	c := v / n
	c -= math.Floor(c) // 取两者余数的小数部分

	if c > .5 {
		c -= 1
	}

	return c * n // 得到余数
}

// 取一下倍数
// 比如: NextMultiples(31, 15) = 45  30的下一个符合15的倍数是45
// NextMultiples(30, 15) = 30
func NextMultiples(v, n float64) float64 {
	c := v / n
	i := math.Floor(c) // 取两者的整数部分
	d := c - i         // 取两者余数的小数部分

	// 有余数
	if !FloatEqual(d, 0, 9) {
		i++
	}

	return n * i
}

// 求2个浮点数是否相等, scale 是对比小数位数
func FloatEqual(x, y float64, scale int) bool {
	return math.Abs(x-y) < (1 / math.Pow(10, float64(scale)))
}
