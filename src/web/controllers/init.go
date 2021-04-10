package controllers

import (
	"go-swe/src/astro"
	"time"
)

var astronomy *astro.Astronomy
var cacheExpired = time.Hour * 24

func init() {
	astronomy = astro.NewAstronomy()
}
