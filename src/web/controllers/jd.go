package controllers

import (
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"go-swe/src/astro"
	"gopkg.in/go-mixed/go-common.v1/web.v1/controllers"
	"time"
)

type JDController struct {
	controllers.Controller
}

type jdConvert struct {
	Date string `uri:"date" binding:"required"`
}

func (c *JDController) JdConvert(ctx *gin.Context) (gin.H, error) {

	date := ctx.DefaultQuery("date", time.Now().Format(time.RFC3339))

	if t, err := dateparse.ParseAny(date); err == nil {
		jd := astro.TimeToJulianDay(t)
		//deltaT := astronomy.DeltaT(jd)
		jdEt := astro.NewEphemerisTime(jd)
		return gin.H{
			"jd":      jd,
			"jd2000":  jd.ToJD2000(),
			"delta_t": jdEt.DeltaT,
			"jd_et":   jdEt.Value(),
			"date":    date,
		}, nil
	} else {
		return nil, controllers.NewResponseException(4001, 400, err.Error())
	}
}
