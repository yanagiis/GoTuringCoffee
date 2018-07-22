package lib

import (
	"math"

	"github.com/globalsign/mgo/bson"
)

const PointInterval = 2

type Cookbook struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Processes   []Process     `json:"processes"`
}

type Process interface {
	ToPoints() []Point
}

type Circle struct {
	Coords      Coordinate `json:"coordinate"`
	ToZ         float64    `json:"toz"`
	Radius      float64    `json:"radius"`
	Cylinder    int64      `json:"cylinder"`
	Time        float64    `json:"time"`
	Water       float64    `json:"water"`
	Temperature float64    `json:"temperature"`
}

func (c *Circle) ToPoints() {
	src := Coordinate{
		X: c.Coords.X + c.Radius,
		Y: c.Coords.Y,
		Z: c.Coords.Z,
	}
	dst := Coordinate{
		X: c.Coords.X + c.Radius,
		Y: c.Coords.Y,
		Z: c.Coords.Z,
	}
}

type Sprial struct {
	Coords      Coordinate `json:"coordinate"`
	ToZ         float64    `json:"toz"`
	Radius      Range      `json:"radius"`
	Cylinder    int        `json:"cylinder"`
	Time        float64    `json:"time"`
	Water       float64    `json:"water"`
	Temperature float64    `json:"temperature"`
}

type Ploygon struct {
	Coords      Coordinate `json:"coordinate"`
	ToZ         float64    `json:"toz"`
	Radius      Range      `json:"radius"`
	Polygon     int        `json:"polygon"`
	Cylinder    int        `json:"cylinder"`
	Time        float64    `json:"time"`
	Water       float64    `json:"water"`
	Temperature float64    `json:"temperature"`
}

type Fixed struct {
	Coords      Coordinate `json:"coordinate"`
	Time        float64    `json:"time"`
	Water       float64    `json:"water"`
	Temperature float64    `json:"temperature"`
}

type Move struct {
	Coords Coordinate `json:"coordinate"`
}

type Wait struct {
	Time float64 `json:"time"`
}

type Mix struct {
	Temperature float64 `json:"temperature"`
}

type Home struct {
}

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (c *Coordinate) delta(other *Coordinate) (float64, float64, float64) {
	return other.X - c.X, other.Y - c.Y, other.Z - c.Z
}

func (c *Coordinate) distance(other *Coordinate) float64 {
	deltaX, deltaY, deltaZ := c.delta(other)
	squareX := math.Pow(deltaX, 2)
	squareY := math.Pow(deltaY, 2)
	squareZ := math.Pow(deltaZ, 2)
	return float64(math.Sqrt(squareX + squareY + squareZ))
}

func (c *Coordinate) rotate(theta float64, p *Coordinate) {
	radian := math.Pi * theta
	cos := math.Cos(radian)
	sin := math.Sin(radian)
	p.X = p.X*cos - p.Y*sin + c.X
	p.Y = p.X*sin + p.Y*cos + c.Y
	p.Z = p.Z
}

type Range struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
}

func makeSpiral(src, dst, center *Coordinate, cylinder int64) *[]Point {
	rotateTheta := float64(cylinder * 360)
	radiusSrc := center.distance(src)
	radiusDst := center.distance(dst)
	radiusPerDegree := (radiusDst - radiusSrc) / rotateTheta
	zPerDegree := (dst.Z - src.Z) / rotateTheta

	var points []Point
	radius := radiusSrc
	for theta := float64(0); theta < rotateTheta; {
		coord := Coordinate{X: radius, Y: 0, Z: src.Z}
		center.rotate(theta, &coord)
		coord.Z += zPerDegree * theta
		points = append(points, Point{
			X: &coord.X,
			Y: &coord.Y,
			Z: &coord.Z,
		})
		stepTheta := (360 * PointInterval) / (2 * math.Pi * radius)
		radius = stepTheta * radiusPerDegree
		theta += stepTheta
	}
	return &points
}

func makeLine(src, dst *Coordinate) *[]Point {
	deltaX, deltaY, deltaZ := src.delta(dst)
	distance := src.distance(dst)
	numPoints := distance / PointInterval
	stepX := deltaX / numPoints
	stepY := deltaY / numPoints
	stepZ := deltaZ / numPoints

	var points []Point
	for i := 0; i < int(numPoints); i += 1 {
		fi := float64(i)
		x := src.X + fi*stepX
		y := src.Y + fi*stepY
		z := src.Z + fi*stepZ
		points = append(points, Point{
			X: &x,
			Y: &y,
			Z: &z,
		})
	}
	return &points
}
