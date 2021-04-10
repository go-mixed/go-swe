package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-common-web/controllers"
	"go-common/cache"
	"go-common/utils"
	"go-swe/src/astro"
)

type SolarController struct {
	controllers.Controller
}

func (c *SolarController) TermsByYear() (gin.H, *controllers.ResponseException) {

	year := utils.Atoi(c.Context.Param("year"), 0)

	if data, err := cache.Remember(fmt.Sprintf("solar/terms/%d", year), cacheExpired, func() (interface{}, error) {
		return astronomy.SolarTerms(year)
	}); err == nil {
		return gin.H{
			"year":         year,
			"result":       data,
			"term_strings": astro.SolarTermsString,
		}, nil
	} else {
		return nil, controllers.NewResponseException(4011, err.Error())
	}
}

//func (c *SolarController) TermsByRange() (gin.H, *ResponseException)  {
//	start := c.ctx.DefaultQuery("start", "")
//	end := c.ctx.DefaultQuery("end", "")
//
//	if res, err := astronomy.SolarTermsRange(astro.JulianDay(start), astro.JulianDay(end)); err == nil {
//		return gin.H{
//			"start": start,
//			"end": end,
//			"start_jd": start_jd,
//			"end_jd": end_jd,
//			"result": res,
//		}, nil
//	} else {
//		return nil, NewResponseException(4012, err.Error())
//	}
//}
