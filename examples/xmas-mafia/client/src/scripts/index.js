import {Client, Emitter} from "../../lib"

async function main() {
    let client = new Client({
        transport: {
            url: 'http://localhost:4545/pubsub?room=who-the-fuck-am-i',
            type: 'sockjs',
        },
        reconnection: {
            fn: t => t * 2,
            initialTimeout: 1000,
        }
    })

    client.on('open', () => console.log('client is open'))

    client.on('close', () => console.log('client is closed'))

    client.on('error', err => console.log('error:', err))

    let emitter = new Emitter(client)

    window.emitter = emitter

    emitter.on('message', (data)=>{
        console.log('new message', data)
    })
    window.client = client
    client.connect()
}
window.addEventListener('load', main)