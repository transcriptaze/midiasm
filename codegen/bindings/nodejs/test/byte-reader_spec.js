import { describe, it } from 'mocha'
import { expect } from 'chai'
import { disassemble, ByteReader } from '../index.js'

const MIDI = [
  0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06,
  0x00, 0x01, 0x00, 0x02, 0x01, 0xe0, 0x4d, 0x54,
  0x72, 0x6b, 0x00, 0x00, 0x00, 0x10, 0x00, 0xff,
  0x03, 0x08, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74,
  0x20, 0x30, 0x00, 0xff, 0x2f, 0x00,
]

describe('ByteReader', function () {
  describe('#peek', function () {
    it('should return the first 4 bytes of the byte array without updating the reader position"', function () {
      const reader = new ByteReader(MIDI)

      expect(reader.peek(4)).to.deep.equal([0x4d, 0x54, 0x68, 0x64])
      expect(reader.peek(4)).to.deep.equal([0x4d, 0x54, 0x68, 0x64])
    })
  })

  describe('#read', function () {
    it('should return the sequential 4 byte arrays"', function () {
      const reader = new ByteReader(MIDI)

      expect(reader.read(4)).to.deep.equal([0x4d, 0x54, 0x68, 0x64])
      expect(reader.read(4)).to.deep.equal([0x00, 0x00, 0x00, 0x06])
    })
  })

  describe('#peekU16', function () {
    it('should return the uint16 values at offset 8"', function () {
      const reader = new ByteReader(MIDI)

      reader.read(8)
      expect(reader.peekU16()).to.equal(1)
      expect(reader.peekU16()).to.equal(1)
    })
  })

  describe('#U16', function () {
    it('should return the uint16 values at offset 8 and offset 10"', function () {
      const reader = new ByteReader(MIDI)

      reader.read(8)
      expect(reader.U16()).to.equal(1)
      expect(reader.U16()).to.equal(2)
    })
  })

  describe('#U32', function () {
    it('should return the uint32 value at offset 4"', function () {
      const reader = new ByteReader(MIDI)

      reader.read(4)
      expect(reader.U32()).to.equal(6)
    })
  })
})
