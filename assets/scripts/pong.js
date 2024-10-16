const socket = new WebSocket("ws://localhost:8080/ws/pong");

const joinBtn = document.getElementById("joinBtn")
const createBtn = document.getElementById("createBtn")
const codeInput = document.getElementById("code")

const room = document.getElementById("waiting-room")
const roomCode = document.getElementById("roomCode")
const player1 = document.getElementById("player1")
const player2 = document.getElementById("player2")

const EventTypeMessage = 1

class Game {
    constructor(code, p1, p2, ball) {
        this.code = code
        this.p1 = p1
        this.p2 = p2
        this.ball = ball
        this.hasStarted = false
    }
}

class Player {
    constructor(username) {
        this.username = username
    }
}

class Bar {
    constructor(x, y) {
        this.x = x
        this.y = y
    }
}

class Ball {
    constructor(x, y, speed) {
       this.x = x
       this.y = y
       this.speed = speed
    }
}


socket.addEventListener("open", (e) => {
    const message = {
      message: "Hello Server!",
    }
    sendEvent(EventTypeMessage, message)
});

function sendEvent(eventType, data) {
    const event = {
       type: eventType,
       data: data,
    }
    socket.send(JSON.stringify(event))
}

socket.addEventListener("message", (event) => {
});

joinBtn.addEventListener("click", function () {
  code = codeInput.value
  const message = {
    action: "Join",
    code: code
  }
  socket.send(JSON.stringify(message))
})

createBtn.addEventListener("click", function () {
  const message = {
    action: "Create",
  }
  socket.send(JSON.stringify(message))
})

function playerJoined(data) {
  console.log("player joined with data: "+data)
}

function joinError(errors) {
  console.log("couldnt join because: "+errors)
}
