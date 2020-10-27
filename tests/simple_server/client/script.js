// WebSocket = undefined
let client = null
const urlParams = new URLSearchParams(window.location.search);
const port = urlParams.get('port')||3000;
let id = ''

const sendBtn = document.querySelector('#send-button')
const messagesBlock = document.querySelector('#messages-block')
const messageInput = document.querySelector('#message-input')

// fetch(`http://localhost:${port}/register`, {
//   body: JSON.stringify({
//       UserID: Math.random().toString(),
//       ID: Math.random().toString(),
//       AppID: 'app-hello',
//   }),
//   method: 'POST'
// }).then(async d => {
//     id = await d.text()
//     connect()
// })

function disconnect() {
    client.close()
}

function makeid(length) {
    var result           = '';
    var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for ( var i = 0; i < length; i++ ) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}

async function connect() {
    let sock = new SockJS(`http://localhost:6565/pubsub`)
    sock.onopen = ()=>{
        sock.send("")
        client = sock
        console.log('socket open')
    }
    sock.onclose = console.log
    sock.onerror = console.log
    sock.onmessage = (mes)=>{
        let mesWrapper = document.createElement('div')
        mesWrapper.className = 'message-wrapper'
        mesWrapper.classList.add('yours')
        let mesEl = document.createElement('div')
        mesEl.className = 'message'
        mesEl.innerHTML = mes.data
        mesWrapper.appendChild(mesEl)
        messagesBlock.appendChild(mesWrapper)
        setTimeout(()=>{
            messagesBlock.scroll({
                top: messagesBlock.scrollHeight,
                behavior: 'smooth'
            })
        }, 80)
    }
}

sendBtn.addEventListener('click', ()=>{
    let value = messageInput.value
    if(!value) return
    messageInput.value = ''
    if(client) client.send(value)
})

messageInput.addEventListener('keyup', (e)=>{
    if(e.key === 'Enter') sendBtn.click()
})
connect()