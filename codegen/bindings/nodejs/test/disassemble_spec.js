import { describe, it } from 'mocha'
import { expect } from 'chai'
import { disassemble, ByteReader } from '../index.js'

const MIDI = [
  0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06,
  0x00, 0x00, 0x00, 0x01, 0x01, 0xe0, 0x4d, 0x54,
  0x72, 0x6b, 0x00, 0x00, 0x00, 0x10, 0x00, 0xff,
  0x03, 0x08, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74,
  0x20, 0x30, 0x00, 0xff, 0x2f, 0x00,
]

describe('disassemble', function () {
  describe('#disassembles a MIDI file', function () {
    it('should fail with "invalid MIDI file reader"', function () {
      expect(() => disassemble()).to.throw(Error, 'invalid MIDI file reader')
    })

    it('should read a MIDI file header"', function () {
      const reader = new ByteReader(MIDI)
      const expected = {
        header: {
          tag: 'MThd',
        },
        tracks: [],
      }

      const midi = disassemble(reader)

      expect(midi).to.deep.equal(expected)
    })
  })
})
