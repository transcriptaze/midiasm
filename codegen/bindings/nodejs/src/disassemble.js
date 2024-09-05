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
  }
  //
  // d.FieldUTF8("tag", 4)
  // length := d.FieldS32("length")
  //
  // d.FramedFn(length*8, func(d *decode.D) {
  //     format := d.FieldU16("format")
  //     if format != 0 && format != 1 && format != 2 {
  //         d.Errorf("invalid MThd format %v (expected 0,1 or 2)", format)
  //     }
  //
  //     tracks := d.FieldU16("tracks")
  //     if format == 0 && tracks > 1 {
  //         d.Errorf("MIDI format 0 expects 1 track (got %v)", tracks)
  //     }
  //
  //     division := d.FieldU16("divisions")
  //     if division&0x8000 == 0x8000 {
  //         SMPTE := (division & 0xff00) >> 8
  //         if SMPTE != 0xe8 && SMPTE != 0xe7 && SMPTE != 0xe6 && SMPTE != 0xe5 {
  //             d.Errorf("invalid MThd division SMPTE timecode type %02X (expected E8,E7, E6 or E5)", SMPTE)
  //         }
  //     }
  // })

  return header
}

function decodeMTrk(reader) {

}
