package main

import (
	"./mercator"
	"fmt"
)

type rect struct {
	minLatitude  float64
	minLongitude float64
	maxLatitude  float64
	maxLongitude float64
	subRects     []*rect
	nodes        []*node
}

type node struct {
	latitude    float64 //Географічна широта
	longitude   float64 //Географічна довгота
	nodeType    string  //Тип
	nodeSubType string  //Підтип
	address     string  //Адреса
}

func main() {
	a := mercator.Coords{
		Latitude:  55.751667,
		Longitude: 37.617778,
	}
	fmt.Println(a.ToMercator())
}
