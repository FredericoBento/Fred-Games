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
        this.paddle = paddle
        this.isConnected = isConnected
    }
}

class Paddle {
    constructor(position, length, width) {
        this.position = position
        this.length = length
        this.width = width
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
const DirectionLeft = "direction_left"
const DirectionRight = "direction_right"

const socket = new WebSocket("ws://localhost:8080/ws/pong");

const joinBtn = document.getElementById("joinBtn")
const createBtn = document.getElementById("createBtn")
const codeInput = document.getElementById("code")

const canvasWidth = 640
const canvasHeight = 360

const p1Paddle = new Paddle(new Position(30, canvasHeight/2), 30, 4)
const p2Paddle = new Paddle(new Position(canvasWidth-30, canvasHeight/2), 30, 4)

const player1 = new Player(getCookie("User"), p1Paddle, true)
const player2 = new Player(getCookie("None"), p2Paddle, false)

const ball = new Ball(new Position(canvasWidth/2, canvasHeight/2), 3.00, 3, DirectionLeft)

const gameState = new GameState("empty", player1, player2, ball)

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

const app = new PIXI.Application();
const canvasDiv = document.getElementById("gameCanvas") 

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

  gameState.ball.drawing = ballDrawing
  gameState.p1.paddle.drawing = p1PaddleDrawing
  gameState.p2.paddle.drawing = p2PaddleDrawing
  

  app.stage.addChild(gameState.ball.drawing);
  app.stage.addChild(gameState.p1.paddle.drawing);
  app.stage.addChild(gameState.p2.paddle.drawing);

  let elapsed = 0.0;
  app.ticker.add((ticker) => {
    elapsed += ticker.deltaTime;
  });
}


initPixi()

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
