
let messageContainer = document.querySelector('#messages')
let messageInput = document.querySelector('#message-input')
let client = new Client({
    transport: new NotifySocket('ws://localhost:3000/pubsub'),
    auth: 'some-auth-data',
    reconnect: {
        fn: t => t * 2,
        initialTimeout: 1000,
    }
})
client.on('open', () => console.log('connection is established'))
client.on('message', data => {
    let mes = document.createElement('div')
    mes.innerHTML = data
    mes.className = 'message'
    messageContainer.appendChild(mes)
})
client.on('close', () => console.log('connection is closed'))
client.on('error', err => console.log('error:', err))
client.connect()

messageInput.addEventListener('keydown', e => {
    if(e.key === 'Enter') {
        client.send(messageInput.value)
        messageInput.value = ''
    }
})