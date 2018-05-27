package lib

type Cookbook struct {
	Name        string
	Description string
	Processes   []Process
}

type Process interface{}

type Circle struct {
	Coords      Coordinate `json:"coordinate"`
	ToZ         float32    `json:"toz"`
	Radius      float32    `json:"radius"`
	Cylinder    int        `json:"cylinder"`
	Time        float32    `json:"time"`
	Water       float32    `json:"water"`
	Temperature float32    `json:"temperature"`
}

type Sprial struct {
	Coords      Coordinate `json:"coordinate"`
	ToZ         float32    `json:"toz"`
	Radius      Range      `json:"radius"`
	Cylinder    int        `json:"cylinder"`
	Time        float32    `json:"time"`
	Water       float32    `json:"water"`
	Temperature float32    `json:"temperature"`
}

type Ploygon struct {
	Coords      Coordinate `json:"coordinate"`
	ToZ         float32    `json:"toz"`
	Radius      Range      `json:"radius"`
	Polygon     int        `json:"polygon"`
	Cylinder    int        `json:"cylinder"`
	Time        float32    `json:"time"`
	Water       float32    `json:"water"`
	Temperature float32    `json:"temperature"`
}

type Fixed struct {
	Coords      Coordinate `json:"coordinate"`
	Time        float32    `json:"time"`
	Water       float32    `json:"water"`
	Temperature float32    `json:"temperature"`
}

type Move struct {
	Coords Coordinate `json:"coordinate"`
}

type Wait struct {
	Time float32 `json:"time"`
}

type Mix struct {
	Temperature float32 `json:"temperature"`
}

type Home struct {
}

type Coordinate struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Range struct {
	From float32 `json:"from"`
	To   float32 `json:"to"`
}
