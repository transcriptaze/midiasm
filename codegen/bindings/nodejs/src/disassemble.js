/**
 * MIDI file disassembler.
 * @module midiasm
 */

/**
   * Parses a MIDI file byte array into an object.
   *
   * @param {object} reader  MIDI file reader.
   *
   * @example
   * midiasm.disassemble(bytes)
   *  .then(object => { console.log(object) })
   *  .catch(err => { console.log(`${err.message}`)
   */
export function disassemble(reader) {
  if (reader == null) {
    throw new Error(`invalid MIDI file reader`)
  }

  // ... decode header

  const header = decodeMThd(reader)

  // ... decode tracks
  const tracks = []
  // d.FieldArray("tracks", func(d *decode.D) {
  //     for d.BitsLeft() > 0 {
  //         d.FieldStruct("track", decodeMTrk)
  //     }
  // })

  return {
    header: header,
    tracks: tracks,
  }
}

function decodeMThd(reader) {
  const header = {}
  const bytes = reader.read(4)
  const tag = Buffer.from(bytes).toString('utf8')

  if (tag !== 'MThd') {
    throw new Error(`missing MThd tag`)
  }
  else {
    header.tag = tag
    header.length = reader.U32()
    header.format = reader.U16()
    header.tracks = reader.U16()

    const division = reader.peekU16()

    if ((division & 0x8000) == 0x8000) {
      header.timecode = {}

      switch ((division >> 8) & 0x00ff) {
        case 0xe8:
          header.timecode.fps = 24
          break

        case 0xe7:
          header.timecode.fps = 25
          break

        case 0xe6:
          header.timecode.fps = 29
          break

        case 0xe5:
          header.timecode.fps = 30
          break
      }

      header.timecode.resolution = division & 0x00ff
    }
    else {
      header.ppqn = division
    }
  }

  return header
}

function decodeMTrk(reader) {

}
