export function parseMessages(buffer) {
    let messages = []
    while(buffer.length > 4) {
        let length = buffer[0]|buffer[1]<<8|buffer[2]<<16|buffer[3]<<24
        let messageBytes = buffer.slice(4, 4+length)
        let message = String.fromCharCode(...messageBytes)
        messages.push(message)
        buffer = buffer.slice(4+length)
    }
    console.log('parsed messages: ', messages.length)
    return messages
}
