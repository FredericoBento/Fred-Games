class Position {
    constructor(x, y) {
       this.x = x
       this.y = y
    }

}

class GameState {
    constructor(code, p1, p2, ball) {
        this.code = code
        this.p1 = p1
        this.p2 = p2
        this.ball = ball
        this.hasStarted = false
    }
}

class Player {
    constructor(username, paddle, isConnected) {
        this.username = username
        this.usernameDrawing = null
        this.paddle = paddle
        this.isConnected = isConnected
    }
}

class Paddle {
    constructor(position, length, width, speed) {
        this.position = position
        this.length = length
        this.width = width
        this.speed = speed
        this.drawing = null
    }
}

class Ball {
    constructor(position, speed, radius, direction) {
      this.position = position
      this.speed = speed
      this.direction = direction
      this.radius = radius
      this.drawing = null
    }
}

const EventTypeMessage = 1
const EventTypeGameSettings = 1

const EventTypeCreateRoom = 21
const EventTypeCreatedRoom = 22

const EventTypeJoinRoom = 23
const EventTypeJoinedRoom	 = 24
const EventTypePlayerJoinedRoom = 25	

const EventTypePaddleUpPressed =31
const EventTypePaddleUpRelease = 32

const EventTypePaddleDownPressed = 33
const EventTypePaddleDownRelease = 34

const EventTypeBallShot = 35

const	EventTypePlayerDisconnected = 4

const DirectionLeft = "direction_left"
const DirectionRight = "direction_right"

const socket = new WebSocket("ws://localhost:8080/ws/pong");

const roomForm = document.getElementById("room-menu")
const roomInfoDiv = document.getElementById("roomInfo")

const joinBtn = document.getElementById("joinBtn")
const createBtn = document.getElementById("createBtn")
const codeInput = document.getElementById("code")

const canvasWidth = 640
const canvasHeight = 360

const gameState = new GameState()

const app = new PIXI.Application();
const canvasDiv = document.getElementById("gameCanvas") 
canvasDiv.style.visibility = "hidden"

function setup() {
    let p1Paddle = new Paddle(new Position(30, canvasHeight/2), 30, 4, 10.00)
    let p2Paddle = new Paddle(new Position(canvasWidth-30, canvasHeight/2), 30, 4, 10.00)

    let player1 = new Player(getCookie("User"), p1Paddle, true)
    let player2 = new Player("Waiting for Player...", p2Paddle, false)

    let ball = new Ball(new Position(canvasWidth/2, canvasHeight/2), 3.00, 3, DirectionLeft)

    gameState.code = "None"
    gameState.p1 = player1
    gameState.p2 = player2
    gameState.ball = ball

    initPixi()
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

socket.addEventListener("message", (e) => {
    eventData = JSON.parse(e.data)
    console.log(eventData)
    if (eventData.isError) {
        console.log("Event "+eventData.type + " had an error")
        return
    }
    switch (eventData.type) {
        case EventTypeCreatedRoom:
            room_created(eventData.data)
            break;
        case EventTypeJoinRoom:
            player_joined(eventData.data)
            break;
        case EventTypeJoinedRoom:
            joined(eventData)
            break;
        case EventTypePlayerJoinedRoom:
            player_joined(eventData.data)
            break;
        default:
            console.log("Unknown Event type: " + eventData.type)
    }
    
});

joinBtn.addEventListener("click", function () {
    code = codeInput.value
    const message = {
        code: code,
        player: getCookie("User"),
    }
    sendEvent(EventTypeJoinRoom, message)
})

createBtn.addEventListener("click", function () {
    sendEvent(EventTypeCreateRoom, null)
})

function room_created(data) {
    gameState.code = data.code
    gameState.p1.username = data.to

    roomForm.style.visibility = "hidden"
    var roomTitle = document.createElement("h1")
    roomTitle.classList.add("subtitle")
    roomTitle.classList.add("is-4")
    roomTitle.innerHTML = "Code: " + data.code    
    roomInfoDiv.insertAdjacentElement("afterbegin", roomTitle)
    canvasDiv.style.visibility = "visible"
}

function player_joined(data) {
  gameState.p2.username = data.player
  draw_usernames(gameState.p1, gameState.p2)
  console.log("player joined with data: "+data.player)
}

function joined(eventData) {
    gameState.code = eventData.code
    gameState.p1.username = eventData.to
    gameState.p2.username = eventData.data.player

    roomForm.style.visibility = "visible"
    var roomTitle = document.createElement("h1")
    roomTitle.classList.add("subtitle")
    roomTitle.classList.add("is-4")
    roomTitle.innerHTML = "Code: " + eventData.code    
    roomInfoDiv.insertAdjacentElement("afterbegin", roomTitle)
    canvasDiv.style.visibility = "visible"
}
function join_error(errors) {
    console.log("couldnt join because: "+errors)
}


async function initPixi() {
    await app.init({ width: canvasWidth, height: canvasHeight})
    canvasDiv.appendChild(app.canvas);

    let ballDrawing = new PIXI.Graphics()
                        .circle(gameState.ball.position.x, gameState.ball.position.y, gameState.ball.radius)
                        .fill("yellow")

    let p1PaddleDrawing = new PIXI.Graphics()
                        .rect(gameState.p1.paddle.position.x, 
                              gameState.p1.paddle.position.y, gameState.p1.paddle.width, gameState.p1.paddle.length)
                        .fill("red")

    let p2PaddleDrawing = new PIXI.Graphics()
                        .rect(gameState.p2.paddle.position.x, 
                              gameState.p2.paddle.position.y, gameState.p2.paddle.width, gameState.p2.paddle.length)
                        .fill("red")

    

    let p1UsernameDrawing = new PIXI.Text({text:gameState.p1.username, style: { fill: "white", fontSize: 18}})

    let p2UsernameDrawing = new PIXI.Text({text:gameState.p2.username, style: { fill: "white", fontSize: 18}})

    p2UsernameDrawing.anchor.set(1,0)
    p2UsernameDrawing.x = app.screen.width - 30
    p2UsernameDrawing.y = 30
    
    p1UsernameDrawing.anchor.set(1,0)
    p1UsernameDrawing.x = 70
    p1UsernameDrawing.y = 30
    

    gameState.ball.drawing = ballDrawing
    gameState.p1.paddle.drawing = p1PaddleDrawing
    gameState.p2.paddle.drawing = p2PaddleDrawing

    gameState.p1.usernameDrawing = p1UsernameDrawing
    gameState.p2.usernameDrawing = p2UsernameDrawing
  

    app.stage.addChild(gameState.ball.drawing);

    app.stage.addChild(gameState.p1.paddle.drawing);
    app.stage.addChild(gameState.p2.paddle.drawing);

    app.stage.addChild(gameState.p1.usernameDrawing);
    app.stage.addChild(gameState.p2.usernameDrawing);

    let elapsed = 0.0;
    app.ticker.add((ticker) => {
        elapsed += ticker.deltaTime;
        if (keys.up) {
            gameState.p1.paddle.position.y -= gameState.p1.paddle.speed * ticker.deltaTime;
        }

        if (keys.down) {
            gameState.p1.paddle.position.y += gameState.p1.paddle.speed * ticker.deltaTime;
        }

        gameState.p1.paddle.position.y = Math.max(0, Math.min(gameState.p1.paddle.position.y, canvasHeight - gameState.p1.paddle.length));

        draw(gameState)
  });
}


function draw(gameState) {  
    draw_ball(gameState.ball)
    draw_paddle(gameState.p1.paddle)
    draw_paddle(gameState.p2.paddle)
    // draw_usernames(gameState.p1, gameState.p2)
}

function draw_ball(ball) {
    ball.drawing.clear()
                .beginFill(0xFFFF00)
                .drawCircle(ball.position.x, ball.position.y, ball.radius)
                .endFill()
}

function draw_paddle(paddle) {
    paddle.drawing.clear()
                 .beginFill(0xFF0000) 
                 .drawRect(paddle.position.x, paddle.position.y, paddle.width, paddle.length)
                 .endFill()
}

function draw_usernames(player1, player2) {
    // player1.usernameDrawing.text = player1.username
    // player2.usernameDrawing.text = player2.username

    // player1.usernameDrawing.clear() 
    player1.usernameDrawing = new PIXI.Text({text:player1.username, style: { fill: "white"}})

    player2.usernameDrawing.clear() 
    player2.usernameDrawing = new PIXI.Text({text:player2.username, style: { align: "right", fill: "white"}})

    
}

let keys = {
    up: false,
    down: false,
}

window.addEventListener("keydown", (event) => {
    if (event.key === "w") keys.up = true;
    if (event.key === "s") keys.down = true;
});

window.addEventListener("keyup", (event) => {
    if (event.key === "w") keys.up = false;
    if (event.key === "s") keys.down = false;
});


function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

setup()
