import {describe,it} from 'mocha'
import { expect } from 'chai'
import {disassemble} from '../index.js'

describe('disassemble', function () {
  describe('#disassembles a MIDI file', function () {
    it('should succeed', function () {
      const bytes = []
      const expected = {}
      const object = disassemble(bytes)

      expect(object).to.equal(expected)
    })
  })
})
