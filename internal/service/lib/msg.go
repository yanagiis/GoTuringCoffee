package lib

import "time"

const (
	CodeGet uint8 = iota
	CodePut
	CodeSuccess
	CodeFailure
***REMOVED***

type Request struct {
	Code uint8
***REMOVED***

type Response struct {
	Code uint8
	Msg  string
***REMOVED***

type BaristaResponse struct {
***REMOVED***

type TempResponse struct {
	Response
	Payload TempRecord
***REMOVED***

type ReplenisherRequest struct {
	Request
	Stop bool
***REMOVED***

type ReplenisherResponse struct {
	Response
	Payload ReplenisherRecord
***REMOVED***

type FullResponse struct {
	Response
	Payload FullRecord
***REMOVED***

type HeaterRequest struct {
	Request
	Temp float64
***REMOVED***

type HeaterResponse struct {
	Response
	Payload HeaterRecord
***REMOVED***

type TempRecord struct {
	Temp float64
	Time time.Time
***REMOVED***

type ReplenisherRecord struct {
	IsReplenishing bool
	Time           time.Time
***REMOVED***

type FullRecord struct {
	IsFull bool
	Time   time.Time
***REMOVED***

type HeaterRecord struct {
	Duty   float64
	Period time.Duration
	Time   time.Time
***REMOVED***

func (r Request***REMOVED*** IsGet(***REMOVED*** bool {
	return r.Code == CodeGet
***REMOVED***

func (r Request***REMOVED*** IsPut(***REMOVED*** bool {
	return r.Code == CodePut
***REMOVED***

func (r Response***REMOVED*** IsFailure(***REMOVED*** bool {
	return r.Code != CodeSuccess
***REMOVED***
