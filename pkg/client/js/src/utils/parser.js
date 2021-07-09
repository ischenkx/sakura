export function parseMessages(buffer) {
    let messages = []
    while (buffer.length > 4) {
        let length = buffer.charCodeAt(0) | buffer.charCodeAt(1) << 8 | buffer.charCodeAt(2) << 16 | buffer.charCodeAt(3) << 24
        let message = buffer.slice(4, 4 + length)
        messages.push(message)
        buffer = buffer.slice(4 + length)
    }
    return messages
}
