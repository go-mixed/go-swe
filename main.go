package main

import (
	"fmt"
	"go-swe/swe"
)

func main() {
	s := swe.NewSwe()
	d, _ := s.JulDay(2020, 1, 1, 0, swe.Gregorian)
	fmt.Printf("%f", d)



}

