package main

import (
	"./mercator"
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const maxNodes = 10

type node struct {
	coords        mercator.Coords
	nodeType      string //Тип
	nodeSubType   string //Підтип
	name          string //Назва
	address       string //Адреса
	addressNumber string //Адреса
}

func (n *node) parseCSV(str string) {
	sl := strings.Split(str, ";")
	//Get coords
	Lat, err := strconv.ParseFloat(sl[0], 64)
	if err != nil {
		log.Fatal("Invalid data:", sl[0], "is not a float value")
	}
	Long, err := strconv.ParseFloat(sl[1], 64)
	if err != nil {
		log.Fatal("Invalid data:", sl[1], "is not a float value")
	}
	geo := mercator.GeoCoords{
		Latitude:  Lat,
		Longitude: Long,
	}
	n.coords = geo.ToMercator()
	n.nodeType = sl[2]
	n.nodeSubType = sl[3]
	n.name = sl[4]
	n.address = sl[5]
	n.addressNumber = sl[6]
}

type rect struct {
	minCoords mercator.Coords // Нижня ліва частина прямокутника
	maxCoords mercator.Coords // Верхня права частина прямокутника
	subRects  []*rect
	nodes     []node
}

func (head *rect) area() float64 {
	return math.Abs((head.maxCoords.Y - head.minCoords.Y) * (head.maxCoords.X - head.minCoords.X))
}

func (head *rect) resize(newNode node) {
	head.minCoords.GetSmallest(newNode.coords, len(head.nodes) == 0)
	head.maxCoords.GetBiggest(newNode.coords, len(head.nodes) == 0)
}

func (head *rect) insert(newNode node) {
	if len(head.subRects) != 0 {
		var minId = 0
		var minArea = head.subRects[0].area()
		var curArea = head.area()
		for i, v := range head.subRects {
			if curArea < minArea {
				minId = i
				minArea = v.area()
			}
		}
		head.subRects[minId].insert(newNode)
		return
	}
	if len(head.nodes) == maxNodes {
		//Split
	}

	head.minCoords.GetSmallest(newNode.coords, len(head.nodes) == 0)
	head.maxCoords.GetBiggest(newNode.coords, len(head.nodes) == 0)

	head.nodes = append(head.nodes, newNode)
}

func () {

}

func main() {
	file, err := os.Open("./src/ukraine_poi.csv")
	if err != nil {
		log.Fatal("No such file")
	}
	var head rect
	reader := bufio.NewScanner(file)
	for i := 0; i < 5; i++ {
		var newNode node
		reader.Scan()
		newNode.parseCSV(reader.Text())
		head.insert(newNode)
		fmt.Println(head.minCoords, head.maxCoords)
	}
}
