package lib

import (
	"time"
)

const (
	CodeGet uint8 = iota
	CodePut
	CodeSuccess
	CodeFailure
)

type Request struct {
	Code uint8
}

type Response struct {
	Code uint8
	Msg  string
}

type BaristaRequest struct {
	Request
	Points []Point
}

type TempResponse struct {
	Response
	Payload TempRecord
}

type DistanceResponse struct {
	Response
	Payload DistanceRecord
}

type ReplenisherRequest struct {
	Request
	Stop bool
}

type ReplenisherResponse struct {
	Response
	Payload ReplenisherRecord
}

type FullResponse struct {
	Response
	Payload FullRecord
}

type HeaterRequest struct {
	Request
	Temp float64
}

type HeaterResponse struct {
	Response
	Payload HeaterRecord
}

type TempRecord struct {
	Temp float64
	Time time.Time
}

type ReplenisherRecord struct {
	IsReplenishing bool
	Time           time.Time
}

type FullRecord struct {
	IsFull bool
	Time   time.Time
}

type HeaterRecord struct {
	Duty   float64
	Target float64
	Period time.Duration
	Time   time.Time
}

type DistanceRecord struct {
	Distance int
	Time     time.Time
}

func (r Request) IsGet() bool {
	return r.Code == CodeGet
}

func (r Request) IsPut() bool {
	return r.Code == CodePut
}

func (r Response) IsFailure() bool {
	return r.Code != CodeSuccess
}
