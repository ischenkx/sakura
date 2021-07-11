const OPCODES = {
    messageCode: 1,
    pingCode: 2,
    pongCode: 3,
    authAckCode: 4,
    authReqCode: 5,
}

export class NotifySocket {
    constructor(cfg) {
        this.ws = null
        this.protocols = cfg.protocols
        this.url = cfg.url
        this.pingInterval = cfg.pingInterval || 30000
        this.pingTimeout = cfg.pingTimeout || 30000
        this.onPong = cfg.onPong
        this.onopen = null;
        this.onerror = null;
        this.onclose = null;
        this.onmessage = null;
        this.onauth = null;
        this.pendingPings = {}
        this.uniqueID = 0
    }

    __getUniqueID() {
        return ++this.uniqueID
    }

    __confirmPong(id) {
        let ok = this.pendingPings[id]
        if (ok) ok(new Date)
        delete this.pendingPings[id]
    }

    __startPinging() {
        setTimeout(() => {
            if (!this.ws) return
            let payload = this.__getUniqueID().toString()
            this.sendServiceMessage(OPCODES.pingCode, payload)
            let sendTime = Date.now()
            new Promise((ok, err) => {
                this.pendingPings[payload] = ok
                setTimeout(() => {
                    delete this.pendingPings[payload]
                    err('pong timeout')
                }, this.pingTimeout)
            }).then(date => {
                let ms = date.getTime() - sendTime
                if (this.onPong) {
                    this.onPong(ms)
                }
                this.__startPinging()
            }).catch((err) => {
                if (this.onerror) this.onerror(new Error(`ping error: ${err}`))
            })
        }, this.pingInterval)
    }

    connect() {
        this.ws = new WebSocket(this.url, this.protocols)
        this.ws.onopen = (...data) => {
            this.__startPinging()
            this.onopen(...data)
        }
        this.ws.onerror = this.onerror
        this.ws.onclose = this.onclose

        this.ws.onmessage = async ev => {

            let payload = ev

            if (typeof payload !== 'string') {
                payload = await ev.data.text()
            }

            let {opCode, message} = decodeMessage(payload)
            switch (opCode) {
                case OPCODES.authAckCode:
                    let err = message === 'ok' ? null : new Error('failed to authenticate')
                    if (this.onauth) this.onauth(err)
                    break
                case OPCODES.messageCode:
                    if (this.onmessage) {
                        this.onmessage(message)
                    }
                    break
                case OPCODES.pongCode:
                    this.__confirmPong(message)
                    break
            }


        }
    }

    authenticate(data) {
        this.sendServiceMessage(OPCODES.authReqCode, data)
    }

    close() {
        if (!this.ws) return
        this.ws.close()
        this.ws = null
    }

    sendServiceMessage(opCode, data) {
        if (!this.ws) throw new Error('no open websocket available')
        this.ws.send(encodeMessage(opCode, data))
    }

    send(data) {
        if (!this.ws) return
        this.sendServiceMessage(OPCODES.messageCode, data)
    }
}

function decodeMessage(data) {
    if (data.length === 0) {
        throw new Error('wrong length')
    }
    return {opCode: data[0].charCodeAt(0), message: data.slice(1)}
}

function encodeMessage(code, data) {
    return String.fromCharCode(code) + data
}