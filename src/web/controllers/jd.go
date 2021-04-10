package controllers

import (
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"go-common-web/controllers"
	"go-swe/src/astro"
	"time"
)

type JDController struct {
	controllers.Controller
}

type jdConvert struct {
	Date string `uri:"date" binding:"required"`
}

func (c *JDController) JdConvert(ctx *gin.Context) (gin.H, *controllers.ResponseException) {

	date := ctx.DefaultQuery("date", time.Now().Format(time.RFC3339))

	if t, err := dateparse.ParseAny(date); err == nil {
		jd := astro.TimeToJulianDay(t)
		//deltaT := astronomy.DeltaT(jd)
		jd_et := astro.NewEphemerisTime(jd)
		return gin.H{
			"jd":      jd,
			"jd2000":  jd.AsJD2000(),
			"delta_t": jd_et.DeltaT,
			"jd_et":   jd_et.Value(),
			"date":    date,
		}, nil
	} else {
		return nil, controllers.NewResponseException(4001, err.Error())
	}
}
