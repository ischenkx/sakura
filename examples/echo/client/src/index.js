import {Client} from 'swirljs'

let client = new Client({
    transport: {
        type: 'ws',
        url: 'ws://localhost:6060/swirl'
    },
    auth: Math.random().toString()
})

client.on('echo', data => console.log('echo:', data))

client.connect()

let input = document.querySelector('#echo-input')
input.addEventListener('keydown', ev => {
    if(ev.key === 'Enter') {
        if(input.value) {
            client.emit('echo', input.value)
        }
    }
})
