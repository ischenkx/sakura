export function readUint16Bytes(b) {
    return b[0] | b[1] << 8
}