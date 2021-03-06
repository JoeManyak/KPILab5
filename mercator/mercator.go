package mercator

import "math"

type GeoCoords struct {
	Latitude  float64
	Longitude float64
}

type Coords struct {
	X float64
	Y float64
}

func (c *Coords) GetSmallest(newCoords Coords) {
	if c.X > newCoords.X {
		c.X = newCoords.X
	}
	if c.Y > newCoords.Y {
		c.Y = newCoords.Y
	}
}

func (c *Coords) GetBiggest(newCoords Coords) {
	if c.X < newCoords.X {
		c.X = newCoords.X
	}
	if c.Y < newCoords.Y {
		c.Y = newCoords.Y
	}
}

func (c *Coords) InRadius(newCoords Coords, radius float64) bool {
	return math.Sqrt(math.Pow(c.Y-newCoords.Y, 2)+math.Pow(c.X-newCoords.X, 2)) <= radius
}

const a = 6378137
const b = 6356752.3142
const f = (a - b) / a

func (c *GeoCoords) Radian() (float64, float64) {
	return toRadian(c.Latitude), toRadian(c.Longitude)
}

func toRadian(num float64) float64 {
	return (num * math.Pi) / 180
}

func e() float64 {
	return math.Sqrt(2*f - math.Pow(f, 2))
}

func (c *GeoCoords) ToMercator() Coords {
	var newCoords Coords
	lat, long := c.Radian()
	newCoords.Y = a * long
	newCoords.X = a * math.Log(math.Tan(math.Pi/4+lat/2)*math.Pow((1-e()*math.Sin(lat))/(1+e()*math.Sin(lat)), e()/2))
	return newCoords
}
