import { Reader } from './midiasm.js'

export class ByteReader extends Reader {
  constructor(bytes) {
    super()

    this.view = new DataView(new Uint8Array(bytes).buffer)
    this.position = 0
  }

  /* Fetches N bytes from the internal byte array without updating the position
   * marker.
   *
   */
  peek(N) {
    const offset = this.position
    const buffer = this.view.buffer
    const slice = new Uint8Array(buffer, offset, N)

    return Array.from(slice)
  }

  /* Reads N bytes from the internal byte array.
   *
   */
  read(N) {
    const offset = this.position
    const buffer = this.view.buffer
    const slice = new Uint8Array(buffer, offset, N)

    this.position += N

    return Array.from(slice)
  }

  /* Reads a big-endian uint16 from the internal byte array without
   * updating the internal position.
   */
  peekU16() {
    const offset = this.position
    const view = this.view
    const u16 = view.getUint16(offset)

    return u16
  }

  /* Reads a big-endian uint16 from the internal byte array.
   *
   */
  U16() {
    const offset = this.position
    const view = this.view
    const u16 = view.getUint16(offset)

    this.position += 2

    return u16
  }

  /* Reads a big-endian uint32 from the internal byte array.
   *
   */
  U32() {
    const offset = this.position
    const view = this.view
    const u32 = view.getUint32(offset)

    this.position += 4

    return u32
  }
}
