package hardware

***REMOVED***
	"errors"
***REMOVED***

// TC Type
const (
	TypeB byte = 0x0
	TypeE      = 0x1
	TypeJ      = 0x2
	TypeK      = 0x3
	TypeN      = 0x4
	TypeR      = 0x5
	TypeS      = 0x6
	TypeT      = 0x7
***REMOVED***

// MAX31856 mode
const (
	MAX31856ModeManual    byte = 0x0
	MAX31856ModeAutomatic      = 0x1
***REMOVED***

// Sample average
const (
	SampleAvg1  byte = 0x0
	SampleAvg2       = 0x1
	SampleAvg4       = 0x2
	SampleAvg8       = 0x3
	SampleAvg16      = 0x4
***REMOVED***

const (
	addrWriteMAX31856 byte = 0x80
	addrCR0                = 0x0
	addrCR1                = 0x1
	addrMASK               = 0x2
	addrCJHF               = 0x3
	addrCJLF               = 0x4
	addrLTHFTH             = 0x5
	addrLTHFTL             = 0x6
	addrLTLFTH             = 0x7
	addrLTLFTL             = 0x8
	addrCJTO               = 0x9
	addrCJTH               = 0xA
	addrCJTL               = 0xB
	addrLTCBH              = 0xC
	addrLTCBM              = 0xD
	addrLTCBL              = 0xE
	addrSR                 = 0xF
***REMOVED***

const (
	resolutionTC = 0.0078125
	resolutionCJ = 0.015625
***REMOVED***

// MAX31856 TC sensor
type MAX31856 struct {
	spi       SPI
	conf      MAX31856Config
	connected bool
***REMOVED***

// MAX31856Config is used to setting max31856 sensor
type MAX31856Config struct {
	TC   byte `mapstructure:"tc"`
	Avg  byte `mapstructure:"avg"`
	Mode byte `mapstructure:"mode"`
***REMOVED***

func NewMAX31856(spi SPI, conf MAX31856Config***REMOVED*** *MAX31856 {
	return &MAX31856{
		spi:       spi,
		conf:      conf,
		connected: false,
***REMOVED***
***REMOVED***

func (m *MAX31856***REMOVED*** Connect(***REMOVED*** error {
	var err error
	var value byte

	if m.connected {
		return nil
***REMOVED***

	if err = m.spi.Open(***REMOVED***; err != nil {
		return err
***REMOVED***

	if err = m.setMode(m.conf.Mode***REMOVED***; err != nil {
		return err
***REMOVED***

	if err = m.setSampleAvg(m.conf.Avg***REMOVED***; err != nil {
		return err
***REMOVED***

	if err = m.setTCType(m.conf.TC***REMOVED***; err != nil {
		return err
***REMOVED***

	value, err = m.getMode(***REMOVED***
***REMOVED***
		return err
***REMOVED***
	if value != m.conf.Mode {
		return errors.New("set mode failed"***REMOVED***
***REMOVED***

	value, err = m.getTCType(***REMOVED***
***REMOVED***
		return err
***REMOVED***
	if value != m.conf.TC {
		return errors.New("set tc-type failed"***REMOVED***
***REMOVED***

	value, err = m.getSampleAvg(***REMOVED***
***REMOVED***
		return err
***REMOVED***
	if value != m.conf.Avg {
		return errors.New("set sample avg failed"***REMOVED***
***REMOVED***

	m.connected = true
	return nil
***REMOVED***

func (m *MAX31856***REMOVED*** Disconnect(***REMOVED*** error {
	if err := m.spi.Close(***REMOVED***; err != nil {
		return err
***REMOVED***
	m.connected = false
	return nil
***REMOVED***

// GetTemperature from max31856
func (m *MAX31856***REMOVED*** GetTemperature(***REMOVED*** (float64, error***REMOVED*** {
	if !m.connected {
		return 0, errors.New("sensor is connected yet"***REMOVED***
***REMOVED***

	buf := make([]byte, 4***REMOVED*** // [t2, t1, t0, fault]
	if err := m.readReg(addrLTCBH, buf***REMOVED***; err != nil {
		return 0, err
***REMOVED***

	adcValue := ((int32(buf[0]***REMOVED*** << 16***REMOVED*** | (int32(buf[1]***REMOVED*** << 8***REMOVED*** | int32(buf[2]***REMOVED******REMOVED*** >> 5

	if (buf[0] & 0x80***REMOVED*** != 0 {
		adcValue -= 0x80000
***REMOVED***

	if adcValue == 0 {
		return 0, errors.New("get zero ADC Value"***REMOVED***
***REMOVED***

	return float64(adcValue***REMOVED*** * resolutionTC, nil
***REMOVED***

func (m *MAX31856***REMOVED*** getMode(***REMOVED*** (byte, error***REMOVED*** {
	buf := make([]byte, 1***REMOVED***
	err := m.readReg(addrCR0, buf***REMOVED***
	return buf[0] >> 7, err
***REMOVED***

func (m *MAX31856***REMOVED*** setMode(mode byte***REMOVED*** error {
	buf := make([]byte, 1***REMOVED***
	if err := m.readReg(addrCR0, buf***REMOVED***; err != nil {
		return err
***REMOVED***
	buf[0] = (buf[0] & 0x7f***REMOVED*** | (mode << 7***REMOVED***
	return m.writeReg(addrCR0, buf***REMOVED***
***REMOVED***

func (m *MAX31856***REMOVED*** getSampleAvg(***REMOVED*** (byte, error***REMOVED*** {
	buf := make([]byte, 1***REMOVED***
	err := m.readReg(addrCR1, buf***REMOVED***
	return (buf[0] % 0x70***REMOVED*** >> 4, err
***REMOVED***

func (m *MAX31856***REMOVED*** setSampleAvg(avg byte***REMOVED*** error {
	buf := make([]byte, 1***REMOVED***
	if err := m.readReg(addrCR0, buf***REMOVED***; err != nil {
		return err
***REMOVED***
	buf[0] = (buf[0] & 0x8f***REMOVED*** | (avg << 4***REMOVED***
	return m.writeReg(addrCR0, buf***REMOVED***
***REMOVED***

func (m *MAX31856***REMOVED*** getTCType(***REMOVED*** (byte, error***REMOVED*** {
	buf := make([]byte, 1***REMOVED***
	err := m.readReg(addrCR1, buf***REMOVED***
	return (buf[0] & 0x0f***REMOVED***, err
***REMOVED***

func (m *MAX31856***REMOVED*** setTCType(t byte***REMOVED*** error {
	buf := make([]byte, 1***REMOVED***
	if err := m.readReg(addrCR1, buf***REMOVED***; err != nil {
		return err
***REMOVED***
	buf[0] = (buf[0] & 0xf0***REMOVED*** | t
	return m.writeReg(addrCR1, buf***REMOVED***
***REMOVED***

func (m *MAX31856***REMOVED*** readReg(addr byte, buf []byte***REMOVED*** error {
	rbuf := append([]byte{addr***REMOVED***, buf...***REMOVED***
	wbuf := make([]byte, len(rbuf***REMOVED******REMOVED***
	return m.spi.Tx(wbuf, rbuf***REMOVED***
***REMOVED***

func (m *MAX31856***REMOVED*** writeReg(addr byte, buf []byte***REMOVED*** error {
	wbuf := append([]byte{addr | addrWriteMAX31856***REMOVED***, buf...***REMOVED***
	rbuf := make([]byte, len(wbuf***REMOVED******REMOVED***
	return m.spi.Tx(wbuf, rbuf***REMOVED***
***REMOVED***
