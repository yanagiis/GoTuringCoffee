//go:generate go-enum -fmax31865.go

package max31865

***REMOVED***
	"errors"
***REMOVED***

	"github.com/yanagiis/GoTuringCoffee/internal/hardware/spiwrap"
***REMOVED***

// RTD wire type
// ENUM(
// 2, 3, 4
// ***REMOVED***
type Wire byte

// MAX31865 Mode
// ENUM(
// Manual, Automatic
// ***REMOVED***
type Mode byte

const (
	addrWriteMAX31865 byte = 0x80
	addrCR                 = 0x0
	addrRTDH               = 0x1
	addrRTDL               = 0x2
	addrHFTH               = 0x3
	addrHFTL               = 0x4
	addrLFTH               = 0x5
	addrLFTL               = 0x6
	addrFault              = 0x7
***REMOVED***

// Config is used to setting max31865 sensor
type Config struct {
	Wire Wire `mapstructure:"wire"`
	Mode Mode `mapstructure:"mode"`
***REMOVED***

// MAX31865 RTD sensor
type MAX31865 struct {
	spi       spiwrap.SPI
	conf      Config
	connected bool
***REMOVED***

func NewMAX31865(spi spiwrap.SPI, conf Config***REMOVED*** *MAX31865 {
	return &MAX31865{
		spi:       spi,
		conf:      conf,
		connected: false,
***REMOVED***
***REMOVED***

func (m *MAX31865***REMOVED*** Connect(***REMOVED*** error {
	var err error
	var mode Mode
	var wire Wire

	if m.connected {
		return nil
***REMOVED***

	if err = m.spi.Open(***REMOVED***; err != nil {
		return err
***REMOVED***

	err = m.setMode(m.conf.Mode***REMOVED***
***REMOVED***
		return err
***REMOVED***
	err = m.setWire(m.conf.Wire***REMOVED***
***REMOVED***
		return err
***REMOVED***

	mode, err = m.getMode(***REMOVED***
***REMOVED***
		return err
***REMOVED***
	if mode != m.conf.Mode {
		return errors.New("set mode failed"***REMOVED***
***REMOVED***

	wire, err = m.getWire(***REMOVED***
***REMOVED***
		return err
***REMOVED***
	if wire != m.conf.Wire {
		return errors.New("set wire failed"***REMOVED***
***REMOVED***

	buf := make([]byte, 1***REMOVED***
	err = m.readReg(addrCR, buf***REMOVED***
***REMOVED***
		return err
***REMOVED***

	buf[0] = (buf[0] & 0x7f***REMOVED*** | (1 << 7***REMOVED***
	err = m.writeReg(addrCR, buf***REMOVED***
***REMOVED***
		return err
***REMOVED***

	m.connected = true
	return nil
***REMOVED***

func (m *MAX31865***REMOVED*** Disconnect(***REMOVED*** error {
	if err := m.spi.Close(***REMOVED***; err != nil {
		return err
***REMOVED***
	m.connected = false
	return nil
***REMOVED***

// GetTemperature get temperature
//
// Callendar-Van Dusen equation:
// R(T***REMOVED*** = R0(1 + aT + bT2 + c(T - 100***REMOVED***T3***REMOVED***
// where:
//     T = temperature (C***REMOVED***
//     R(T***REMOVED*** = resistance at T
//     R0 = resistance at T = 0C
//     IEC 751 specifies α = 0.00385055 and the following
//     Callendar-Van Dusen coefficient values:
//         a = 3.90830 x 10^-3
//         b = -5.77500 x 10^-7
//         c = -4.18301 x 10^-12 for -200C < T < 0C, 0 for 0C < T < +850C
//
// Linearizing Temperature Data
// For a temperature range of -100C to +100C, a good
// approximation of temperature can be made by simply
// using the RTD data as shown below:
// Temperature (C***REMOVED*** ≈ (ADC code/32***REMOVED*** – 256
func (m *MAX31865***REMOVED*** GetTemperature(***REMOVED*** (float64, error***REMOVED*** {

	if !m.connected {
		return 0, errors.New("sensor is not connected yet"***REMOVED***
***REMOVED***

	buf := make([]byte, 2***REMOVED***
	err := m.readReg(addrRTDH, buf***REMOVED***
***REMOVED***
		return 0, err
***REMOVED***

	if (buf[1] & 0x1***REMOVED*** != 0 {
		var value byte
		value, err = m.getError(***REMOVED***
	***REMOVED***
			return 0, err
	***REMOVED***

		err = m.clearError(***REMOVED***
	***REMOVED***
			return 0, err
	***REMOVED***

		return 0, fmt.Errorf("max31865 error %02x", value***REMOVED***
***REMOVED***

	adcValue := ((int32(buf[0]***REMOVED*** << 8***REMOVED*** | (int32(buf[1]***REMOVED******REMOVED******REMOVED*** >> 1
	if adcValue == 0 {
		return 0, errors.New("get zero adcValue"***REMOVED***
***REMOVED***

	return float64(adcValue***REMOVED***/32 - 256, nil
***REMOVED***

func (m *MAX31865***REMOVED*** getMode(***REMOVED*** (Mode, error***REMOVED*** {
	buf := make([]byte, 1***REMOVED***
	err := m.readReg(addrCR, buf***REMOVED***
	return Mode((buf[0] & 0x40***REMOVED*** >> 6***REMOVED***, err
***REMOVED***

func (m *MAX31865***REMOVED*** setMode(mode Mode***REMOVED*** error {
	buf := make([]byte, 1***REMOVED***
	err := m.readReg(addrCR, buf***REMOVED***
***REMOVED***
		return err
***REMOVED***
	buf[0] = (buf[0] & 0xbf***REMOVED*** | (byte(mode***REMOVED*** << 6***REMOVED***
	return m.writeReg(addrCR, buf***REMOVED***
***REMOVED***

func (m *MAX31865***REMOVED*** getWire(***REMOVED*** (Wire, error***REMOVED*** {
	buf := make([]byte, 1***REMOVED***
	err := m.readReg(addrCR, buf***REMOVED***
	return Wire((buf[0] & 0x10***REMOVED*** >> 4***REMOVED***, err
***REMOVED***

func (m *MAX31865***REMOVED*** setWire(wire Wire***REMOVED*** error {
	buf := make([]byte, 1***REMOVED***
	err := m.readReg(addrCR, buf***REMOVED***
***REMOVED***
		return err
***REMOVED***
	buf[0] = (buf[0] & 0xef***REMOVED*** | (byte(wire***REMOVED*** << 4***REMOVED***
	return m.writeReg(addrCR, buf***REMOVED***
***REMOVED***

func (m *MAX31865***REMOVED*** getError(***REMOVED*** (byte, error***REMOVED*** {
	buf := make([]byte, 1***REMOVED***
	err := m.readReg(addrFault, buf***REMOVED***
	return buf[0], err
***REMOVED***

func (m *MAX31865***REMOVED*** clearError(***REMOVED*** error {
	buf := make([]byte, 1***REMOVED***
	err := m.readReg(addrCR, buf***REMOVED***
***REMOVED***
		return err
***REMOVED***
	buf[0] |= 0x2
	return m.writeReg(addrCR, buf***REMOVED***
***REMOVED***

func (m *MAX31865***REMOVED*** readReg(addr byte, buf []byte***REMOVED*** error {
	wbuf := append([]byte{addr***REMOVED***, buf...***REMOVED***
	rbuf := make([]byte, len(wbuf***REMOVED******REMOVED***
	if err := m.spi.Tx(wbuf, rbuf***REMOVED***; err != nil {
		return err
***REMOVED***
	copy(buf, rbuf[1:]***REMOVED***
	return nil
***REMOVED***

func (m *MAX31865***REMOVED*** writeReg(addr byte, buf []byte***REMOVED*** error {
	wbuf := append([]byte{addr | addrWriteMAX31865***REMOVED***, buf...***REMOVED***
	rbuf := make([]byte, len(wbuf***REMOVED******REMOVED***
	return m.spi.Tx(wbuf, rbuf***REMOVED***
***REMOVED***
