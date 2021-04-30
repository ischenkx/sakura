import {Client} from '../../../../pkg/client/js/index'

let client = new Client({
    transport: {
        type: 'ws',
        url: 'ws://localhost:3434'
    },
    auth: Math.random().toString(),
})

class Emitter {
    constructor(client) {
        this.client = client
        this.handlers = {}
        this.__initEvents()
    }

    on(name, handler) {
        let handlers = this.handlers[name]
        if (!handlers) {
            handlers = []
            this.handlers[name] = handlers
        }
        handlers.push(handler)
    }

    emit(name, data) {
        this.client.send(this.__encodeData(name, data))
    }

    __encodeData(name, data) {
        return String.fromCharCode(name.length)+name+JSON.stringify(data)
    }

    __handleEvent(n, ...data) {
        if (this.handlers[n]) {
            this.handlers[n].forEach(handler => handler(...data))
        } else if (this.handlers['_']) {
            this.handlers['_'].forEach(handler => handler(n, ...data))
        }
    }

    __initEvents() {
        this.client.on('message', (mes, batch) => {
            let len = mes[0].charCodeAt(0)
            let name = mes.slice(1, len + 1)
            let data = JSON.parse(mes.slice(len + 1))
            this.__handleEvent(name, data, batch)
        })
    }
}

let emitter = new Emitter(client)

emitter.on('chat.error', console.log)
emitter.on('chat.message', console.log)
emitter.on('_', console.log)

window.emitter = emitter
window.client = client

client.connect()