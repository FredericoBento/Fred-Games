var Direction;
(function (Direction) {
    Direction[Direction["Up"] = 1] = "Up";
    Direction[Direction["Down"] = 2] = "Down";
    Direction[Direction["Left"] = 3] = "Left";
    Direction[Direction["Right"] = 4] = "Right";
})(Direction || (Direction = {}));
var GameStatus;
(function (GameStatus) {
    GameStatus[GameStatus["Running"] = 0] = "Running";
    GameStatus[GameStatus["Paused"] = 1] = "Paused";
    GameStatus[GameStatus["NotStarted"] = 2] = "NotStarted";
})(GameStatus || (GameStatus = {}));
export class Player {
    username;
    paddle;
    isConnected;
    usernameLabel;
    label_font;
    constructor(username = "Not Connected...", paddle, isConnected = false) {
        this.username = username;
        this.paddle = paddle;
        this.isConnected = isConnected;
        this.label_font = "19px Arial";
    }
    draw_label(ctx, x, y) {
        ctx.font = this.label_font;
        ctx.fillStyle = "white";
        ctx.fillText(this.username, x, y);
    }
    get_label_width(ctx) {
        let old_font = ctx.font;
        ctx.font = this.label_font;
        ctx.font = old_font;
        return ctx.measureText(this.username).width;
    }
}
export class Paddle {
    position;
    length;
    width;
    speed;
    constructor(position, length, width, speed) {
        this.position = position;
        this.length = length;
        this.width = width;
        this.speed = speed;
    }
    draw(ctx) {
        ctx.fillStyle = "blue";
        ctx.fillRect(this.position.x, this.position.y, this.width, this.length);
    }
}
export class Ball {
    position;
    radius;
    speed;
    direction;
    constructor(position, radius, speed, direction) {
        this.position = position;
        this.radius = radius;
        this.speed = speed;
        this.direction = direction;
    }
    draw(ctx) {
        ctx.fillStyle = "yellow";
        ctx.arc(this.position.x, this.position.y, this.radius, 0, 2 * Math.PI);
        ctx.fill();
    }
}
export class GameState {
    code;
    ball;
    p1;
    p2;
    status;
    constructor(code, ball, p1, p2) {
        this.code = code;
        this.ball = ball;
        this.p1 = p1;
        this.p2 = p2;
        this.status = GameStatus.NotStarted;
    }
}
var EventType;
(function (EventType) {
    EventType[EventType["GameSettings"] = 0] = "GameSettings";
    EventType[EventType["Message"] = 1] = "Message";
    EventType[EventType["CreateRoom"] = 21] = "CreateRoom";
    EventType[EventType["CreatedRoom"] = 22] = "CreatedRoom";
    EventType[EventType["JoinRoom"] = 23] = "JoinRoom";
    EventType[EventType["JoinedRoom"] = 24] = "JoinedRoom";
    EventType[EventType["PlayerJoinedRoom"] = 25] = "PlayerJoinedRoom";
    EventType[EventType["PaddleUpPressed"] = 31] = "PaddleUpPressed";
    EventType[EventType["PaddleUpRelease"] = 32] = "PaddleUpRelease";
    EventType[EventType["PaddleDownPressed"] = 33] = "PaddleDownPressed";
    EventType[EventType["PaddleDownRelease"] = 34] = "PaddleDownRelease";
    EventType[EventType["BallShot"] = 35] = "BallShot";
    EventType[EventType["PlayerDisconnected"] = 4] = "PlayerDisconnected";
})(EventType || (EventType = {}));
const canvas_width = 640;
const canvas_height = 360;
// const canvas_div = document.getElementById("canvasDiv") as HTMLDivElement;
const canvas = document.getElementById("gameCanvas");
const ctx = canvas.getContext("2d");
const socket = new WebSocket("ws://localhost:8080/ws/pong");
const room_form = document.getElementById("room-menu");
const room_info_div = document.getElementById("roomInfo");
const join_btn = document.getElementById("joinBtn");
const create_btn = document.getElementById("createBtn");
const code_input = document.getElementById("code");
join_btn?.addEventListener("click", join_game);
let keys = {
    up: false,
    down: false,
};
let game_state;
let lastTimestamp;
function main() {
    // Get Settings from server
    // Set Ball
    // Set Player1 Object
    // then create layer
    const paddle = new Paddle({ x: 30, y: canvas_height / 2 }, 40, 4, 200);
    const paddle2 = new Paddle({ x: canvas_width - 30, y: canvas_height / 2 }, 40, 4, 200);
    const player1 = new Player("My Name", paddle, true);
    const player2 = new Player("To be connected...", paddle2, false);
    const ball = new Ball({ x: canvas_width / 2, y: canvas_height / 2 }, 7, 6, Direction.Left);
    game_state = new GameState("CCDD", ball, player1, player2);
    lastTimestamp = 0;
    window.requestAnimationFrame(animate);
}
function draw_state(state) {
    if (ctx == null) {
        return false;
    }
    // window.requestAnimationFrame(update);
    ctx.clearRect(0, 0, canvas_width, canvas_height);
    ctx.fillStyle = "black";
    ctx.fillRect(0, 0, canvas_width, canvas_height);
    ctx.save();
    state.ball.draw(ctx);
    state.p1.paddle.draw(ctx);
    state.p2.paddle.draw(ctx);
    state.p1.draw_label(ctx, 5, 25);
    state.p2.draw_label(ctx, canvas_width - state.p2.get_label_width(ctx), 20);
    return true;
}
function update(deltaTime) {
    if (keys.up) {
        game_state.p1.paddle.position.y -= game_state.p1.paddle.speed * deltaTime;
    }
    if (keys.down) {
        game_state.p1.paddle.position.y += game_state.p1.paddle.speed * deltaTime;
    }
    game_state.p1.paddle.position.y = Math.max(0, Math.min(game_state.p1.paddle.position.y, canvas_height - game_state.p1.paddle.length));
}
window.addEventListener("keydown", (event) => {
    if (event.key == "w") {
        keys.up = true;
    }
    if (event.key == "s") {
        keys.down = true;
    }
});
window.addEventListener("keyup", (event) => {
    if (event.key == "w") {
        keys.up = false;
    }
    if (event.key == "s") {
        keys.down = false;
    }
});
function animate(timestamp) {
    if (!lastTimestamp) {
        lastTimestamp = timestamp;
        window.requestAnimationFrame(animate);
        return;
    }
    let deltaTime = timestamp - lastTimestamp;
    deltaTime = deltaTime / 1000;
    update(deltaTime);
    draw_state(game_state);
    lastTimestamp = timestamp;
    window.requestAnimationFrame(animate);
}
function handle_event() {
}
function send_event() {
}
function join_game() {
    console.log("Clicked joined");
    let code = code_input?.value;
    if (code == "") {
        console.log("abc");
    }
}
main();
//# sourceMappingURL=ponggame.js.map