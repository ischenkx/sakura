import SockJS from "sockjs-client"

export class NotifySockJS {
    constructor(url, protocols) {
        this.sock = null
        this.protocols = protocols
        this.url = url
    }
    connect() {
        this.sock = new SockJS(this.url, this.protocols)
        this.sock.onopen = (...data) => {
            if(this.onopen)
                this.onopen(...data)
        }
        this.sock.onerror = (...data) => {
            if(this.onerror)
                this.onerror(...data)
        }
        this.sock.onclose = (...data) => {
            if(this.onclose)
                this.onclose(...data)
        }
        this.sock.onmessage = event => {
            if(this.onmessage) {
                let data = Uint8Array.from(event.data, c => c.charCodeAt(0))
                this.onmessage(data)
            }
        }
    }
    close() {
        if(!this.sock) return
        this.sock.close()
    }
    send(data) {
        if(!this.sock) return
        this.sock.send(data)
    }
}
