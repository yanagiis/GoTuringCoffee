package vl6180x

import (
	"bytes"
	"encoding/binary"
	"log"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

const DEFAULT_ADDRESS = 0x29

type VL6180X struct {
	device *i2c.Dev
	bus    i2c.BusCloser
}

func (v *VL6180X) Init(i2cPath string) {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	bus, err := i2creg.Open(i2cPath)
	if err != nil {
		log.Fatal(err)
	}
	v.bus = bus

	// Dev is a valid conn.Conn.
	v.device = &i2c.Dev{Addr: DEFAULT_ADDRESS, Bus: bus}

	// Initialize device
	v.WriteByte(0x0207, 0x01)
	v.WriteByte(0x0208, 0x01)
	v.WriteByte(0x0096, 0x00)
	v.WriteByte(0x0097, 0xfd)
	v.WriteByte(0x00e3, 0x00)
	v.WriteByte(0x00e4, 0x04)
	v.WriteByte(0x00e5, 0x02)
	v.WriteByte(0x00e6, 0x01)
	v.WriteByte(0x00e7, 0x03)
	v.WriteByte(0x00f5, 0x02)
	v.WriteByte(0x00d9, 0x05)
	v.WriteByte(0x00db, 0xce)
	v.WriteByte(0x00dc, 0x03)
	v.WriteByte(0x00dd, 0xf8)
	v.WriteByte(0x009f, 0x00)
	v.WriteByte(0x00a3, 0x3c)
	v.WriteByte(0x00b7, 0x00)
	v.WriteByte(0x00bb, 0x3c)
	v.WriteByte(0x00b2, 0x09)
	v.WriteByte(0x00ca, 0x09)
	v.WriteByte(0x0198, 0x01)
	v.WriteByte(0x01b0, 0x17)
	v.WriteByte(0x01ad, 0x00)
	v.WriteByte(0x00ff, 0x05)
	v.WriteByte(0x0100, 0x05)
	v.WriteByte(0x0199, 0x05)
	v.WriteByte(0x01a6, 0x1b)
	v.WriteByte(0x01ac, 0x3e)
	v.WriteByte(0x01a7, 0x1f)
	v.WriteByte(0x0030, 0x00)

	// Recommended : Public registers - See data sheet for more detail

	v.WriteByte(0x0011, 0x10) // Enables polling for ‘New Sample ready’ when measurement completes
	v.WriteByte(0x010a, 0x30) // Set the averaging sample period (compromise between lower noise and increased execution time)
	v.WriteByte(0x003f, 0x46) // Sets the light and dark gain (upper nibble). Dark gain should not be changed.
	v.WriteByte(0x0031, 0xFF) // sets the # of range measurements after which auto calibration of system is performed
	v.WriteByte(0x0040, 0x63) // Set ALS integration time to 100ms
	v.WriteByte(0x002e, 0x01) // perform a single temperature calibratio of the ranging sensor
	v.WriteByte(0x001b, 0x09) // Set default ranging inter-measurement period to 100ms
	v.WriteByte(0x003e, 0x31) // Set default ALS inter-measurement period to 500ms
	v.WriteByte(0x0014, 0x24) // Configures interrupt on ‘New Sample Ready threshold event’

	v.WriteByte(0x016, 0x00)
}

func (v *VL6180X) Close() {
	if v.device != nil {
		v.bus.Close()
	}
}

func SetScaling(newScaling int){
    int scalerValues[] = {0, 253, 127, 84};
    int defaultCrosstalkValidHeight = 20;
    if (new_scaling < 1 || new_scaling > 3) { return; }

    int ptp_offset = read_byte(handle,0x24);

    write_two_bytes(handle,0x96,scalerValues[new_scaling]);
    write_byte(handle,0x24,ptp_offset / new_scaling);
    write_byte(handle,0x21, defaultCrosstalkValidHeight / new_scaling);
    int rce = read_byte(handle,0x2d);
    write_byte(handle,0x2d, (rce & 0xFE) | (new_scaling == 1));
}

func (v *VL6180X) WriteByte(reg int16, data byte) {
	v.WriteBytes(reg, []byte{data})
}

func (v *VL6180X) WriteBytes(reg int16, data []byte) {
	v.device.Write(v.Encode(reg, data))
}

func (v *VL6180X) Encode(reg int16, data []byte) []byte {
	bytesBuffer := new(bytes.Buffer)
	binary.Write(bytesBuffer, binary.LittleEndian, reg)
	binary.Write(bytesBuffer, binary.LittleEndian, data)

	return bytesBuffer.Bytes()
}
