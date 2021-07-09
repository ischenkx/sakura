const clientEvents = `type ClientEvents interface {
	OnEmit(func(EventOptions))
	OnEvent(func(EventOptions))
	OnUnsubscribe(func(string, int64))
	OnSubscribe(func(string, int64))
	OnReconnect(func(int64))
	OnDisconnect(func(int64))
	Close()
}`

const userEvents = `type UserEvents interface {
	OnEmit(func(EventOptions))
	OnEvent(func(EventOptions))
	OnUnsubscribe(func(string, int64))
	OnSubscribe(func(string, int64))
	OnClientConnect(func(Client, int64))
	OnClientReconnect(func(Client, int64))
	OnClientDisconnect(func(string, int64))
	Close()
}`

const topicEvents = `type TopicEvents interface {
	OnClientSubscribe(func(Client))
	OnUserSubscribe(func(User))
	OnClientUnsubscribe(func(Client))
	OnUserUnsubscribe(func(User))
	OnEmit(func(EventOptions))
	Close()
}`

const regex = /On(?<name>\w+)\s*\(func\((?<args>[\w\s,*]+)\)\s*\)/gm


function generateClientEvents() {
    let structName = `localClientsEventsRegistry`
    let shortName = `reg`

    let iter = clientEvents.matchAll(regex)
    let res = iter.next()

    let fields = []
    let methods = []

    fields.push('_presence map[string]string')
    fields.push('_handlers map[string][]string')
    let names = []

    while (!res.done) {
        let {name, args} = res.value.groups
        names.push(name)
        fields.push(`_${name} map[string][]func(${args})`)
        let addHandlerDecl = `func (${shortName} *${structName}) handle${name}(uid string, f func(${args})) {
            ${shortName}.mu.Lock()
            defer ${shortName}.mu.Unlock()
            _, ok := ${shortName}._presence[uid]
            if !ok {
                return
            }
            ${shortName}._${name}[uid] = append(${shortName}._${name}[uid], f)
        }`

        let callHandlerDecl = `func (${shortName} *${structName}) call${name}(__id string, ${
            args.split(',')
                .map(x => x.trim())
                .map((x, i) => `arg${i} ${x}`).join(', ')
        }) {
            ${shortName}.mu.RLock()
            uids, ok := ${shortName}._handlers[__id]
            if !ok {
                ${shortName}.mu.RUnlock()
                return
            }
            var handlers []func(${args})
            for _, uid := range uids {
                if hs, ok := ${shortName}._${name}[uid]; ok {
                    handlers = append(handlers, hs...)
                }
            }
            ${shortName}.mu.RUnlock()
            
            for _, h := range handlers {
                h(${args.split(',').map(x => x.trim()).map((x, i) => `arg${i}`).join(', ')})
            }
        }`

        methods.push(addHandlerDecl, callHandlerDecl)
        res = iter.next()
    }

    let deleteEventsDecl = `func (${shortName} *${structName}) deleteEvents(uid string) {
            ${shortName}.mu.Lock()
            defer ${shortName}.mu.Unlock()
            __id, ok := ${shortName}._presence[uid]
            if !ok {
                return
            }
            handlers := ${shortName}._handlers[__id]
            
            for i := 0; i < len(handlers); i++ {
                if handlers[i] == uid {
                    handlers[i] = handlers[len(handlers)-1]
                    handlers = handlers[:len(handlers)-1]
                    break
                }
            }
            ${shortName}._handlers[__id] = handlers
            ${
        names.map(name => `delete(${shortName}._${name}, uid)`).join('\n')
    }
        }`

    let createEventsDecl = `func (${shortName} *${structName}) createEvents(__id string) localClientEvents {
            ${shortName}.mu.Lock()
            defer ${shortName}.mu.Unlock()
            
            uid := uuid.New().String()
            
            ${shortName}._presence[uid] = __id
            ${shortName}._handlers[__id] = append(${shortName}._handlers[__id], uid)
            
            return localClientEvents{reg: ${shortName}, clientID: __id, id: uid}
        }`
    methods.push(createEventsDecl, deleteEventsDecl)

    let structDecl = `type ${structName} struct {
    ${fields.join('\n')}
}`

    let methodsDecl = methods.join('\n')

    return `${structDecl}\n\n${methodsDecl}`

}


function generateTopicEvents() {
    let structName = `localTopicsEventsRegistry`
    let shortName = `reg`

    let iter = topicEvents.matchAll(regex)
    let res = iter.next()

    let fields = []
    let methods = []

    fields.push('mu sync.RWMutex')
    fields.push('_presence map[string]string')
    fields.push('_handlers map[string][]string')
    let names = []

    while (!res.done) {
        let {name, args} = res.value.groups
        names.push(name)
        fields.push(`_${name} map[string][]func(${args})`)
        let addHandlerDecl = `func (${shortName} *${structName}) handle${name}(uid string, f func(${args})) {
            ${shortName}.mu.Lock()
            defer ${shortName}.mu.Unlock()
            _, ok := ${shortName}._presence[uid]
            if !ok {
                return
            }
            ${shortName}._${name}[uid] = append(${shortName}._${name}[uid], f)
        }`

        let callHandlerDecl = `func (${shortName} *${structName}) call${name}(__id string, ${
            args.split(',')
                .map(x => x.trim())
                .map((x, i) => `arg${i} ${x}`).join(', ')
        }) {
            ${shortName}.mu.RLock()
            uids, ok := ${shortName}._handlers[__id]
            if !ok {
                ${shortName}.mu.RUnlock()
                return
            }
            var handlers []func(${args})
            for _, uid := range uids {
                if hs, ok := ${shortName}._${name}[uid]; ok {
                    handlers = append(handlers, hs...)
                }
            }
            ${shortName}.mu.RUnlock()
            
            for _, h := range handlers {
                h(${args.split(',').map(x => x.trim()).map((x, i) => `arg${i}`).join(', ')})
            }
        }`

        methods.push(addHandlerDecl, callHandlerDecl)
        res = iter.next()
    }

    let deleteEventsDecl = `func (${shortName} *${structName}) deleteEvents(uid string) {
            ${shortName}.mu.Lock()
            defer ${shortName}.mu.Unlock()
            __id, ok := ${shortName}._presence[uid]
            if !ok {
                return
            }
            handlers := ${shortName}._handlers[__id]
            
            for i := 0; i < len(handlers); i++ {
                if handlers[i] == uid {
                    handlers[i] = handlers[len(handlers)-1]
                    handlers = handlers[:len(handlers)-1]
                    break
                }
            }
            ${shortName}._handlers[__id] = handlers
            ${
        names.map(name => `delete(${shortName}._${name}, uid)`).join('\n')
    }
        }`

    let createEventsDecl = `func (${shortName} *${structName}) createEvents(__id string) localTopicEvents {
            ${shortName}.mu.Lock()
            defer ${shortName}.mu.Unlock()
            
            uid := uuid.New().String()
            
            ${shortName}._presence[uid] = __id
            ${shortName}._handlers[__id] = append(${shortName}._handlers[__id], uid)
            
            return localTopicEvents{reg: ${shortName}, topicID: __id, id: uid}
        }`
    methods.push(createEventsDecl, deleteEventsDecl)

    let structDecl = `type ${structName} struct {
    ${fields.join('\n')}
}`

    let methodsDecl = methods.join('\n')

    return `${structDecl}\n\n${methodsDecl}`

}

function generateUserEvents() {
    let structName = `localUsersEventsRegistry`
    let shortName = `reg`

    let iter = userEvents.matchAll(regex)
    let res = iter.next()

    let fields = []
    let methods = []

    fields.push('mu sync.RWMutex')
    fields.push('_presence map[string]string')
    fields.push('_handlers map[string][]string')
    let names = []

    while (!res.done) {
        let {name, args} = res.value.groups
        names.push(name)
        fields.push(`_${name} map[string][]func(${args})`)
        let addHandlerDecl = `func (${shortName} *${structName}) handle${name}(uid string, f func(${args})) {
            ${shortName}.mu.Lock()
            defer ${shortName}.mu.Unlock()
            _, ok := ${shortName}._presence[uid]
            if !ok {
                return
            }
            ${shortName}._${name}[uid] = append(${shortName}._${name}[uid], f)
        }`

        let callHandlerDecl = `func (${shortName} *${structName}) call${name}(__id string, ${
            args.split(',')
                .map(x => x.trim())
                .map((x, i) => `arg${i} ${x}`).join(', ')
        }) {
            ${shortName}.mu.RLock()
            uids, ok := ${shortName}._handlers[__id]
            if !ok {
                ${shortName}.mu.RUnlock()
                return
            }
            var handlers []func(${args})
            for _, uid := range uids {
                if hs, ok := ${shortName}._${name}[uid]; ok {
                    handlers = append(handlers, hs...)
                }
            }
            ${shortName}.mu.RUnlock()
            
            for _, h := range handlers {
                h(${args.split(',').map(x => x.trim()).map((x, i) => `arg${i}`).join(', ')})
            }
        }`

        methods.push(addHandlerDecl, callHandlerDecl)
        res = iter.next()
    }

    let deleteEventsDecl = `func (${shortName} *${structName}) deleteEvents(uid string) {
            ${shortName}.mu.Lock()
            defer ${shortName}.mu.Unlock()
            __id, ok := ${shortName}._presence[uid]
            if !ok {
                return
            }
            handlers := ${shortName}._handlers[__id]
            
            for i := 0; i < len(handlers); i++ {
                if handlers[i] == uid {
                    handlers[i] = handlers[len(handlers)-1]
                    handlers = handlers[:len(handlers)-1]
                    break
                }
            }
            ${shortName}._handlers[__id] = handlers
            ${
        names.map(name => `delete(${shortName}._${name}, uid)`).join('\n')
    }
        }`

    let createEventsDecl = `func (${shortName} *${structName}) createEvents(__id string) localUserEvents {
            ${shortName}.mu.Lock()
            defer ${shortName}.mu.Unlock()
            
            uid := uuid.New().String()
            
            ${shortName}._presence[uid] = __id
            ${shortName}._handlers[__id] = append(${shortName}._handlers[__id], uid)
            
            return localUserEvents{reg: ${shortName}, userID: __id, id: uid}
        }`
    methods.push(createEventsDecl, deleteEventsDecl)

    let structDecl = `type ${structName} struct {
    ${fields.join('\n')}
}`

    let methodsDecl = methods.join('\n')

    return `${structDecl}\n\n${methodsDecl}`

}

console.log(generateUserEvents())