package lib

import (
	"fmt"
	"math"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	PointInterval = float64(2)
	DefaultF      = float64(5000)
)

// Cookbook Coffee cookbook
type Cookbook struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Notes       []string  `json:"notes"`
	Setup       []Process `json:"setup"`
	Processes   []Process `json:"processes"`
	Teardown    []Process `json:"teardown"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToPoints Generate points for all processes
func (c *Cookbook) ToPoints() []Point {
	var points []Point
	for index, p := range c.Processes {
		log.Info().Msgf("Process %d ToPoints", index)
		points = append(points, p.ToPoints()...)
	}
	return points
}

// GetTotalWater Total water
func (c *Cookbook) GetTotalWater() float64 {
	var totalWater float64

	for _, p := range c.Processes {
		totalWater += p.GetWater()
	}
	return totalWater
}

// GetTotalTime Total time
func (c *Cookbook) GetTotalTime() float64 {
	var totalTime float64

	for _, p := range c.Processes {
		totalTime += p.GetTime()
	}
	return totalTime
}

// ProcessImpl Process implementation
type ProcessImpl interface {
	ToPoints() []Point
	GetWater() float64
	GetTime() float64
	GetTemperature() float64
}

// Process Cookbook process
type Process struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Impl      ProcessImpl `json:"impl"`
}

// GetImpl Get the process implementation class
func (p *Process) GetImpl() ProcessImpl {
	return p.Impl
}

// ToPoints Generate points from process implementation
func (p *Process) ToPoints() []Point {
	return p.GetImpl().ToPoints()
}

// GetWater Calcuate water usage
func (p *Process) GetWater() float64 {
	return p.GetImpl().GetWater()
}

// GetTime Calcuate required time of process
func (p *Process) GetTime() float64 {

	return p.GetImpl().GetTime()
}

// GetTemperature Get the water temperature in process
func (p *Process) GetTemperature() float64 {

	return p.GetImpl().GetTemperature()
}

func NewProcessImpl(processName string) (ProcessImpl, error) {
	var processImpl ProcessImpl

	switch processName {
	case "Circle":
		processImpl = new(Circle)
	case "Spiral":
		processImpl = new(Spiral)
	case "Polygon":
		processImpl = new(Polygon)
	case "Fixed":
		processImpl = new(Fixed)
	case "Move":
		processImpl = new(Move)
	case "Wait":
		processImpl = new(Wait)
	case "Mix":
		processImpl = new(Mix)
	case "Home":
		processImpl = new(Home)
	default:
		return processImpl, fmt.Errorf("Not support process '%s'", processName)
	}

	return processImpl, nil
}

type Circle struct {
	Coords      Coordinate `json:"coordinate"`
	ToZ         float64    `json:"toz" unit:"mm"`
	Radius      Range      `json:"radius" unit:"mm"`
	Cylinder    int64      `json:"cylinder"`
	Time        float64    `json:"time" unit:"s"`
	Water       float64    `json:"water" unit:"ml"`
	Temperature float64    `json:"temperature" unit:"celsius"`
}

// ToPoints Use process params to generate points
func (c *Circle) ToPoints() []Point {
	src := Coordinate{
		X: c.Coords.X + c.Radius.From,
		Y: c.Coords.Y,
		Z: c.Coords.Z,
	}
	dst := Coordinate{
		X: c.Coords.X + c.Radius.From,
		Y: c.Coords.Y,
		Z: c.ToZ,
	}

	points := makeSpiral(&src, &dst, &c.Coords, c.Cylinder)
	pathlen := float64(len(points)-1) * PointInterval
	feedrate := pathlen / (c.Time / 60)
	pointwater := c.Water / float64(len(points))
	for i := range points {
		points[i].T = &c.Temperature
		points[i].F = &feedrate
		points[i].E = &pointwater
	}

	points = append(makeMove(&src), points...)
	return points
}

func (p *Circle) GetTime() float64 {
	return p.Time
}

func (p *Circle) GetWater() float64 {
	return p.Water
}

func (p *Circle) GetTemperature() float64 {
	return p.Temperature
}

type Spiral struct {
	Coords      Coordinate `json:"coordinate"`
	ToZ         float64    `json:"toz" unit:"mm"`
	Radius      Range      `json:"radius" unit:"mm"`
	Cylinder    int64      `json:"cylinder"`
	Time        float64    `json:"time" unit:"s"`
	Water       float64    `json:"water" unit:"ml"`
	Temperature float64    `json:"temperature" unit:"celsius"`
}

func (s *Spiral) ToPoints() []Point {
	src := Coordinate{
		X: s.Coords.X + s.Radius.From,
		Y: s.Coords.Y,
		Z: s.Coords.Z,
	}
	dst := Coordinate{
		X: s.Coords.X + s.Radius.To,
		Y: s.Coords.Y,
		Z: s.ToZ,
	}

	log.Info().Msgf("Spiral: Src %f -> Dst %f", s.Radius.From, s.Radius.To)
	points := makeSpiral(&src, &dst, &s.Coords, s.Cylinder)
	pathlen := float64(len(points)-1) * PointInterval
	feedrate := pathlen / (s.Time / 60)
	pointwater := s.Water / float64(len(points))
	for i := range points {
		points[i].T = &s.Temperature
		points[i].F = &feedrate
		points[i].E = &pointwater
	}

	points = append(makeMove(&src), points...)
	return points
}

func (p *Spiral) GetTime() float64 {
	return p.Time
}

func (p *Spiral) GetWater() float64 {
	return p.Water
}

func (p *Spiral) GetTemperature() float64 {
	return p.Temperature
}

type Polygon struct {
	Coords      Coordinate `json:"coordinate"`
	ToZ         float64    `json:"toz" unit:"mm"`
	Radius      Range      `json:"radius" unit:"mm"`
	Polygon     int64      `json:"polygon" unit:"mm"`
	Cylinder    int64      `json:"cylinder"`
	Time        float64    `json:"time" unit:"s"`
	Water       float64    `json:"water" unit:"ml"`
	Temperature float64    `json:"temperature" unit:"celsius"`
}

func (p *Polygon) ToPoints() []Point {
	var points []Point

	angel := float64(360) / float64(p.Polygon)
	theta := float64(360) / float64(p.Cylinder)

	base := Coordinate{
		X: p.Coords.X + p.Radius.From,
		Y: p.Coords.Y,
		Z: p.Coords.Z,
	}

	for i := int64(0); i < p.Cylinder; i++ {
		src := p.Coords.rotate(theta*float64(i), &base)
		points = append(points, makeMove(&src)...)
		for j := int64(0); j < p.Polygon; j++ {
			dst := p.Coords.rotate(angel, &src)
			points = append(points, makeLine(&src, &dst)...)
			src = dst
		}
	}

	pathlen := float64(int64(len(points))-p.Cylinder) * PointInterval
	feedrate := pathlen / (p.Time / 60)
	pointwater := p.Water / float64(int64(len(points))-p.Cylinder)
	for i := range points {
		if points[i].F != nil {
			points[i].T = &p.Temperature
			points[i].F = &feedrate
			points[i].E = &pointwater
		}
	}

	return points
}

func (p *Polygon) GetTime() float64 {
	return p.Time
}

func (p *Polygon) GetWater() float64 {
	return p.Water
}

func (p *Polygon) GetTemperature() float64 {
	return p.Temperature
}

type Fixed struct {
	Coords      Coordinate `json:"coordinate"`
	Time        float64    `json:"time" unit:"mm"`
	Water       float64    `json:"water" unit:"ml"`
	Temperature float64    `json:"temperature" unit:"celsius"`
}

func (f *Fixed) ToPoints() []Point {
	points := makeMove(&f.Coords)
	waterPerPoint := f.Water / (f.Time * float64(10))
	feedrate := float64(0.1)
	numOfPoint := int(f.Time) * 10
	for i := 0; i < numOfPoint; i++ {
		points = append(points, Point{
			Type: PointT,
			E:    &waterPerPoint,
			F:    &feedrate,
			T:    &f.Temperature,
		})
	}
	return points
}

func (p *Fixed) GetTime() float64 {
	return p.Time
}

func (p *Fixed) GetWater() float64 {
	return p.Water
}

func (p *Fixed) GetTemperature() float64 {
	return p.Temperature
}

type Move struct {
	Coords Coordinate `json:"coordinate"`
}

func (m *Move) ToPoints() []Point {
	return makeMove(&m.Coords)
}

func (p *Move) GetTime() float64 {
	return 0
}

func (p *Move) GetWater() float64 {
	return 0
}

func (p *Move) GetTemperature() float64 {
	return 0
}

type Wait struct {
	Time float64 `json:"time" unit:"s"`
}

func (w *Wait) ToPoints() []Point {
	return []Point{
		Point{
			Type: WaitT,
			Time: &w.Time,
		},
	}
}

func (p *Wait) GetTime() float64 {
	return 0
}

func (p *Wait) GetWater() float64 {
	return 0
}

func (p *Wait) GetTemperature() float64 {
	return 0
}

type Mix struct {
	Temperature float64 `json:"temperature" unit:"celsius"`
}

func (m *Mix) ToPoints() []Point {
	return []Point{
		Point{
			Type: MixT,
			T:    &m.Temperature,
		},
	}
}

func (p *Mix) GetTime() float64 {
	return 0
}

func (p *Mix) GetWater() float64 {
	return 0
}

func (p *Mix) GetTemperature() float64 {
	return p.Temperature
}

type Home struct {
}

func (h *Home) ToPoints() []Point {
	return []Point{
		Point{
			Type: HomeT,
		},
	}
}

func (p *Home) GetTime() float64 {
	return 0
}

func (p *Home) GetWater() float64 {
	return 0
}

func (p *Home) GetTemperature() float64 {
	return 0
}

type Coordinate struct {
	X float64 `json:"x" unit:"mm"`
	Y float64 `json:"y" unit:"mm"`
	Z float64 `json:"z" unit:"mm"`
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

func (c *Coordinate) rotate(theta float64, p *Coordinate) Coordinate {
	radian := math.Pi * (theta / 180)
	cos := math.Cos(radian)
	sin := math.Sin(radian)
	return Coordinate{
		X: p.X*cos - p.Y*sin + c.X,
		Y: p.X*sin + p.Y*cos + c.Y,
		Z: p.Z,
	}
}

type Range struct {
	From float64 `json:"from" unit:"mm"`
	To   float64 `json:"to" unit:"mm"`
}

func makeSpiral(src, dst, center *Coordinate, cylinder int64) []Point {
	rotateTheta := float64(cylinder * 360)
	srcXY := &Coordinate{
		X: src.X,
		Y: src.Y,
	}
	dstXY := &Coordinate{
		X: dst.X,
		Y: dst.Y,
	}
	centerXY := &Coordinate{
		X: center.X,
		Y: center.Y,
	}
	radiusSrc := centerXY.distance(srcXY)
	radiusDst := centerXY.distance(dstXY)
	radiusPerDegree := (radiusDst - radiusSrc) / rotateTheta
	log.Info().Msgf("Center XY: %f, %f", center.X, center.Y)
	log.Info().Msgf("Cylinder %d, radiusSrc: %f radiusDst: %f radiusPerDegree: %f", cylinder, radiusSrc, radiusDst, radiusPerDegree)
	zPerDegree := (dst.Z - src.Z) / rotateTheta

	var points []Point
	radius := radiusSrc
	theta := float64(0)
	for {
		if radius == 0 {
			radius = 0.1
		}

		/*
		   |INFO| Process 14 ToPoints
		    2020-01-05T20:53:16+08:00 |INFO| Spiral: Src 0.000000 -> Dst 15.000000
		    2020-01-05T20:53:16+08:00 |INFO| Center XY: 0.000000, 0.000000
		    2020-01-05T20:53:16+08:00 |INFO| Cylinder 3, radiusSrc: 0.000000 radiusDst: 15.000000 radiusPerDegree: 0.013889
		    2020-01-05T20:53:16+08:00 |INFO| 360*PointInterval=(720.000000) / 2 * math.Pi * radius=(0.628319) => 1145.915590
		    2020-01-05T20:53:16+08:00 |INFO| stepTheta:1145.915590, PointInterval:2.000000, radius:16.015494
		    2020-01-05T20:53:16+08:00 |WARN| 0.000000 15.000000 - points: 0 => theta 1145.915590 > rotateTheta 1080.000000
		    2020-01-05T20:53:16+08:00 |WARN| Points size is zero

		*/
		stepTheta := float64(360*PointInterval) / float64(2*math.Pi*radius)
		log.Info().Msgf("360*PointInterval=(%f) / 2 * math.Pi * radius=(%f) => %f", 360*PointInterval, 2*math.Pi*radius, (360*PointInterval)/(2*math.Pi*radius))

		radius += stepTheta * radiusPerDegree
		theta += stepTheta

		if theta > rotateTheta {
			log.Info().Msgf("stepTheta:%f, PointInterval:%f, radius:%f", stepTheta, PointInterval, radius)
			log.Warn().Msgf("%f %f - points: %d => theta %f > rotateTheta %f", src.X, dst.X, len(points), theta, rotateTheta)
			break
		}

		base := Coordinate{X: radius, Y: 0, Z: src.Z}
		coord := center.rotate(theta, &base)
		coord.Z += zPerDegree * theta
		points = append(points, Point{
			Type: PointT,
			X:    &coord.X,
			Y:    &coord.Y,
			Z:    &coord.Z,
		})
	}

	if len(points) > 0 {
		lastPoint := points[len(points)-1]
		if *lastPoint.X != dst.X || *lastPoint.Y != dst.Y || *lastPoint.Z != dst.Z {
			points = append(points, Point{
				Type: PointT,
				X:    &dst.X,
				Y:    &dst.Y,
				Z:    &dst.Z,
			})
		}
	} else {
		log.Warn().Msg("Points size is zero")
	}

	return points
}

func makeLine(src, dst *Coordinate) []Point {
	deltaX, deltaY, deltaZ := src.delta(dst)
	distance := src.distance(dst)
	numPoints := distance / PointInterval
	stepX := deltaX / numPoints
	stepY := deltaY / numPoints
	stepZ := deltaZ / numPoints

	var points []Point
	for i := 0; i <= int(numPoints); i += 1 {
		fi := float64(i)
		x := src.X + fi*stepX
		y := src.Y + fi*stepY
		z := src.Z + fi*stepZ
		points = append(points, Point{
			Type: PointT,
			X:    &x,
			Y:    &y,
			Z:    &z,
		})
	}
	return points
}

func makeMove(dst *Coordinate) []Point {
	f := DefaultF
	return []Point{
		Point{
			Type: PointT,
			Z:    &dst.Z,
			F:    &f,
		},
		Point{
			Type: PointT,
			X:    &dst.X,
			Y:    &dst.Y,
			F:    &f,
		},
	}
}
