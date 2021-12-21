package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-common-cache"
	"go-common-web/controllers"
	"go-common/utils/conv"
	"go-swe/src/astro"
	"time"
)

type SolarController struct {
	controllers.Controller
}

func (c *SolarController) TermsByYear() (gin.H, error) {

	year := conv.Atoi(c.Context.Param("year"), 0)
	timezone := c.Context.Query("tz")
	var err error
	var tz *time.Location
	if tz, err = time.LoadLocation(timezone); err != nil {
		tz = time.UTC
	}

	if data, err := cache.Remember(fmt.Sprintf("solar/terms/%d", year), cacheExpired, func() (interface{}, error) {
		return astronomy.SolarTerms(year)
	}); err == nil {
		jds := data.([]*astro.JulianDayExtra)
		var terms map[string]string = map[string]string{}
		for _, jd := range jds {
			terms[astro.SolarTermsString[jd.Index]] = jd.JdUT.ToTime(tz).Format(time.RFC3339)
		}

		return gin.H{
			"year":         year,
			"jds":          data,
			"solar_terms":  terms,
			"term_strings": astro.SolarTermsString,
		}, nil
	} else {
		return nil, controllers.NewResponseException(4011, 400, err.Error())
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
