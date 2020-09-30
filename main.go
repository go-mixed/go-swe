package main

import (
	"fmt"
	"go-swe/astro"
	"time"
)

func main() {

	//for i := 1990; i < 2010; i++ {
	//	t := time.Date(i, 1, 1, 4, 0, 0, 0, time.UTC)
	//	jd := astro.TimeToJulianDay(t)
	//	deltaT1 := astro.DeltaT(jd)
	//	deltaT2 := swe.NewSwe().DeltaT(float64(jd))
	//	fmt.Printf("JD: %v -> %f deltaT %v / %v \n", t, jd, deltaT1, deltaT2)
	//
	//}

	year, month, day := time.Now().Date()
	t := time.Date(year, month, day, 4, 0, 0, 0, time.UTC)

	jd := astro.TimeToJulianDay(t)
	deltaT := astro.DeltaT(jd)
	et := jd.ToEphemerisTime(deltaT)
	etT := et.ToTime(time.UTC)
	fmt.Printf("JD: %f at %v \n", jd, jd.ToTime(time.UTC))
	fmt.Printf("ET: %f at %v deltaT: %v\n", et, etT, deltaT)

	long, _ := astro.StringToDegrees("116°23'")
	lat, _ := astro.StringToDegrees("39°54'")
	fmt.Printf("Geo: %f %f\n", long, lat)

	geo := &astro.GeographicCoordinates{
		Longitude: astro.ToRadians(long),
		Latitude:  astro.ToRadians(lat),
	}
	tz, _ := time.LoadLocation("Asia/Shanghai")

	sunTimes := astro.SunTwilight(jd, geo, false)

	fmt.Printf("Sun Rise: %v\n", sunTimes.Rise.ToTime(nil).In(tz))
	fmt.Printf("Sun Set: %v\n", sunTimes.Set.ToTime(nil).In(tz))
	fmt.Printf("Sun Culmination: %v | %v\n", sunTimes.Culmination.ToTime(nil).In(tz), sunTimes.LowerCulmination.ToTime(nil).In(tz))
	fmt.Printf("Sun Civil : %v | %v\n", sunTimes.Civil.Dawn.ToTime(nil).In(tz), sunTimes.Civil.Dusk.ToTime(nil).In(tz))

	moonTimes := astro.MoonTwilight(jd, geo, false)

	fmt.Printf("Moon Rise: %v\n", moonTimes.Rise.ToTime(nil).In(tz))
	fmt.Printf("Moon Set: %v\n", moonTimes.Set.ToTime(nil).In(tz))
	fmt.Printf("Moon Culmination: %v | %v\n", moonTimes.Culmination.ToTime(nil).In(tz), moonTimes.LowerCulmination.ToTime(nil).In(tz))
}
