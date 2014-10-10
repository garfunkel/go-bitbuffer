package bitbuffer

import (
	"encoding/binary"
	"bytes"
	"io"
)

// BitBuffer represents a buffer, which is filled with bytes where each bit can be read as a single unit.
type BitBuffer struct {
	buffer []byte
	pos    uint8
}

// NewBitBuffer constructs a new BitBuffer.
func NewBitBuffer() (bitBuffer *BitBuffer) {
	return &BitBuffer{
		pos: 0,
	}
}

// Feed data bytes into the buffer.
func (bitBuffer *BitBuffer) Feed(data []byte) {
	bitBuffer.buffer = append(bitBuffer.buffer, data...)

	return
}

// Read a number of bits from the buffer and return them as a byte array.
func (bitBuffer *BitBuffer) Read(numBits uint64) (data []byte, err error) {
	if uint64(len(bitBuffer.buffer) * 8 - int(bitBuffer.pos)) < numBits {
		err = io.EOF
	}

	for numBits > 0 && len(bitBuffer.buffer) > 0 {
		data = append(data, bitBuffer.buffer[0])
		data[len(data) - 1] <<= bitBuffer.pos

		if len(bitBuffer.buffer) > 1 {
			shifter := bitBuffer.buffer[1] >> (8 - bitBuffer.pos)
			data[len(data) - 1] ^= shifter
		}

		if numBits < 8 {
			data[len(data) - 1] >>= (8 - numBits)
			data[len(data) - 1] <<= (8 - numBits)

			if uint64(bitBuffer.pos) + numBits > 7 {
				bitBuffer.buffer = bitBuffer.buffer[1 :]
			}

			bitBuffer.pos = uint8((uint64(bitBuffer.pos) + numBits) % 8)
			numBits = 0

			return
		} else {
			bitBuffer.buffer = bitBuffer.buffer[1 :]
		}

		numBits -= 8
	}

	return
}

// Read an unsigned integer from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint64(numBits uint8) (data uint64, err error) {
	dataBytes, err := bitBuffer.Read(uint64(numBits))

	if err != nil {
		return
	}

	dataBytes = append(make([]byte, 8 - len(dataBytes)), dataBytes...)
	err = binary.Read(bytes.NewBuffer(dataBytes), binary.BigEndian, &data)

	if err != nil {
		return
	}

	shifter := uint8(0)

	if numBits % 8 > 0 {
		shifter = 8 - (numBits % 8)
	}

	data >>= shifter

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
