package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-common-cache"
	"go-common-web/controllers"
	"go-common/utils/conv"
	"go-swe/src/astro"
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
		return gin.H{
			"year":          year,
			"result":        data,
			"month_strings": astro.LunarMonthStrings[0:12],
			"day_strings":   astro.LunarDayStrings,
			"leap_string":   astro.LunarMonthStrings[12],
		}, nil
	} else {
		return nil, controllers.NewResponseException(4022, 400, err.Error())
	}
}
