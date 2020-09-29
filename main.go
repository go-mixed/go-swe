package main

import (
	"fmt"
	"go-swe/astro"
	"time"
)

func main() {
	//s := swe.NewSwe()

	fmt.Printf("%v %d \n ", time.Now(), time.Now().Nanosecond())

	t := time.Date(2020, 9, 30, 4, 0, 0, 0, time.UTC)

	jd := astro.TimeToJulianDay(t)
	deltaT := astro.DeltaT(jd)
	et := jd.ToEphemerisTime(deltaT)
	etT := et.ToTime(time.UTC)
	fmt.Printf("JD: %f at %v \n", jd, jd.ToTime(time.UTC))
	fmt.Printf("ET: %f at %v deltaT: %v\n", et, etT, deltaT)

	long, _ := astro.StringToDegree("116°23'")
	lat, _ := astro.StringToDegree("39°54'")
	fmt.Printf("Geo: %f %f\n", long, lat)

	geo := &astro.GeographicCoordinates{
		Longitude: astro.ToRadian(long),
		Latitude:  astro.ToRadian(lat),
	}
	tz, _ := time.LoadLocation("Asia/Shanghai")

	times := astro.SunTwilight(jd, geo, false)

	fmt.Printf("Sun Rise: %v\n", times.Rise.ToTime(nil).In(tz))
	fmt.Printf("Sun Set: %v\n", times.Set.ToTime(nil).In(tz))
	fmt.Printf("Sun Noon: %v\n", times.Noon.ToTime(nil).In(tz))
	fmt.Printf("Sun Midnight: %v\n", times.Midnight.ToTime(nil).In(tz))
	fmt.Printf("Sun Civil : %v %v\n", times.Civil.Dawn.ToTime(nil).In(tz), times.Civil.Dusk.ToTime(nil).In(tz))

}
