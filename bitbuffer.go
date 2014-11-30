// Package bitbuffer is a simple and easy library used to read data on the bit-level from a buffer.
package bitbuffer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// BitBuffer represents a buffer, which is filled with bytes where each bit can be read as a single unit.
type BitBuffer struct {
	buffer    []byte
	pos       uint8
	ByteOrder binary.ByteOrder
}

// TooManyBitsError occurs when attempting to read too many bits from a buffer.
type TooManyBitsError uint8

func (err TooManyBitsError) Error() string {
	return fmt.Sprintf("bitbuffer: too many bits requested: %v", err)
}

// NewBitBuffer constructs a new BitBuffer.
func NewBitBuffer(byteOrder binary.ByteOrder) (bitBuffer *BitBuffer) {
	return &BitBuffer{
		pos:       0,
		ByteOrder: byteOrder,
	}
}

// Feed data bytes into the buffer.
func (bitBuffer *BitBuffer) Feed(data []byte) {
	bitBuffer.buffer = append(bitBuffer.buffer, data...)

	return
}

// Clear buffer.
func (bitBuffer *BitBuffer) Clear() {
	bitBuffer.buffer = []byte{}
	bitBuffer.pos = 0
}

// Read a number of bits from the buffer and return them as a byte array.
func (bitBuffer *BitBuffer) Read(numBits uint64) (data []byte, err error) {
	if uint64(len(bitBuffer.buffer)*8-int(bitBuffer.pos)) < numBits {
		err = io.EOF
	}

	for numBits > 0 && len(bitBuffer.buffer) > 0 {
		data = append(data, bitBuffer.buffer[0])
		data[len(data)-1] <<= bitBuffer.pos

		if len(bitBuffer.buffer) > 1 {
			shifter := bitBuffer.buffer[1] >> (8 - bitBuffer.pos)
			data[len(data)-1] ^= shifter
		}

		if numBits < 8 {
			data[len(data)-1] >>= (8 - numBits)
			data[len(data)-1] <<= (8 - numBits)

			if uint64(bitBuffer.pos)+numBits > 7 {
				bitBuffer.buffer = bitBuffer.buffer[1:]
			}

			bitBuffer.pos = uint8((uint64(bitBuffer.pos) + numBits) % 8)
			numBits = 0

			return
		}

		bitBuffer.buffer = bitBuffer.buffer[1:]
		numBits -= 8
	}

	return
}

// ReadUint64 reads a uint64 from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint64(numBits uint8) (data uint64, err error) {
	if numBits > 64 {
		err = TooManyBitsError(numBits)

		return
	}

	dataBytes, err := bitBuffer.Read(uint64(numBits))

	if err != nil {
		return
	}

	if bitBuffer.ByteOrder == binary.BigEndian {
		dataBytes = append(make([]byte, 8-len(dataBytes)), dataBytes...)
	} else if bitBuffer.ByteOrder == binary.LittleEndian {
		dataBytes = append(dataBytes, make([]byte, 8-len(dataBytes))...)
	}

	err = binary.Read(bytes.NewBuffer(dataBytes), bitBuffer.ByteOrder, &data)

	if err != nil {
		return
	}

	shifter := uint8(0)

	if numBits%8 > 0 {
		shifter = 8 - (numBits % 8)
	}

	data >>= shifter

	return
}

// ReadUint a uint from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint(numBits uint8) (data uint, err error) {
	wordSize := 32 << (^uint(0) >> 32 & 1)

	if numBits > uint8(wordSize) {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = uint(rawData)

	return
}

// ReadBit reads a single bit as a boolean and returns the value.
func (bitBuffer *BitBuffer) ReadBit() (data bool, err error) {
	rawData, err := bitBuffer.ReadUint8(1)

	if err != nil {
		return
	}

	data = rawData != 0

	return
}

// ReadUint8 reads a uint8 from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint8(numBits uint8) (data uint8, err error) {
	if numBits > 8 {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = uint8(rawData)

	return
}

// ReadUint16 reads a uint16 from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint16(numBits uint8) (data uint16, err error) {
	if numBits > 16 {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = uint16(rawData)

	return
}

// ReadUint32 reads a uint32 from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint32(numBits uint8) (data uint32, err error) {
	if numBits > 32 {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = uint32(rawData)

	return
}

// ReadByte reads a byte from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadByte(numBits uint8) (data byte, err error) {
	if numBits > 8 {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = byte(rawData)

	return
}

// ReadString reads numBits bits from the buffer and returns the value as a string.
func (bitBuffer *BitBuffer) ReadString(numBits uint64) (data string, err error) {
	dataBytes, err := bitBuffer.Read(numBits)

	if err != nil {
		return
	}

	data = string(dataBytes)

	return
}
