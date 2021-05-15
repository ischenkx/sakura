import {Client} from '../../../../../pkg/client/js'

const CLIENT_ID = Math.random().toString()

let client = new Client({
    transport: {
        type: 'ws',
        url: 'ws://localhost:8080',
        pingInterval: 15000,
    },
    auth: CLIENT_ID,
})

window.send = mes => {
    client.emit('chat.message', {
        payload: mes,
        code: 100,
    })
}

window.client = client
Object.defineProperty(window, 'cls', {
    get() {
        console.clear()
    }
})
client.on('chat.error', console.log)
client.on('_', console.log)
client.on('chat.message', console.log)

client.hook('open', ()=> {
    console.log('client is open')
    client.emit('chat.join', {name: 'chat', code: 100})
})

client.connect()

