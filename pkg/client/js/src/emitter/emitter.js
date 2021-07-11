import {writeUin16Bytes} from "../utils/write_bytes";

export class Emitter {
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

    emit(name, ...data) {
        this.client.__send(this.__encodeData(name, ...data))
    }

    __encodeData(name, ...data) {
        let bts = String.fromCharCode(name.length) + name
        for (let arg of data) {
            let enc = JSON.stringify(arg)
            bts += String.fromCharCode(...writeUin16Bytes(enc.length))
            bts += enc
        }
        return bts
    }

    __handleEvent(n, ...data) {
        if (this.handlers[n]) {
            this.handlers[n].forEach(handler => handler(...data))
        } else if (this.handlers['_']) {
            this.handlers['_'].forEach(handler => handler(n, ...data))
        }
    }

    __initEvents() {
        this.client.hook('message', mes => {
            let len = mes[0].charCodeAt(0)
            if(mes.length < len+1) {
                return
            }
            let name = mes.slice(1, len + 1)
            try {
                let data = JSON.parse(mes.slice(len + 1))
                if(!Array.isArray(data)) return
                this.__handleEvent(name, ...data)
            } catch(ex) {

            }

        })
    }
}