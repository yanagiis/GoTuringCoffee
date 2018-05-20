package lib

type Cookbook struct {
	Name        string
	Description string
	Processes   []Process
***REMOVED***

type Process interface{***REMOVED***

type Circle struct {
	Coords      Coordinate `json:"coordinate"`
	Z           Range      `json:"z"`
	Radius      float32    `json:"radius"`
	Cylinder    int        `json:"cylinder"`
	Time        float32    `json:"time"`
	Water       float32    `json:"water"`
	Temperature float32    `json:"temperature"`
***REMOVED***

type Sprial struct {
	Coords      Coordinate `json:"coordinate"`
	Z           Range      `json:"z"`
	Radius      Range      `json:"radius"`
	Cylinder    int        `json:"cylinder"`
	Time        float32    `json:"time"`
	Water       float32    `json:"water"`
	Temperature float32    `json:"temperature"`
***REMOVED***

type Ploygon struct {
	Coords      Coordinate `json:"coordinate"`
	Z           Range      `json:"z"`
	Radius      Range      `json:"radius"`
	Polygon     int        `json:"polygon"`
	Cylinder    int        `json:"cylinder"`
	Time        float32    `json:"time"`
	Water       float32    `json:"water"`
	Temperature float32    `json:"temperature"`
***REMOVED***

type Fixed struct {
	Coords      Coordinate `json:"coordinate"`
	Time        float32    `json:"time"`
	Water       float32    `json:"water"`
	Temperature float32    `json:"temperature"`
***REMOVED***

type Move struct {
	Coords Coordinate `json:"coordinate"`
***REMOVED***

type Wait struct {
	Time float32 `json:"time"`
***REMOVED***

type Mix struct {
	Temperature float32 `json:"temperature"`
***REMOVED***

type Home struct {
***REMOVED***

type Coordinate struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
***REMOVED***

type Range struct {
	From float32 `json:"from"`
	To   float32 `json:"to"`
***REMOVED***
