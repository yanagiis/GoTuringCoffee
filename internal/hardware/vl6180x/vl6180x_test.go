package vl6180x

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple testing what different between Fatal and Error
func TestNew(t *testing.T) {
	var device VL6180X
	result := device.WriteBytes(0x0207, []byte{0x01})
	assert.Equal(t, len(result), 3, "The binary length should be 3")
	assert.Equal(t, result[0], uint8(7), "The binary length should be 3")
	assert.Equal(t, result[1], uint8(2), "The binary length should be 3")
	assert.Equal(t, result[2], uint8(1), "The binary length should be 3")
}
