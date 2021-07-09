export function writeUin16Bytes(n) {
    let byte = 0b1111_1111
    let b = [0, 0]
    b[0] = n & byte
    b[1] = (n >> 8) & byte
    return b
}