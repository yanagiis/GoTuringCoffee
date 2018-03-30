package hardware

***REMOVED***
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

// SPI is SPI's configuration
type SPIDevice struct {
	Path   string   `mapstructure:"path"`
	Hz     int64    `mapstructure:"hz"`
	Mode   spi.Mode `mapstructure:"mode"`
	Bits   int      `mapstructure:"bits"`
	closer spi.PortCloser
	conn   spi.Conn
***REMOVED***

func (s *SPIDevice***REMOVED*** Open(***REMOVED*** error {
	var err error
	if s.closer == nil {
		s.closer, err = spireg.Open(s.Path***REMOVED***
	***REMOVED***
			return err
	***REMOVED***
***REMOVED***
	if s.conn == nil {
		if s.conn, err = s.closer.Connect(s.Hz, s.Mode, s.Bits***REMOVED***; err != nil {
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
		fmt.Errorf("Not opened"***REMOVED***
***REMOVED***
	return s.conn.Tx(w, r***REMOVED***
***REMOVED***
