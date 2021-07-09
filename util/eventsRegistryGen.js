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

let structName = 'eventsRegistry'
let structShortName = 'reg'
let usePointerReceiver = true

let methods = []
let fields = []

fields.push(`app *App`)
fields.push(`mu sync.RWMutex`)
fields.push('appEventHandlers []*localAppEvents')
fields.push('localEntities *localEntitiesEventsRegistry')

let regex = /On(?<name>\w+)\s*\(f func\((?<args>[\w\s,*]+)\)\s*\)/gm
let iter = eventsInterface.matchAll(regex)
let res = iter.next()

while (!res.done) {
    let {name, args} = res.value.groups
    let callerName = `call${name}`

    let rawArgs = args.split(',').map(x => x.trim())
        .filter(x => x !== '*App')

    let argsDecl = rawArgs
        .map((x, i) => `arg${i} ${x}`).join(', ')

    let callDecl = `func (${structShortName} ${usePointerReceiver ? '*' : ''}${structName}) ${callerName}(${argsDecl}) {
        ${structShortName}.mu.RLock()
        handlers := make([]*localAppEvents, len(${structShortName}.appEventHandlers))
        copy(handlers, ${structShortName}.appEventHandlers)
        ${structShortName}.mu.RUnlock()
        for _, h := range handlers {
            h.${callerName}(${rawArgs.map((_, i) => `arg${i}`)})
        }
    }`

    methods.push(callDecl)

    res = iter.next()
}

let structDecl = `type ${structName} struct {
    ${fields.join('\n')}
}`

let methodsDecl = methods.join('\n')

require('fs').writeFileSync('eventsRegistryGenOut.txt', `${structDecl}\n${methodsDecl}`)