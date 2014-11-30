package bitbuffer

import (
	"encoding/binary"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BitBufferTestSuite struct {
	suite.Suite
	leBitBuffer *BitBuffer
	beBitBuffer *BitBuffer
	data1 []byte
	data2 []byte
	dataStr []byte
	assert *assert.Assertions
}

func (suite *BitBufferTestSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
	suite.leBitBuffer = NewBitBuffer(binary.LittleEndian)
	suite.beBitBuffer = NewBitBuffer(binary.BigEndian)
	suite.data1 = []byte{'\x0F', '\xAC', '\x2B'}
	suite.data2 = []byte{'\x92', '\xDB'}
	suite.dataStr = []byte{'a', 'k', 'G'}

	suite.leBitBuffer.Feed(suite.data1)
	suite.leBitBuffer.Feed(suite.data2)
	suite.leBitBuffer.Feed(suite.dataStr)

	suite.beBitBuffer.Feed(suite.data2)
	suite.beBitBuffer.Feed(suite.data1)
	suite.beBitBuffer.Feed(suite.dataStr)
}

func (suite *BitBufferTestSuite) TestFeed() {
	suite.assert.Equal(append(append(suite.data1, suite.data2...), suite.dataStr...), suite.leBitBuffer.buffer)
	suite.assert.Equal(append(append(suite.data2, suite.data1...), suite.dataStr...), suite.beBitBuffer.buffer)
}

func (suite *BitBufferTestSuite) TestClear() {
	suite.leBitBuffer.Clear()
	suite.beBitBuffer.Clear()

	suite.assert.Empty(suite.leBitBuffer.buffer)
	suite.assert.Empty(suite.leBitBuffer.pos)

	suite.assert.Empty(suite.beBitBuffer.buffer)
	suite.assert.Empty(suite.beBitBuffer.pos)
}

func (suite *BitBufferTestSuite) TestRead() {
	leData, err := suite.leBitBuffer.Read(26)

	if err != nil {
		return
	}

	correct := []byte{suite.data1[0], suite.data1[1], suite.data1[2],
		suite.data2[0] >> 6 << 6}

	suite.assert.Equal(correct, leData)
	suite.assert.Equal(2, suite.leBitBuffer.pos)

	beData, err := suite.beBitBuffer.Read(15)

	if err != nil {
		return
	}

	correct = []byte{suite.data2[0], suite.data2[1] >> 1 << 1}

	suite.assert.Equal(correct, beData)
	suite.assert.Equal(7, suite.beBitBuffer.pos)
}

func (suite *BitBufferTestSuite) TestReadUint64() {
	data, err := suite.leBitBuffer.ReadUint64(32)

	if err != nil {
		return
	}

	suite.assert.Equal(0x922bac0f, data)

	data, err = suite.beBitBuffer.ReadUint64(32)

	if err != nil {
		return
	}

	suite.assert.Equal(0x92db0fac, data)
}

func (suite *BitBufferTestSuite) TestReadString() {
	_, err := suite.leBitBuffer.Read(5 * 8)

	if err != nil {
		return
	}

	leData, err := suite.leBitBuffer.ReadString(16)

	if err != nil {
		return
	}

	suite.assert.Equal("ak", string(leData))

	_, err = suite.beBitBuffer.Read(6 * 8)

	beData, err := suite.beBitBuffer.ReadString(16)

	if err != nil {
		return
	}

	suite.assert.Equal("kG", string(beData))
}

func TestBitBufferTestSuite(t *testing.T) {
	suite.Run(t, new(BitBufferTestSuite))
}
