import { describe, it } from 'mocha'
import { expect } from 'chai'
import { disassemble } from '../index.js'

describe('disassemble', function () {
  describe('#disassembles a MIDI file', function () {
    it('should fail with "invalid MIDI file"', function () {
      expect(() => disassemble()).to.throw(Error, 'invalid MIDI file')
    })
  })
})
