import {Client} from '../../../../pkg/client/js/index'

const CLIENT_ID = Math.random().toString()
const USER_ID = Math.random().toString()
const AUTH_URl = `http://localhost:4646/token?user=${USER_ID}&client=${CLIENT_ID}`

fetch(AUTH_URl)
    .catch(err => console.log(`failed to fetch auth token: ${err}`))
    .then(data => data.text())
    .then(token => {
        let client = new Client({
            transport: {
                type: 'ws',
                url: 'ws://localhost:4545',
                pingInterval: 15000,
            },
            auth: token,
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
    })

