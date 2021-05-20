import {Client} from '../../../../../pkg/client/js'

const CLIENT_ID = Math.random().toString()

let client = new Client({
    transport: {
        type: 'ws',
        url: 'ws://localhost:8080',
        pingInterval: 15000,
        onPong: ms => console.log(`[ping] ${ms}ms`)
    },
    auth: CLIENT_ID,
})

window.client = client

Object.defineProperty(window, 'cls', {
    get() {
        console.clear()
    }
})

let requestRegistry = {}

window.req = payload => {
    return new Promise((ok, err) => {
        let id = Math.random().toString()
        requestRegistry[id] = ok
        client.emit('request', id, payload)
    })
}

client.on('response', (id, payload) => {
    let ok = requestRegistry[id]
    if(ok) {
        ok(payload)
        delete requestRegistry[id]
    }
})

client.hook('open', () => {
    console.log('client is open...')
})

client.hook('error', err => console.log(`error: ${err}`))
client.connect()

client.on('_', (name, ...args) => {
    console.log(`%s: %o`, name, args)
})