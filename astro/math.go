package astro

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// 每弧度的角秒数，此处的3600是因为 1° == 3600″
// 所以可读的表达式为 角° = 1 弧度 * 180 / Pi; 角″ = 角° * 3600
const DegreeSecondsPerRadian float64 = 1 * 180 * 3600. / math.Pi
const Radian90 float64 = math.Pi / 2.
const Radian180 float64 = math.Pi
const Radian360 float64 = math.Pi * 2

const (
	DegreesString = '\u00B0'
	MinutesString = '\''
	SecondsString = '"'
	PointString   = '.'
)

// ToDegree 弧度转换为角度
func ToDegree(r float64) float64 {
	return r * 180. / math.Pi
}

// 对超过0~2PI的角度转为0~2PI,也就是超过360°转换为360°以内
func RadianMod360(r float64) float64 {
	_r := math.Mod(r, Radian360)

	if _r < 0 {
		return _r + Radian360
	}
	return _r
}

// 对超过-PI到PI的角度转为-PI到PI, 即转化到-180°~180°
func RadianMod180(r float64) float64 {
	_r := math.Mod(r, Radian360)
	if _r < -Radian180 {
		return _r + Radian360
	} else if _r > Radian180 {
		return _r - Radian360
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
	min = int(math.Mod(absDegree*60, 60))
	sec = math.Mod(absDegree*3600, 60)

	return fmt.Sprintf("%d%c%d%c%.4f%c", deg, DegreesString, min, MinutesString, sec, SecondsString)
}

// ParseDMS parses a coordinate in degrees, minutes, seconds.
// - e.g. 33° 23' 22"
func StringToDegree(s string) (float64, error) {
	degrees := 0
	minutes := 0
	seconds := 0.0
	// Whether a number has finished parsing (i.e whitespace after it)
	endNumber := false
	// Temporary parse buffer.
	var tmpBytes []byte
	var err error

	s = strings.ReplaceAll(s, "′", string(MinutesString))
	s = strings.ReplaceAll(s, "″", string(SecondsString))

	for i, r := range s {
		if unicode.IsNumber(r) || r == PointString {
			if !endNumber {
				tmpBytes = append(tmpBytes, s[i])
			} else {
				return 0, errors.New("parse error (no delimiter)")
			}
		} else if unicode.IsSpace(r) && len(tmpBytes) > 0 {
			endNumber = true
		} else if r == DegreesString {
			if degrees, err = strconv.Atoi(string(tmpBytes)); err != nil {
				return 0, errors.New("parse error (degrees)")
			}
			tmpBytes = tmpBytes[:0]
			endNumber = false
		} else if s[i] == MinutesString {
			if minutes, err = strconv.Atoi(string(tmpBytes)); err != nil {
				return 0, errors.New("parse error (minutes)")
			}
			tmpBytes = tmpBytes[:0]
			endNumber = false
		} else if s[i] == SecondsString {
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

// 临界余数(a与最近的整倍数b相差的距离)
func Mod2(a, b float64) float64 {
	c := a / b
	c -= math.Floor(c)

	if c > .5 {
		c -= 1
	}

	return c * b
}
