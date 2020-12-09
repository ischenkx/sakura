
import { parseMessages } from './utils/parser'
import { Reconnector } from './reconnector/reconnector'
import { NotifySocket } from './transports/ws'
import { NotifySockJS } from './transports/sockjs'

const OPEN_EVENT = 'open'
const MESSAGE_EVENT = 'message'
const CLOSE_EVENT = 'close'
const ERROR_EVENT = 'error'

const FORCIBLY_CLOSED_CLIENT = 'fc_client'
const OPEN_CLIENT = 'o_client'
const CLOSED_CLIENT = 'c_client'
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

export class Client {
    constructor(config) {
        this.__initialize(config||{})
    }

    __initialize(config) {
        this.auth = config.auth || ''
        if(!config.transport) {
            throw new Error("notify client: no transport provided")
        }
        
        switch(config.transport.type) {
            case "ws":
                this.transport = new NotifySocket(config.transport.url, config.transport.protocols)
            break
            case "sockjs":
                this.transport = new NotifySockJS(config.transport.url, config.transport.options)
            break
        }

        this.events = {}
        this.queue = []
        this.reconnector = new Reconnector(config.reconnection)
        this.__state = CLOSED_CLIENT
    }

    __authenticate() {
        this.transport.send(this.auth)
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


