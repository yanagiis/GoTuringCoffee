package hardware

***REMOVED***
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
***REMOVED***

type ExtruderPort interface {
	io.ReadWriter
	io.Closer
	Open(***REMOVED*** error
	IsOpen(***REMOVED*** bool
***REMOVED***

// Extruder setting
type Extruder struct {
	port ExtruderPort
	io   *bufio.ReadWriter
***REMOVED***

// ConnectExtruder is used to check extruder is alive or not.
func NewExtruder(port ExtruderPort***REMOVED*** *Extruder {
	return &Extruder{
		port: port,
		io:   bufio.NewReadWriter(bufio.NewReader(port***REMOVED***, bufio.NewWriter(port***REMOVED******REMOVED***,
***REMOVED***
***REMOVED***

func (e *Extruder***REMOVED*** Connect(***REMOVED*** error {
	if err := e.port.Open(***REMOVED***; err != nil {
		return err
***REMOVED***

	e.io.Flush(***REMOVED***
	if e.Writeline(""***REMOVED*** == nil {
		line, err := e.Readline(***REMOVED***
	***REMOVED***
			return err
	***REMOVED***
		if strings.Compare(line, "Ok"***REMOVED*** == 0 {
			return nil
	***REMOVED***
***REMOVED***
	return errors.New("no response"***REMOVED***
***REMOVED***

// Disconnect extruder
func (e *Extruder***REMOVED*** Disconnect(***REMOVED*** error {
	if err := e.io.Flush(***REMOVED***; err != nil {
		return err
***REMOVED***
	if err := e.port.Close(***REMOVED***; err != nil {
		return err
***REMOVED***
	return nil
***REMOVED***

// Writeline is used to write a line to extruder
func (e *Extruder***REMOVED*** Writeline(msg string***REMOVED*** error {
	var buffer bytes.Buffer
	var err error

	if _, err = buffer.WriteString(msg***REMOVED***; err != nil {
		return err
***REMOVED***
	buffer.WriteByte('\n'***REMOVED***

	if _, err := buffer.WriteTo(e.io***REMOVED***; err != nil {
		return err
***REMOVED***
	return nil
***REMOVED***

// Readline is used to read a line from extruder
func (e *Extruder***REMOVED*** Readline(***REMOVED*** (string, error***REMOVED*** {
	line, isPrefix, err := e.io.ReadLine(***REMOVED***
***REMOVED***
		return "", err
***REMOVED***
	if isPrefix {
		return "", errors.New("Line is too long"***REMOVED***
***REMOVED***
	return string(line***REMOVED***, err
***REMOVED***
