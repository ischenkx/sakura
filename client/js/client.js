
const OPEN_EVENT = 'open'
const MESSAGE_EVENT = 'message'
const CLOSE_EVENT = 'close'
const ERROR_EVENT = 'error'

const FORCIBLY_CLOSED_CLIENT = 'fc_client'
const OPEN_CLIENT = 'o_client'
const CLOSED_CLIENT = 'c_client'

function parseMessages(buffer) {
    let messages = []
    while(buffer.length > 4) {
        let length = buffer[0]|buffer[1]<<8|buffer[2]<<16|buffer[3]<<24
        let messageBytes = buffer.slice(4, 4+length)
        let message = String.fromCharCode(...messageBytes)
        messages.push(message)
        buffer = buffer.slice(4+length)
    }
    return messages
}

class Reconnector {
    constructor(config) {
        config = config || {}
        this.timeout = config.initialTimeout || 1000
        this.maxRetries = config.retries || 10
        this.retriesOut = 0
        this.fn = config.fn || (t => t * 2)
    }

    reset() {
        this.timeout = 100
        this.retriesOut = 0
    }

    next() {
        return new Promise((ok, err) => {
            if(this.retriesOut >= this.maxRetries) {
                err('retries are out')
                return
            }
            this.retriesOut += 1
            setTimeout(ok, this.timeout)
            this.timeout = this.fn(this.timeout)
        })
    }
}

/*
        Transport should have methods:
            - onmessage
            - onerror
            - onopen
            - onclose
            - connect
            - close
            - send
*/

class Client {
    constructor(config) {
        config = config||{}
        this.authData = config.auth||''
        this.transport = config.transport
        this.events = {}
        this.queue = []
        this.reconnector = new Reconnector(config.reconnect)
        this.__state = CLOSED_CLIENT
    }

    __authenticate() {
        this.transport.send(this.authData)
    }

    __enqueueData(data) {
        this.queue.push(data)
    }

    __flushQueue() {
        let newQueue = []
        for(let data of this.queue) {
            try {
                this.transport.send(data)
            } catch {
                newQueue.push(data)
            }
        }
        this.queue = newQueue
    }

    __emitEvent(name, ...data) {
        let handlers = this.events[name]
        if(handlers) {
            for(let h of handlers) {
                h(...data)
            }
        }
    }

    __setupConnection() {
        this.transport.onmessage = async (data)=>{
            let messages = parseMessages(data)
            for(let mes of messages) {
                this.__emitEvent(MESSAGE_EVENT, mes)
            }
        }

        this.transport.onopen = async ()=>{
            this.__state = OPEN_CLIENT
            this.__emitEvent(OPEN_EVENT)
            this.__authenticate()
            this.__flushQueue()
            this.reconnector.reset()
        }

        this.transport.onerror = async err =>{
            this.__emitEvent(ERROR_EVENT, err)
        }

        this.transport.onclose = async ()=>{
            let forciblyClosed = this.__state === FORCIBLY_CLOSED_CLIENT
            this.__state = CLOSED_CLIENT
            this.__emitEvent(CLOSE_EVENT)
            if(!forciblyClosed) {
                try {
                    await this.reconnector.next()
                    this.__setupConnection()
                } catch(err) {
                    console.log(err)
                }
            }   
        }
        this.transport.connect()
    }

    on(name, h) {
        let handlers = this.events[name]
        if(!handlers) {
            handlers = []
            this.events[name] = handlers
        }
        handlers.push(h)
    }

    off(name, h) {
        let handlers = this.events[name]
        if(!handlers) {
            return
        }
        this.events[name] = handlers.filter(handler => handler !== h)
    }

    send(data) {
        if(this.__state !== OPEN_CLIENT) {
            this.__enqueueData(data)
            return
        }
        this.transport.send(data)
    }

    connect() {
        this.__setupConnection()
    }

    disconnect() {
        this.transport.close()
    }
}


class NotifySocket {
    constructor(url, protocols) {
        this.ws = null
        this.protocols = protocols
        this.url = url
    }
    connect() {
        this.ws = new WebSocket(this.url, this.protocols)
        this.ws.onopen = this.onopen
        this.ws.onerror = this.onerror
        this.ws.onclose = this.onclose
        this.ws.onmessage = async event => {
            if(this.onmessage) {
                this.onmessage(new Uint8Array(await event.data.arrayBuffer()))
            }
        }
    }
    close() {
        if(!this.ws) return
        this.ws.close()
    }
    send(data) {
        if(!this.ws) return
        this.ws.send(data)
    }
}