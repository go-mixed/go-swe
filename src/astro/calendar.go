package astro

import (
	"fmt"
	"math"
)

type LunarMonth struct {
	// 朔时间
	JdUT JulianDay `json:"jd_ut"`
	// 第几个月
	Index int `json:"index"`
	// 该月有多少天
	Days int `json:"days"`
	// 是否是闰月
	Leap bool `json:"leap"`
}

// LunarMonths 阴历
//
// 夏正：以冬至日必须在子月（寅正十一月），上个冬至月（寅正十一月）到下个冬至月如有12个月就不置闰，如有13个月就要置闰，以上个冬至月之后第一个无中气的月份为闰月
func (astro *Astronomy) LunarMonths(year int) ([]*LunarMonth, error) {
	// 去年/今年冬至日
	winterSolstices, err := astro.SolarEclipticLongitudesToTimes(DateToJulianDay(year-1, 12, 1, 0, 0, 0), []float64{ToRadians(270), ToRadians(270)})
	if err != nil {
		return nil, fmt.Errorf("LunarMonths WinterSolstices: %w", err)
	}
	// 去年冬至之前的第一个朔日
	lastNewMoon, err := astro.LastNewMoons(winterSolstices[0])
	if err != nil {
		return nil, fmt.Errorf("LunarMonths LastNewMoon: %w", err)
	}
	// 今年冬至日的第一个朔日
	nextNewMoon, err := astro.LastNewMoons(winterSolstices[1])
	if err != nil {
		return nil, fmt.Errorf("LunarMonths NextNewMoon: %w", err)
	}
	// 从去年冬至日之前的朔日, 推14个月
	newMoons, err := astro.NextNewMoons(lastNewMoon, 14)
	if err != nil {
		return nil, fmt.Errorf("LunarMonths NewMoons: %w", err)
	}
	// 两个冬至之间的节气
	solarTerms, err := astro.SolarTermsRange(winterSolstices[0], winterSolstices[1])
	if err != nil {
		return nil, fmt.Errorf("LunarMonths SolarTerms: %w", err)
	}

	// 该年有多少个月
	has13 := FloatEqual(float64(newMoons[len(newMoons)-1]), float64(nextNewMoon), 9)

	// 是否有中气
	// 也需要按照东八区计算
	var hasMiddleChi = func(start, end JulianDayWithLocation) bool {

		for _, jdExtra := range solarTerms {
			cst := jdExtra.JdUT.ToCST()
			if cst >= start && cst < end {
				// 雨水、春分、谷雨、小满、夏至、大暑、处暑、秋分、霜降、小雪、冬至和大寒
				// 也就是可以整除2的
				if jdExtra.Index%2 == 0 {
					return true
				}
			}
		}
		return false
	}

	months := make([]*LunarMonth, IfThenElse(has13, 13, 12).(int))
	var leapMonth = math.MaxInt8
	for i := range months { // 遍历
		var newMoon = newMoons[i]
		var nextNewMoon = newMoons[i+1]

		// 转化成东八区的0点
		var start = newMoon.ToCST().StartOfDay()
		var end = nextNewMoon.ToCST().StartOfDay()

		// 计算是第几月 子月是11月
		var index = 10 + i
		// 还没闰月
		if has13 && leapMonth == math.MaxInt8 && !hasMiddleChi(start, end) {
			leapMonth = index
		}
		// 该月有多少天，按照东八区0点计算
		var days = int(end - start)

		months[i] = &LunarMonth{
			JdUT:  newMoon,
			Index: IfThenElse(index >= leapMonth, index-1, index).(int) % 12,
			Days:  days,
			Leap:  leapMonth == index,
		}
	}

	return months, nil
}

// DogDays 伏天
// 从夏至开始，依照干支纪日的排列，第3个庚日为初伏，第4个庚日为中伏，立秋后第1个庚日为末伏。当夏至与立秋之间出现4个庚日时中伏为10天，出现5个庚日则为20天
func (astro *Astronomy) DogDays(year int) {

}

// Winter9Days 九天
//冬至逢壬日为起点，每“九天”算一“九”，
func (astro *Astronomy) Winter9Days(year int) {

}
