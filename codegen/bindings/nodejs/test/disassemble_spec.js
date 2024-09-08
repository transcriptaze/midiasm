import { describe, it } from 'mocha'
import { expect } from 'chai'
import { disassemble, ByteReader } from '../index.js'

const MThd_PPQN = [
  0x4d, 0x54, 0x68, 0x64,
  0x00, 0x00, 0x00, 0x06,
  0x00, 0x01, 0x00, 0x02, 0x01, 0xe0,
]

const MThd_SMPTE = [
  0x4d, 0x54, 0x68, 0x64,
  0x00, 0x00, 0x00, 0x06,
  0x00, 0x01, 0x00, 0x02, 0xe6, 0x04,
]

const MIDI = [
  0x4d, 0x54, 0x68, 0x64,
  0x00, 0x00, 0x00, 0x06,
  0x00, 0x01, 0x00, 0x02, 0x01, 0xe0,
  0x4d, 0x54, 0x72, 0x6b,
  0x00, 0x00, 0x00, 0x00,
  0x4d, 0x54, 0x72, 0x6b,
  0x00, 0x00, 0x00, 0x00,
]

describe('disassemble', function () {
  it('should fail with "invalid MIDI file reader"', function () {
    expect(() => disassemble()).to.throw(Error, 'invalid MIDI file reader')
  })

  describe('#disassembles a MIDI MThd chunk', function () {
    it('should read a MIDI MThd chunk with a PPQN division field"', function () {
      const reader = new ByteReader(MThd_PPQN)

      const expected = {
        header: {
          tag: 'MThd',
          length: 6,
          format: 1,
          tracks: 2,
          ppqn: 480,
        },
        tracks: [],
      }

      const midi = disassemble(reader)

      expect(midi).to.deep.equal(expected)
    })

    it('should read a MIDI MThd chunk with an SMPTE timecode division field"', function () {
      const reader = new ByteReader(MThd_SMPTE)

      const expected = {
        header: {
          tag: 'MThd',
          length: 6,
          format: 1,
          tracks: 2,
          timecode: {
            fps: 29,
            resolution: 4,
          },
        },
        tracks: [],
      }

      const midi = disassemble(reader)

      expect(midi).to.deep.equal(expected)
    })
  })

  describe('#disassembles MIDI MTrk chunk', function () {
    it('should extract the MIDI MTrk chunks"', function () {
      const reader = new ByteReader(MIDI)

      const expected = {
        header: {
          tag: 'MThd',
          length: 6,
          format: 1,
          tracks: 2,
          ppqn: 480,
        },
        tracks: [
          {
            tag: 'MTrk',
            length: 0,
            events: [],
          },
          {
            tag: 'MTrk',
            length: 0,
            events: [],
          },
        ],
      }

      const midi = disassemble(reader)

      expect(midi).to.deep.equal(expected)
    })
  })
})
