package mercator

import "math"

type Coords struct {
	Latitude  float64
	Longitude float64
}

const a = 6378137
const b = 6356752.3142
const f = (a - b) / a

func (c *Coords) Radian() (float64, float64) {
	return toRadian(c.Latitude), toRadian(c.Longitude)
}

func toRadian(num float64) float64 {
	return (num * math.Pi) / 180
}

func e() float64 {
	return math.Sqrt(2*f - math.Pow(f, 2))
}

func (c *Coords) ToMercator() (float64, float64) {
	lat, long := c.Radian()
	mercY := a * long
	mercX := a * math.Log(math.Tan(math.Pi/4+lat/2)*math.Pow((1-e()*math.Sin(lat))/(1+e()*math.Sin(lat)), e()/2))
	return mercX, mercY
}
