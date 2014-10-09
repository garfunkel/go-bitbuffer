package bitbuffer

import (
	"encoding/binary"
	"bytes"
)

type BitBuffer struct {
	buffer []byte
	pos    uint64
}

func NewBitBuffer() (bitBuffer *BitBuffer) {
	return &BitBuffer{
		pos: 0,
	}
}

func (bitBuffer *BitBuffer) Feed(data []byte) {
	bitBuffer.buffer = append(bitBuffer.buffer, data...)

	return
}

func (bitBuffer *BitBuffer) Read(numBits uint64) (data []byte, err error) {
	for numBits > 0 {
		data = append(data, bitBuffer.buffer[0])
		data[len(data)-1] <<= bitBuffer.pos
		shifter := bitBuffer.buffer[1] >> (8 - bitBuffer.pos)
		data[len(data)-1] ^= shifter

		if numBits < 8 {
			data[len(data)-1] >>= (8 - numBits)
			data[len(data)-1] <<= (8 - numBits)

			if bitBuffer.pos+numBits > 7 {
				bitBuffer.buffer = bitBuffer.buffer[1:]
			}

			bitBuffer.pos = (bitBuffer.pos + numBits) % 8
			numBits = 0

			break
		} else {
			bitBuffer.buffer = bitBuffer.buffer[1:]
		}

		numBits -= 8
	}

	return
}

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

func (bitBuffer *BitBuffer) ReadString(numBits uint64) (data string, err error) {
	dataBytes, err := bitBuffer.Read(numBits)

	if err != nil {
		return
	}

	data = string(dataBytes)

	return
}

/*
func main() {
	buffer := NewBitBuffer()

	// 01100110 01101100 01100001 01100011 00100000 00110101 00110001 00110010
	b := []byte("flac 512")

	buffer.Feed(b)

	// Read ints.
	intData, err := buffer.Read(11)
	fmt.Println(strconv.FormatUint(intData, 2))
	intData, err = buffer.Read(7)
	fmt.Println(strconv.FormatUint(intData, 2))

	buffer = NewBitBuffer()

	b = []byte("flac 512")

	buffer.Feed(b)

	stringData, err := buffer.ReadString(11)

	fmt.Println(stringData, err)

	fmt.Println(err)
}*/
