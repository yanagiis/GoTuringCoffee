package hardware

***REMOVED***
***REMOVED***
	"net"
***REMOVED***

type TCP struct {
	Addr string `mapstructure:"addr"`
	Port uint16 `mapstructure:"port"`
	conn net.Conn
***REMOVED***

func (t *TCP***REMOVED*** Open(***REMOVED*** error {
	var err error
	if t.conn, err = net.Dial("tcp", fmt.Sprint("%s:%u", t.Addr, t.Port***REMOVED******REMOVED***; err != nil {
		return err
***REMOVED***
	return nil
***REMOVED***

func (t *TCP***REMOVED*** IsOpen(***REMOVED*** bool {
	return t.conn != nil
***REMOVED***

func (t *TCP***REMOVED*** Close(***REMOVED*** error {
	if t.conn != nil {
		if err := t.conn.Close(***REMOVED***; err != nil {
			return err
	***REMOVED***
		t.conn = nil
***REMOVED***
	return nil
***REMOVED***

func (t *TCP***REMOVED*** Read(p []byte***REMOVED*** (int, error***REMOVED*** {
	return t.conn.Read(p***REMOVED***
***REMOVED***

func (t *TCP***REMOVED*** Write(p []byte***REMOVED*** (int, error***REMOVED*** {
	return t.conn.Write(p***REMOVED***
***REMOVED***
