package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-swe/src/astro"
	"gopkg.in/go-mixed/go-common.v1/cache.v1"
	"gopkg.in/go-mixed/go-common.v1/utils/conv"
	"gopkg.in/go-mixed/go-common.v1/web.v1/controllers"
	"time"
)

type LunarController struct {
	controllers.Controller
}

func (c *LunarController) PhasesByYear() (gin.H, error) {
	year := conv.Atoi(c.Context.Param("year"), 0)

	if data, err := cache.Remember(fmt.Sprintf("lunar/phases/%d", year), cacheExpired, func() (interface{}, error) {
		return astronomy.LunarPhases(year)
	}); err == nil {
		return gin.H{
			"year":          year,
			"result":        data,
			"phase_strings": astro.LunarPhaseStrings,
		}, nil
	} else {
		return nil, controllers.NewResponseException(4021, 400, err.Error())
	}
}

func (c *LunarController) MonthsByYear() (gin.H, error) {
	year := conv.Atoi(c.Context.Param("year"), 0)

	if data, err := cache.Remember(fmt.Sprintf("lunar/months/%d", year), cacheExpired, func() (interface{}, error) {
		return astronomy.LunarMonths(year)
	}); err == nil {
		lunarMonths := data.([]*astro.LunarMonth)
		tz, _ := time.LoadLocation("Asia/Shanghai")

		type lunarMonth struct {
			JdUT astro.JulianDay `json:"jd_ut"`
			At   string          `json:"at"`
			Days int             `json:"days"`
			Leap bool            `json:"leap"`
		}

		var _lunarMonths map[string]lunarMonth = map[string]lunarMonth{}
		for _, month := range lunarMonths {
			_lunarMonths[astro.LunarMonthStrings[month.Index]] = lunarMonth{
				JdUT: month.JdUT,
				At:   month.JdUT.ToTime(tz).Format(time.RFC3339),
				Leap: month.Leap,
				Days: month.Days,
			}
		}

		return gin.H{
			"year":         year,
			"lunar_months": _lunarMonths,
		}, nil
	} else {
		return nil, controllers.NewResponseException(4022, 400, err.Error())
	}
}
