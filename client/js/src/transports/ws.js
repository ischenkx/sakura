export class NotifySocket {
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