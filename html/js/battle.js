const field = [];
const clazz = {0: "free", 1: "boat", 2: "lose", 3: "boom"}

function point(f, i, j) {
    let e = document.getElementById(f + '.' + i + '.' + j)
    if (f === 1) {
        e.className = "void"
        e.onclick = function () {
            onClick(f, i, j)
        }
    }
    e.innerHTML = '&#x25CF;'
    return e
}

function initField() {
    for (let f = 0; f < 2; f++) {
        field[f] = []
        for (let i = 0; i < 10; i++) {
            field[f][i] = []
            for (let j = 0; j < 10; j++) {
                field[f][i][j] = point(f, i, j)
            }
        }
    }
}

let seq = 0
let socket = new WebSocket(location.origin.replace(/^http/, 'ws'))
let handle = new Map

function query(done, id, method, path, body) {
    socket.send([id.toString(), method, path, JSON.stringify(body)].join('\r\n'))
    handle.set(id.toString(), done)
}

function reply(p) {
    if (!handle.has(p[0]))
        return
    if (p[1] === '200')
        try {
            handle.get(p[0])(JSON.parse(p[3]))
        } catch (e) {
        }
    if (p[2] === '#')
        handle.delete(p[0])
}

function onClick(f, x, y) {
    query(onReply, seq++, "GET", "/click", {X: x, Y: y})
}

function onReply(e = {
    F: 0,
    X: 0,
    Y: 0,
    C: 0,
}) {
    field[e.F][e.X][e.Y].className = clazz[e.C]
    field[e.F][e.X][e.Y].onclick = undefined
}

socket.onopen = function () {
    query(onReply, seq++, "GET", "/begin")
}

socket.onmessage = function (e) {
    reply(e.data.toString().split(/\r?\n/))
}

socket.onclose = function () {
    handle.clear()
}

socket.onerror = function () {
    handle.clear()
}

window.onload = function () {
    initField()
}