package hardware

***REMOVED***
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
***REMOVED***

var initCmds = [...]string{"G28", "G21", "G90", "M83"***REMOVED***

type SmoothiePort interface {
	io.ReadWriter
	io.Closer
	Open(***REMOVED*** error
	IsOpen(***REMOVED*** bool
***REMOVED***

// Smoothie setting
type Smoothie struct {
	port SmoothiePort
	io   *bufio.ReadWriter
***REMOVED***

// ConnectSmoothie is used to check extruder is alive or not.
func NewSmoothie(port SmoothiePort***REMOVED*** *Smoothie {
	return &Smoothie{
		port: port,
		io:   bufio.NewReadWriter(bufio.NewReaderSize(port, 60***REMOVED***, bufio.NewWriter(port***REMOVED******REMOVED***,
***REMOVED***
***REMOVED***

func (s *Smoothie***REMOVED*** Connect(***REMOVED*** error {
	if err := s.port.Open(***REMOVED***; err != nil {
		return err
***REMOVED***

	s.io.Flush(***REMOVED***
	if s.Writeline("G"***REMOVED*** {
		line, err := s.Readline(***REMOVED***
	***REMOVED***
			return err
	***REMOVED***
		if strings.Compare(line, "Ok"***REMOVED*** == 0 {
			return nil
	***REMOVED***
***REMOVED***

	for _, cmd := range initCmds {
		var line string
		var err error

		if s.Writeline(cmd***REMOVED*** {
			return errors.New("initial failed"***REMOVED***
	***REMOVED***
		if line, err = s.Readline(***REMOVED***; err != nil {
			return err
	***REMOVED***
		if strings.Compare(line, "ok"***REMOVED*** != 0 {
			return errors.New("initial failed"***REMOVED***
	***REMOVED***
***REMOVED***

	return errors.New("no response"***REMOVED***
***REMOVED***

// Disconnect extruder
func (s *Smoothie***REMOVED*** Disconnect(***REMOVED*** error {
	s.io.Flush(***REMOVED***
	if err := s.port.Close(***REMOVED***; err != nil {
		return err
***REMOVED***
	return nil
***REMOVED***

// Writeline is used to write a line to extruder
func (s *Smoothie***REMOVED*** Writeline(msg string***REMOVED*** bool {
	var buffer bytes.Buffer

	buffer.WriteString(msg***REMOVED***
	buffer.WriteByte('\n'***REMOVED***

	if _, err := buffer.WriteTo(s.io***REMOVED***; err != nil {
		return false
***REMOVED***
	return true
***REMOVED***

// Readline is used to read a line from extruder
func (s *Smoothie***REMOVED*** Readline(***REMOVED*** (string, error***REMOVED*** {
	line, isPrefix, err := s.io.ReadLine(***REMOVED***
***REMOVED***
		return "", err
***REMOVED***
	if isPrefix {
		return "", errors.New("Line is too long"***REMOVED***
***REMOVED***
	return string(line***REMOVED***, err
***REMOVED***
