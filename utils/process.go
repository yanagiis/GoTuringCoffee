package main

import (
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func genCircle() {
	circle := lib.Circle{
		Coords: lib.Coordinate{
			X: 0,
			Y: 0,
			Z: 0,
		},
		ToZ: 0,
		Radius: lib.Range{
			From: 20,
		},
		Cylinder: 8,
		Time:     float64(30),
	}

	points := circle.ToPoints()
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	pts := make(plotter.XYs, len(points))
	for i := range pts {
		if points[i].X != nil && points[i].Y != nil {
			pts[i].X = *points[i].X
			pts[i].Y = *points[i].Y
		}
	}

	p.Title.Text = "Circle"
	if err = plotutil.AddScatters(p, pts); err != nil {
		panic(err)
	}

	if err = p.Save(4*vg.Inch, 4*vg.Inch, "circle.png"); err != nil {
		panic(err)
	}
}

func genSpiral() {
	spiral := lib.Spiral{
		Coords: lib.Coordinate{
			X: 0,
			Y: 0,
			Z: 0,
		},
		ToZ: 0,
		Radius: lib.Range{
			From: 5,
			To:   30,
		},
		Cylinder: 4,
		Time:     float64(30),
	}

	points := spiral.ToPoints()
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	pts := make(plotter.XYs, len(points))
	for i := range pts {
		if points[i].X != nil && points[i].Y != nil {
			pts[i].X = *points[i].X
			pts[i].Y = *points[i].Y
		}
	}

	p.Title.Text = "Spiral"
	if err = plotutil.AddScatters(p, pts); err != nil {
		panic(err)
	}

	if err = p.Save(4*vg.Inch, 4*vg.Inch, "spiral.png"); err != nil {
		panic(err)
	}
}

func genPolygon(n int) {
	polygon := lib.Polygon{
		Coords: lib.Coordinate{
			X: 0,
			Y: 0,
			Z: 0,
		},
		ToZ: 0,
		Radius: lib.Range{
			From: 20,
		},
		Polygon:  int64(n),
		Cylinder: 2,
		Time:     float64(30),
	}

	points := polygon.ToPoints()
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	pts := make(plotter.XYs, len(points))
	for i := range pts {
		if points[i].X != nil && points[i].Y != nil {
			pts[i].X = *points[i].X
			pts[i].Y = *points[i].Y
		}
	}

	p.Title.Text = "Polygon"
	if err = plotutil.AddScatters(p, pts); err != nil {
		panic(err)
	}

	if err = p.Save(4*vg.Inch, 4*vg.Inch, "polygon.png"); err != nil {
		panic(err)
	}
}

func main() {
	genCircle()
	genSpiral()
	genPolygon(5)
}
