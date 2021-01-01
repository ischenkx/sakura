export const JSONCodec = {
    marshal(data) {
        return JSON.stringify(data)
    },
    unmarshal(data) {
        return JSON.parse(data)
    }
}

export class Emitter {
    constructor(client, codec=JSONCodec) {
        this.__client = client
        this.__handlers = {}
        this.__codec = codec
        this.__client.on('message', data => {
            try {
                let event = this.__parseEvent(data)
                let handlers = this.__handlers[event.name]
                if(handlers) {
                    for(let h of handlers) h(event.data)
                }
            } catch(e) {
                console.error('error while processing the incoming event:', e)
            }
        })
    }

    __parseEvent(payload) {
        if(payload.length < 2) {
            throw new Error('received message too short')
        }
        let nameLength = payload[0].charCodeAt(0)
        if(payload.length < nameLength + 1) {
            throw new Error("name length is too big")
        }
        let name = payload.slice(1, 1+nameLength)
        let data = this.__codec.unmarshal(payload.slice(1+nameLength))
        return {name, data}
    }

    off(event, handler) {
        let eventHandlers = this.__handlers[event]
        if(!eventHandlers) {
            return
        }
        let i = eventHandlers.indexOf(handler)
        while(i >= 0) {
            eventHandlers[i] = eventHandlers[eventHandlers.length-1]
            i = eventHandlers.indexOf(handler)
        }
        if(eventHandlers.length === 0) {
            delete this.__handlers[event]
        }
    }

    on(event, handler) {
        let eventHandlers = this.__handlers[event]
        if(!eventHandlers) {
            eventHandlers = []
            this.__handlers[event] = eventHandlers
        }
        eventHandlers.push(handler)
    }

    emit(event, data) {
        try {
            let marshaledData = this.__codec.marshal(data)
            let eventLength = event.length
            if(eventLength > 255) {
                console.error('event length is bigger than 255')
                return
            }
            this.__client.send(String.fromCharCode(eventLength)+event+marshaledData)
        } catch(e) {
            console.error("error while emitting:", e)
        }
    }
}