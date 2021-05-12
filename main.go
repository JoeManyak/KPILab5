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
		if head.subRects[0].area() < head.subRects[1].area() {
			head.subRects[0].insert(newNode)
		} else {
			head.subRects[1].insert(newNode)
		}
		return
	}
	if len(head.nodes) == maxNodes {
		head.split()
		head.insert(newNode)
		return
	}

	head.minCoords.GetSmallest(newNode.coords, len(head.nodes) == 0)
	head.maxCoords.GetBiggest(newNode.coords, len(head.nodes) == 0)

	head.nodes = append(head.nodes, newNode)
}

func (head *rect) split() {
	if len(head.nodes) == 0 {
		log.Fatal("Nothing to split")
	}
	*head = subdivide(rect{}, rect{}, *head)
}

func (head *rect) deleteNode(id int) {
	temp := make([]node, len(head.nodes)-1)
	copy(temp[:id], head.nodes[:id])
	copy(temp[id:], head.nodes[id+1:])
	head.nodes = temp
}

func subdivide(leftRect rect, rightRect rect, head rect) rect {
	var selected *rect
	var minId int
	var minArea = -1.0
	for i, v := range head.nodes {
		tempRect := leftRect
		tempRect.nodes = append(tempRect.nodes, v)
		if minArea == -1.0 || minArea > tempRect.area() {
			minArea = tempRect.area()
			minId = i
			selected = &leftRect
		}
		tempRect = rightRect
		tempRect.nodes = append(tempRect.nodes, v)
		if minArea == -1.0 || minArea > tempRect.area() {
			minArea = tempRect.area()
			minId = i
			selected = &rightRect
		}
	}
	if selected != nil {
		(*selected).resize(head.nodes[minId])
		(*selected).nodes = append((*selected).nodes, head.nodes[minId])
		head.deleteNode(minId)
	}
	if len(head.nodes) == 0 {
		head.subRects = []*rect{&leftRect, &rightRect}
		return head
	}
	return subdivide(leftRect, rightRect, head)
}

func (head *rect) show() {
	head.showUtil(0)
}

func (head *rect) showUtil(number int) {
	fmt.Printf("%s|%d block: \n", strings.Repeat("->", number*3), number)
	if len(head.subRects) != 0 {
		head.subRects[0].showUtil(number + 1)
		head.subRects[1].showUtil(number + 1)
	} else {
		for _, v := range head.nodes {
			fmt.Printf("%s|%v \n", strings.Repeat("->", number*3), v)
		}
	}
}

func main() {
	file, err := os.Open("./src/ukraine_poi.csv")
	if err != nil {
		log.Fatal("No such file")
	}
	var head rect
	reader := bufio.NewScanner(file)
	for i := 0; i < 500; i++ {
		var newNode node
		reader.Scan()
		newNode.parseCSV(reader.Text())
		head.insert(newNode)
	}
	fmt.Println(head)
	head.show()
}
