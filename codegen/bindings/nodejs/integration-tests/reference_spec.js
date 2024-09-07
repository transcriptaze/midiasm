import fs from 'node:fs'

import { describe, it } from 'mocha'
import { expect } from 'chai'
import { disassemble, ByteReader } from '../index.js'

describe('reference file', function () {
  it('#disassembles the reference MIDI file', function () {
    const buffer = fs.readFileSync('./integration-tests/midi/reference.mid')
    const bytes = Uint8Array.from(buffer)
    const expected = JSON.parse(fs.readFileSync('./integration-tests/json/reference.json', 'utf8'))

    const reader = new ByteReader(Array.from(bytes))
    const midi = disassemble(reader)

    expect(midi).to.deep.equal(expected)
  })
})
