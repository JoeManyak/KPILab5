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

func (n *node) parseCSVNoMercator(str string) {
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
	n.coords = mercator.Coords{
		X: Lat,
		Y: Long,
	}
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

func (r rect) fitArea(anotherRect rect) float64 {
	left := math.Max(r.minCoords.X, anotherRect.minCoords.X)
	top := math.Min(r.maxCoords.Y, anotherRect.maxCoords.Y)
	right := math.Min(r.maxCoords.X, r.maxCoords.X)
	bottom := math.Max(r.minCoords.Y, anotherRect.minCoords.Y)

	width := right - left
	height := top - bottom

	if width < 0 || height < 0 {
		return 0
	}

	return width * height
}

func (r rect) addedArea(newNode node) float64 {
	newR := r.rectWithNode(newNode)
	return newR.area() - r.area()
}

func (r rect) area() float64 {
	return math.Abs((r.maxCoords.Y - r.minCoords.Y + 1) * (r.maxCoords.X - r.minCoords.X + 1))
	//return math.Abs((r.maxCoords.Y - r.minCoords.Y) + (r.maxCoords.X - r.minCoords.X))*2
}

func (r rect) rectWithNode(checkNode node) rect {
	r.resize(checkNode)
	return r
}

func (r *rect) resize(newNode node) {
	r.minCoords.GetSmallest(newNode.coords, len(r.nodes) == 0)
	r.maxCoords.GetBiggest(newNode.coords, len(r.nodes) == 0)
}

func (r *rect) insert(newNode node) {
	if len(r.subRects) != 0 {
		/*if r.subRects[0].rectWithNode(newNode).fitArea(*r.subRects[1]) >
			r.subRects[1].rectWithNode(newNode).fitArea(*r.subRects[0]) {
			r.subRects[1].insert(newNode)
		} else if r.subRects[0].rectWithNode(newNode).fitArea(*r.subRects[1]) <
			r.subRects[1].rectWithNode(newNode).fitArea(*r.subRects[0]) {
			r.subRects[0].insert(newNode)
		} else {*/
		if r.subRects[0].addedArea(newNode) < r.subRects[1].addedArea(newNode) {
			r.subRects[0].insert(newNode)
		} else {
			r.subRects[1].insert(newNode)
		}
		//}
		return
	}

	if len(r.nodes) == maxNodes {
		r.split()
		r.insert(newNode)
		return
	}

	r.resize(newNode)

	r.nodes = append(r.nodes, newNode)
}

func (r *rect) takeFurther() (node, node) {
	var maxNode, minNode node
	var maxId, minId int
	minNode = r.nodes[0]
	maxNode = r.nodes[0]
	for i, v := range r.nodes {
		if minNode.coords.X+minNode.coords.Y > v.coords.X+v.coords.Y {
			minId = i
			minNode = v
		}
		if maxNode.coords.X+maxNode.coords.Y < v.coords.X+v.coords.Y {
			maxId = i
			maxNode = v
		}
	}
	r.deleteNode(minId)
	if minId < maxId {
		maxId--
	}
	r.deleteNode(maxId)
	return minNode, maxNode
}

func (r *rect) split() {
	if len(r.nodes) == 0 {
		log.Fatal("Nothing to split")
	}
	left, right := r.takeFurther()
	leftRect := rect{
		minCoords: left.coords,
		maxCoords: left.coords,
		nodes:     []node{left},
	}
	rightRect := rect{
		minCoords: right.coords,
		maxCoords: right.coords,
		nodes:     []node{right},
	}
	*r = subdivide(leftRect, rightRect, *r)
}

func (r *rect) deleteNode(id int) {
	temp := make([]node, len(r.nodes)-1)
	copy(temp[:id], r.nodes[:id])
	copy(temp[id:], r.nodes[id+1:])
	r.nodes = temp
}

func subdivide(leftRect rect, rightRect rect, head rect) rect {
	var selected *rect
	var minId int
	var minArea = -1.0
	for i, v := range head.nodes {
		if minArea == -1.0 || minArea > leftRect.addedArea(v) {
			minArea = leftRect.addedArea(v)
			minId = i
			selected = &leftRect
		}
		if minArea == -1.0 || minArea > rightRect.addedArea(v) {
			minArea = leftRect.addedArea(v)
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

func (r *rect) show() {
	r.showUtil(0)
}

func (r *rect) showUtil(number int) {
	fmt.Printf("%s|%d block [%.2f:%.2f,%.2f:%.2f]: \n", strings.Repeat("->", number*3), number, r.minCoords.X,
		r.minCoords.Y, r.maxCoords.X, r.maxCoords.Y)
	if len(r.subRects) != 0 {
		r.subRects[0].showUtil(number + 1)
		r.subRects[1].showUtil(number + 1)
	} else {
		for _, v := range r.nodes {
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
	/*	var i = 0
		for reader.Scan() {
			i++*/
	for i := 0; i < 30000; i++ {
		var newNode node
		reader.Scan()
		//newNode.parseCSVNoMercator(reader.Text())
		newNode.parseCSV(reader.Text())
		head.insert(newNode)
	}
	fmt.Println(head)
}
