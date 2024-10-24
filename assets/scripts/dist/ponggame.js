"use strict";
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
class Player {
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
        ctx.fillText(this.username, x, y);
    }
    get_label_width(ctx) {
        let old_font = ctx.font;
        ctx.font = this.label_font;
        ctx.font = old_font;
        return ctx.measureText(this.username).width;
    }
}
class Paddle {
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
        ctx.fillRect(this.position.x, this.position.y, this.width, this.length);
    }
}
const endAngle = 2 * Math.PI;
class Ball {
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
        ctx.beginPath();
        ctx.arc(this.position.x, this.position.y, this.radius, 0, 
        // 2 * Math.PI
        endAngle);
        ctx.fill();
    }
}
class GameState {
    _code;
    ball;
    p1;
    p2;
    status;
    fps = 0;
    width = 640;
    height = 360;
    canvas;
    ctx;
    offscreen;
    offscreen_ctx;
    constructor(code, ball, p1, p2, canvas, offscreen, ctx, offscreen_ctx) {
        this._code = code;
        this.ball = ball;
        this.p1 = p1;
        this.p2 = p2;
        this.status = GameStatus.NotStarted;
        this.canvas = canvas;
        this.offscreen = offscreen;
        this.ctx = ctx;
        this.offscreen_ctx = offscreen_ctx;
    }
    set code(theCode) {
        this._code = theCode;
    }
    get code() {
        return this._code;
    }
    draw_fps(fps, x = 640, y = 360) {
        this.ctx.font = "11px Arial";
        let fpsStr = fps.toString();
        let width = this.offscreen_ctx.measureText(fpsStr).width + 5;
        this.ctx?.fillText(fpsStr, x - width, y);
    }
    draw_code(x = 320, y = 20) {
        let width = this.offscreen_ctx.measureText(this._code).width + 5;
        this.ctx?.fillText(this._code, x - width, y);
    }
}
const canvas_width = 640;
const canvas_height = 360;
const socket = new WebSocket("ws://localhost:8080/ws/pong");
const room_form = document.getElementById("room-menu");
const room_info_div = document.getElementById("roomInfo");
const join_btn = document.getElementById("joinBtn");
const create_btn = document.getElementById("createBtn");
const code_input = document.getElementById("code");
join_btn?.addEventListener("click", join_game);
create_btn?.addEventListener("click", create_room);
let keys = {
    up: false,
    down: false,
};
let swapped = false;
const canvas = document.getElementById("gameCanvas");
let game_state;
let raf;
function main() {
    // Get Settings from server
    // Set Ball
    // Set Player1 Object
    // then create layer
    const ctx = canvas.getContext("2d");
    const offscreen = new OffscreenCanvas(canvas_width, canvas_height);
    let offscreen_ctx = offscreen.getContext("2d");
    const paddle = new Paddle({ x: 30, y: canvas_height / 2 }, 40, 4, 5);
    const paddle2 = new Paddle({ x: canvas_width - 30, y: canvas_height / 2 }, 40, 4, 5);
    const player1 = new Player("My Name", paddle, true);
    const player2 = new Player("To be connected...", paddle2, false);
    const ball = new Ball({ x: canvas_width / 2, y: canvas_height / 2 }, 7, 6, Direction.Left);
    if (ctx == null || offscreen_ctx == null) {
        console.log("Error: Could not put canvas context to work");
        return;
    }
    game_state = new GameState("CCDD", ball, player1, player2, canvas, offscreen, ctx, offscreen_ctx);
    raf = requestAnimationFrame(animate);
}
function draw_state(state) {
    state.ctx.clearRect(0, 0, state.width, state.height);
    state.ctx.fillStyle = "blue";
    state.ctx.strokeStyle = "black";
    state.ctx.lineWidth = 2;
    if (!swapped) {
        state.p1.draw_label(state.ctx, 10, 20);
        state.p2.draw_label(state.ctx, state.width - state.p2.get_label_width(state.ctx) - 10, 20);
    }
    else {
        state.p1.draw_label(state.ctx, state.width - state.p1.get_label_width(state.ctx) - 10, 20);
        state.p2.draw_label(state.ctx, 10, 20);
    }
    state.draw_code();
    state.p1.paddle.draw(state.ctx);
    state.p2.paddle.draw(state.ctx);
    state.ctx.fillStyle = "yellow";
    state.ctx.strokeStyle = "1px black";
    state.ball.draw(state.ctx);
}
function update(deltaTime) {
    // const speed = game_state.p1.paddle.speed
    const speed = 300;
    if (keys.up) {
        game_state.p1.paddle.position.y -= speed * deltaTime;
        // game_state.p1.paddle.position.y -= game_state.p1.paddle.speed
    }
    if (keys.down) {
        game_state.p1.paddle.position.y += speed * deltaTime;
        // game_state.p1.paddle.position.y += game_state.p1.paddle.speed
    }
    game_state.p1.paddle.position.y = Math.max(0, Math.min(game_state.p1.paddle.position.y, canvas_height - game_state.p1.paddle.length));
    game_state.p1.paddle.position.y = Math.round(game_state.p1.paddle.position.y);
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
let lastTime = 0;
function animate(time) {
    if (!time) {
        lastTime = time;
        window.requestAnimationFrame(animate);
        return;
    }
    const deltaTime = (time - lastTime) / 1000;
    update(deltaTime);
    draw_state(game_state);
    lastTime = time;
    raf = window.requestAnimationFrame(animate);
}
socket.addEventListener("open", (e) => {
    const msg_event = {
        type: EventType.Message,
        data: {
            message: "Hello Server",
        }
    };
    send_event(msg_event);
});
socket.addEventListener("message", (e) => {
    const event = parse_event(e.data);
    if (event.isError !== undefined) {
        if (event.isError == true) {
            handle_event_error(event);
            return;
        }
    }
    handle_event(event);
});
function parse_event(data) {
    const event_data = JSON.parse(data);
    const event = event_data;
    if (!event.type) {
        throw new Error("Receive event without type");
    }
    return event;
}
function handle_event(event) {
    switch (event.type) {
        case EventType.CreatedRoom:
            handle_room_created(event);
            break;
        case EventType.JoinedRoom:
            handle_joined(event);
            break;
        case EventType.PlayerJoinedRoom:
            handle_player_joined(event);
            break;
        default:
            console.log("Unknown Event type: " + event.type);
            console.log(event);
            break;
    }
}
function handle_event_error(event) {
    switch (event.type) {
        case EventType.JoinRoom:
            join_game_error(event);
            break;
        default:
            console.log("Unknown Event type: " + event.type);
            console.log(event);
            break;
    }
}
function send_event(ev) {
    socket.send(JSON.stringify(ev));
}
function handle_other_player_joined(event) {
}
function handle_joined(event) {
    if (event.data) {
        console.log(event);
        room_form.style.visibility = "hidden";
        var roomTitle = document.createElement("h1");
        roomTitle.classList.add("subtitle");
        roomTitle.classList.add("is-4");
        roomTitle.innerHTML = "Code: " + event.data.code;
        room_info_div.insertAdjacentElement("afterbegin", roomTitle);
        canvas.style.visibility = "visible";
        if (event.to) {
            game_state.p1.username = event.to;
        }
        game_state.p1.isConnected = true;
        game_state.p2.isConnected = true;
        game_state.p2.username = event.data.player;
        game_state.code = event.data.code;
        swap_paddle_position();
    }
    else {
        console.log("no data", event);
    }
}
function handle_player_joined(event) {
    if (event.data) {
        console.log("player joined", event);
        game_state.p2.isConnected = true;
        game_state.p2.username = event.data.player;
    }
    else {
        console.log("no data", event);
    }
}
function handle_room_created(event) {
    game_state.code = event.data.code;
    if (event.to) {
        game_state.p1.username = event.to;
    }
    room_form.style.visibility = "hidden";
    var roomTitle = document.createElement("h1");
    roomTitle.classList.add("subtitle");
    roomTitle.classList.add("is-4");
    roomTitle.innerHTML = "Code: " + event.data.code;
    room_info_div.insertAdjacentElement("afterbegin", roomTitle);
    canvas.style.visibility = "visible";
}
function create_room() {
    const event = {
        type: EventType.CreateRoom,
    };
    send_event(event);
}
function join_game() {
    console.log("Clicked joined");
    let code = code_input?.value;
    if (code == "") {
        console.log("Empty code");
        return;
    }
    const event = {
        type: EventType.JoinRoom,
        data: {
            code: code,
        }
    };
    send_event(event);
}
function join_game_error(event) {
    if (event.data) {
        alert("Could not join room: " + event.data.message);
        return;
    }
    alert("Could not join room");
}
main();
function swap_paddle_position() {
    const aux = game_state.p1.paddle;
    game_state.p1.paddle = game_state.p2.paddle;
    game_state.p2.paddle = aux;
    swapped = true;
}
//# sourceMappingURL=ponggame.js.map