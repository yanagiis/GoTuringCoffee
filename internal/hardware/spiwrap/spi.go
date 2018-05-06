//go:generate go-enum -fspi.go

package spiwrap

***REMOVED***
	"errors"
***REMOVED***
	"io"

	"github.com/yanagiis/periph/conn/spi"
	"github.com/yanagiis/periph/conn/spi/spireg"
***REMOVED***

type SPI interface {
	io.Closer
	Open(***REMOVED*** error
	Tx(w, r []byte***REMOVED*** error
***REMOVED***

type Config struct {
	Path  string   `mapstructure:"path"`
	Speed int64    `mapstructure:"speed"`
	Mode  spi.Mode `mapstructure:"mode"`
	Bits  int      `mapstructure:"bits"`
***REMOVED***

type SPIDevice struct {
	Conf   Config
	closer spi.PortCloser
	conn   spi.Conn
***REMOVED***

func (s *SPIDevice***REMOVED*** Open(***REMOVED*** error {
	var err error
	if s.closer == nil {
		s.closer, err = spireg.Open(s.Conf.Path***REMOVED***
	***REMOVED***
			return err
	***REMOVED***
***REMOVED***
	if s.conn == nil {
		if s.conn, err = s.closer.Connect(s.Conf.Speed, s.Conf.Mode, s.Conf.Bits***REMOVED***; err != nil {
			return err
	***REMOVED***
***REMOVED***
	return nil
***REMOVED***

func (s *SPIDevice***REMOVED*** IsOpen(***REMOVED*** bool {
	return s.conn != nil
***REMOVED***

func (s *SPIDevice***REMOVED*** Close(***REMOVED*** error {
	err := s.closer.Close(***REMOVED***
***REMOVED***
		return err
***REMOVED***
	s.conn = nil
	s.closer = nil
	return nil
***REMOVED***

func (s *SPIDevice***REMOVED*** Tx(w, r []byte***REMOVED*** error {
	if s.conn == nil {
		fmt.Errorf("Not open"***REMOVED***
		return errors.New("Not open"***REMOVED***
***REMOVED***
	return s.conn.Tx(w, r***REMOVED***
***REMOVED***
