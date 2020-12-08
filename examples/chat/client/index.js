import {Client} from "./lib/index"

let messageContainer = document.querySelector('#messages')
let messageInput = document.querySelector('#message-input')
let client = new Client({
    transport: {
        url: 'http://localhost:3000/pubsub',
        type: 'sockjs',
        options: {

        },
    },

    auth: 'some-auth-data',
    reconnection: {
        fn: t => t * 2,
        initialTimeout: 1000,
    }
})

console.log(client)

client.on('open', () => {
    console.log('client is open')
})
client.on('message', data => {
    let mes = document.createElement('div')
    mes.innerHTML = data
    mes.className = 'message'
    messageContainer.appendChild(mes)
})
client.on('close', () => {
    console.log('client is closed')
})
client.on('error', err => {
    console.log('error:', err)
})
client.connect()

messageInput.addEventListener('keydown', e => {
    if(e.key === 'Enter') {
        client.send(messageInput.value)
        messageInput.value = ''
    }
})

