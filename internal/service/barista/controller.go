package barista

***REMOVED***
	"bytes"
	"errors"
***REMOVED***
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
***REMOVED***

type Controller interface {
	Connect(***REMOVED*** error
	Disconnect(***REMOVED*** error
	Do(p *lib.Point***REMOVED*** error
***REMOVED***

type SEController struct {
	Smoothie *hardware.Smoothie
	Extruder *hardware.Extruder
***REMOVED***

func (se *SEController***REMOVED*** Connect(***REMOVED*** (err error***REMOVED*** {
	if err = se.Extruder.Connect(***REMOVED***; err != nil {
		return err
***REMOVED***
	defer func(***REMOVED*** {
	***REMOVED***
			se.Extruder.Disconnect(***REMOVED***
	***REMOVED***
***REMOVED***(***REMOVED***

	if err = se.Smoothie.Connect(***REMOVED***; err != nil {
		return err
***REMOVED***
	return nil
***REMOVED***

func (se *SEController***REMOVED*** Disconnect(***REMOVED*** error {
	se.Extruder.Disconnect(***REMOVED***
	se.Smoothie.Disconnect(***REMOVED***
	return nil
***REMOVED***

func (se *SEController***REMOVED*** Do(p *lib.Point***REMOVED*** error {
	gcode, gerr := pointToGCode(p***REMOVED***
	hcode, herr := pointToHCode(p***REMOVED***

	if gerr == nil {
		se.Smoothie.Writeline(gcode***REMOVED***
***REMOVED***
	if herr == nil {
		se.Extruder.Writeline(hcode***REMOVED***
***REMOVED***

	if gerr == nil {
		resp, err := se.Smoothie.Readline(***REMOVED***
		if strings.Compare(resp, "ok"***REMOVED*** != 0 {
			log.Error(***REMOVED***.Err(err***REMOVED***
	***REMOVED***
***REMOVED***
	if herr == nil {
		resp, err := se.Extruder.Readline(***REMOVED***
		if strings.Compare(resp, "ok"***REMOVED*** != 0 {
			log.Error(***REMOVED***.Err(err***REMOVED***
	***REMOVED***
***REMOVED***

	return nil
***REMOVED***

func pointToGCode(p *lib.Point***REMOVED*** (string, error***REMOVED*** {
	if p.X == nil && p.Y == nil && p.Z == nil {
		return "", errors.New("no x, y, and z"***REMOVED***
***REMOVED***

	var buffer bytes.Buffer
	buffer.WriteString("G1"***REMOVED***
	if p.X != nil {
		buffer.WriteString(fmt.Sprintf(" X%0.5f", *p.X***REMOVED******REMOVED***
***REMOVED***
	if p.Y != nil {
		buffer.WriteString(fmt.Sprintf(" Y%0.5f", *p.Y***REMOVED******REMOVED***
***REMOVED***
	if p.Z != nil {
		buffer.WriteString(fmt.Sprintf(" Z%0.5f", *p.Z***REMOVED******REMOVED***
***REMOVED***
	buffer.WriteString(fmt.Sprintf(" F%0.5f", *p.F***REMOVED******REMOVED***
	return buffer.String(***REMOVED***, nil
***REMOVED***

func pointToHCode(p *lib.Point***REMOVED*** (string, error***REMOVED*** {
	if p.Time == nil {
		return "", errors.New("no time"***REMOVED***
***REMOVED***

	var buffer bytes.Buffer
	buffer.WriteString("H"***REMOVED***
	if p.E1 != nil && *p.E1 != 0 {
		buffer.WriteString(fmt.Sprintf(" E0 %05f", *p.E1***REMOVED******REMOVED***
***REMOVED***
	if p.E2 != nil && *p.E2 != 0 {
		buffer.WriteString(fmt.Sprintf(" E1 %05f", *p.E2***REMOVED******REMOVED***
***REMOVED***
	buffer.WriteString(fmt.Sprintf(" T%0.5f", *p.Time***REMOVED******REMOVED***
	return buffer.String(***REMOVED***, nil
***REMOVED***
