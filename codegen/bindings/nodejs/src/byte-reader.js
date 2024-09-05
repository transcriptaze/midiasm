import { Reader } from './midiasm.js'

export class ByteReader extends Reader {
  constructor(bytes) {
    super()

    this.buffer = new Uint8Array(bytes)
    this.position = 0
  }

  peek(N) {
    const start = this.position
    const end = start + N

    return Array.from(this.buffer.slice(start, end))
  }

  read(N) {
    const start = this.position
    const end = start + N
    const bytes = Array.from(this.buffer.slice(start, end))

    this.position += N

    return bytes
  }
}
