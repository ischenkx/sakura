const eventsInterface = `OnEmit(f func(*App, EmitOptions))
	OnEvent(f func(*App, Client, EventOptions))
	OnConnect(f func(*App, ConnectOptions, Client))
	OnDisconnect(f func(*App, Client))
	OnReconnect(f func(*App, ConnectOptions, Client))
	OnInactivate(f func(*App, Client))
	OnError(f func(*App, error))
	OnChange(f func(*App, ChangeLog))
	OnClientSubscribe(f func(*App, Client, string, int64))
	OnUserSubscribe(f func(*App, User, string, int64))
	OnClientUnsubscribe(f func(*App, Client, string, int64))
	OnUserUnsubscribe(f func(*App, User, string, int64))`

let structName = 'localAppEvents'
let structShortName = 'evs'
let usePointerReceiver = true

let regex = /On(?<name>\w+)\s*\(f func\((?<args>[\w\s,*]+)\)\s*\)/gm
let iter = eventsInterface.matchAll(regex)
let res = iter.next()

let methods = []
let fields = []

fields.push(`priority Priority`)
fields.push(`registry *eventsRegistry`)
fields.push(`id string`)
fields.push(`mu sync.RWMutex`)

while (!res.done) {
    let {name, args} = res.value.groups

    let handlersStorageName = `_${name}`

    fields.push(`${handlersStorageName} []func(${args})`)

    let callerName = `call${name}`

    let rawArgs = args.split(',').map(x => x.trim())
        .filter(x => x !== '*App')

    let argsDecl = rawArgs
        .map((x, i) => `arg${i} ${x}`).join(', ')

    let callDecl = `func (${structShortName} ${usePointerReceiver ? '*' : ''}${structName}) ${callerName}(${argsDecl}) {
        ${structShortName}.mu.RLock()
        handlers := make([]func(${args}), len(${structShortName}.${handlersStorageName}))
        copy(handlers, ${structShortName}.${handlersStorageName})
        ${structShortName}.mu.RUnlock()
        for _, h := range handlers {
            h(${structShortName}.registry.app, ${rawArgs.map((_, i) => `arg${i}`)})
        }
    }`

    let handleName = `On${name}`

    let handleDecl = `func (${structShortName} ${usePointerReceiver ? '*' : ''}${structName}) ${handleName}(f func(${args})) {
        ${structShortName}.mu.Lock()
        ${structShortName}.${handlersStorageName} = append(${structShortName}.${handlersStorageName}, f)
        ${structShortName}.mu.Unlock()
    }`

    methods.push(callDecl, handleDecl)

    res = iter.next()
}

let structDecl = `type ${structName} struct {
    ${fields.join('\n')}
}`

let methodsDecl = methods.join('\n')

require('fs').writeFileSync('localEventsGenOut.txt', `${structDecl}\n${methodsDecl}`)

