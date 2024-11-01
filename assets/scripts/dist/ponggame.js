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
    EventType[EventType["Ping"] = 98] = "Ping";
    EventType[EventType["Pong"] = 99] = "Pong";
    EventType[EventType["GameSettings"] = 0] = "GameSettings";
    EventType[EventType["Message"] = 1] = "Message";
    EventType[EventType["CreateRoom"] = 21] = "CreateRoom";
    EventType[EventType["CreatedRoom"] = 22] = "CreatedRoom";
    EventType[EventType["JoinRoom"] = 23] = "JoinRoom";
    EventType[EventType["JoinedRoom"] = 24] = "JoinedRoom";
    EventType[EventType["PlayerJoinedRoom"] = 25] = "PlayerJoinedRoom";
    EventType[EventType["PaddleMoved"] = 35] = "PaddleMoved";
    EventType[EventType["BallShot"] = 36] = "BallShot";
    EventType[EventType["BallUpdate"] = 37] = "BallUpdate";
    EventType[EventType["Goal"] = 38] = "Goal";
    EventType[EventType["PlayerDisconnected"] = 4] = "PlayerDisconnected";
})(EventType || (EventType = {}));
class Player {
    username;
    paddle;
    isConnected;
    label;
    score;
    constructor(username = "Not Connected...", paddle, isConnected = false) {
        this.username = username;
        this.paddle = paddle;
        this.isConnected = isConnected;
        this.label = {
            x: 0,
            y: 0,
            content: this.username,
            font: "19px Arial"
        };
        this.score = 0;
    }
    draw_label(ctx, x = this.label.x, y = this.label.y) {
        ctx.font = this.label.font;
        if (this.isConnected) {
            ctx.fillStyle = "white";
        }
        else {
            ctx.fillStyle = "red";
        }
        ctx.fillText(this.username, x, y);
    }
    get_label_width(ctx) {
        let old_font = ctx.font;
        ctx.font = this.label.font;
        const w = ctx.measureText(this.username).width;
        ctx.font = old_font;
        return w;
    }
}
class Paddle {
    position;
    length;
    width;
    speed;
    keys;
    last_move_time;
    last_interpolated_y;
    constructor(position, length, width, speed) {
        this.position = position;
        this.length = length;
        this.width = width;
        this.speed = speed;
        this.keys = {
            up: false,
            down: false,
        };
        this.last_move_time = 0;
        this.last_interpolated_y = this.position.y;
    }
    draw(ctx) {
        ctx.fillStyle = "white";
        ctx.fillRect(this.position.x, this.position.y, this.width, this.length);
    }
    update(canvas_height, deltaTime, speed = this.speed, keys = this.keys) {
        if (keys.up) {
            this.position.y -= speed * deltaTime;
        }
        if (keys.down) {
            this.position.y += speed * deltaTime;
        }
        this.position.y = Math.max(25, Math.min(this.position.y, (canvas_height - 25) - this.length));
        this.position.y = Math.round(this.position.y);
    }
    move(y) {
        const factor = 0.9;
        const interpolatedY = this.lerp(this.last_interpolated_y, y, factor);
        this.last_interpolated_y = interpolatedY;
        this.position.y = Math.max(0, Math.min(interpolatedY, canvas_height - this.length));
    }
    lerp(current_y, target_y, interpolation_factor) {
        return current_y + (target_y - current_y) * interpolation_factor;
    }
}
class Ball {
    position;
    last_interpolated_position;
    radius;
    speed;
    dx;
    dy;
    endAngle = 2 * Math.PI;
    constructor(position, radius, speed, dx, dy) {
        this.last_interpolated_position = position;
        this.position = position;
        this.radius = radius;
        this.speed = speed;
        this.dx = dx;
        this.dy = dy;
    }
    draw(ctx) {
        ctx.strokeStyle = "black";
        ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.arc(this.position.x, this.position.y, this.radius, 0, this.endAngle);
        ctx.closePath();
        ctx.fill();
    }
    center(width, height) {
        this.position.x = width / 2;
        this.position.y = height / 2;
    }
    move(position) {
        const factor = 0.9;
        const interpolatedY = this.lerp(this.last_interpolated_position.y, position.y, factor);
        this.last_interpolated_position.y = interpolatedY;
        const interpolatedX = this.lerp(this.last_interpolated_position.x, position.x, factor);
        this.last_interpolated_position.x = interpolatedX;
        // this.position.y = Math.max(0, Math.min(interpolatedY, canvas_height - this.radius));
        // this.position.x = Math.max(0, Math.min(interpolatedY, canvas_height - this.radius));
    }
    lerp(current, target, interpolation_factor) {
        return current + (target - current) * interpolation_factor;
    }
}
class GameState {
    _code;
    ball;
    p1;
    p2;
    status;
    swap_players_position = false;
    fps = 0;
    width = 640;
    height = 360;
    canvas;
    ctx;
    score_canvas;
    score_canvas_ctx;
    constructor(code, ball, p1, p2, canvas, ctx) {
        this._code = code;
        this.ball = ball;
        this.p1 = p1;
        this.p2 = p2;
        this.status = GameStatus.NotStarted;
        this.canvas = canvas;
        this.ctx = ctx;
        this.score_canvas = new OffscreenCanvas(125, 25);
        const score_ctx = this.score_canvas.getContext("2d");
        if (score_ctx) {
            this.score_canvas_ctx = score_ctx;
            this.update_scores();
        }
        else {
            throw new Error("Could not setup score_canvas context");
        }
    }
    set code(theCode) {
        this._code = theCode;
    }
    get code() {
        return this._code;
    }
    update_scores() {
        this.score_canvas_ctx.fillStyle = "black";
        this.score_canvas_ctx.fillRect(0, 0, 25, 25);
        this.score_canvas_ctx.fillRect(100, 0, 25, 25);
        // this.score_canvas_ctx.strokeStyle = "white"
        this.score_canvas_ctx.fillStyle = "white";
        this.score_canvas_ctx.font = "16px Arial";
        this.score_canvas_ctx.fillText(this.p1.score.toString(), 8, 17);
        this.score_canvas_ctx.fillText(this.p2.score.toString(), 109, 17);
    }
    draw_scores() {
        this.ctx.drawImage(this.score_canvas, (this.width / 2) - 65, 0);
    }
    draw_fps(fps, x = 640, y = 360) {
        this.ctx.font = "11px Arial";
        let fpsStr = fps.toString();
        let width = this.ctx.measureText(fpsStr).width + 5;
        this.ctx?.fillText(fpsStr, x - width, y);
    }
    draw_ms(ms, x = 640, y = 360) {
        this.ctx.font = "11px Arial";
        let fpsStr = ms.toFixed(0) + " ms";
        let width = this.ctx.measureText(fpsStr).width + 5;
        this.ctx?.fillText(fpsStr, x - width, y);
    }
    draw_code(x = 320, y = 20) {
        let width = this.ctx.measureText(this._code).width + 5;
        this.ctx?.fillText(this._code, x - width, y);
    }
    swap_players() {
        const aux = this.p1.paddle;
        this.p1.paddle = this.p2.paddle;
        this.p2.paddle = aux;
        this.swap_players_position = true;
        let label_aux = this.p1.label.x;
        this.p1.label.x = this.width - (this.p1.get_label_width(this.ctx) + 10);
        this.p2.label.x = label_aux;
        label_aux = this.p1.label.y;
        this.p1.label.y = this.p2.label.y;
        this.p2.label.y = label_aux;
    }
    update_paddle_p1(deltaTime) {
        this.p1.paddle.update(this.height, deltaTime);
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
const canvas = document.getElementById("gameCanvas");
let game_state;
let raf;
let ms = 52;
join_btn?.addEventListener("click", join_game);
create_btn?.addEventListener("click", create_room);
function measure_latency() {
    const event_ping = {
        type: EventType.Ping,
        data: {
            timestamp: performance.now().toString(),
        }
    };
    send_event(event_ping);
}
function main() {
    const ctx = canvas.getContext("2d");
    const offscreen = new OffscreenCanvas(canvas_width, canvas_height);
    let offscreen_ctx = offscreen.getContext("2d");
    const paddle_y = (canvas_height / 2) - 20;
    const paddle = new Paddle({ x: 30, y: paddle_y }, 40, 4, 300);
    const paddle2 = new Paddle({ x: canvas_width - 30, y: paddle_y }, 40, 4, 300);
    const player1 = new Player("", paddle, true);
    const player2 = new Player("", paddle2, false);
    const ball = new Ball({ x: canvas_width / 2, y: (canvas_height) / 2 }, 7, 6, 0, 0);
    // ball.center(canvas_width, canvas_height)
    if (ctx == null || offscreen_ctx == null) {
        console.log("Error: Could not put canvas context to work");
        return;
    }
    game_state = new GameState("XXXX", ball, player1, player2, canvas, ctx);
    player1.label.x = 10;
    player1.label.y = 20;
    player2.label.x = game_state.width - (player2.get_label_width(game_state.ctx) + 10);
    player2.label.y = 20;
    raf = requestAnimationFrame(animate);
}
function draw_state(state) {
    state.ctx.fillStyle = "#36454F";
    state.ctx.fillRect(0, 0, state.width, 25);
    state.ctx.fillRect(0, state.height - 25, state.width, 25);
    state.ctx.fillStyle = "#91A3B0";
    state.ctx.fillRect(0, 25, state.width, state.height - 50);
    state.ctx.fillStyle = "white";
    state.ctx.strokeStyle = "white";
    state.ctx.setLineDash([10, 5]);
    state.ctx.beginPath();
    state.ctx.lineTo(state.width / 2, 25);
    state.ctx.lineTo(state.width / 2, state.height - 25);
    state.ctx.stroke();
    state.p1.draw_label(state.ctx);
    state.p2.draw_label(state.ctx);
    state.draw_scores();
    // state.draw_code()
    state.p1.paddle.draw(state.ctx);
    state.p2.paddle.draw(state.ctx);
    state.ctx.fillStyle = "yellow";
    state.ball.draw(state.ctx);
    state.draw_ms(ms, state.width, state.height - 2);
}
function update(deltaTime) {
    game_state.update_paddle_p1(deltaTime);
}
window.addEventListener("keydown", async (event) => {
    if (game_state.status === GameStatus.Running) {
        if (event.key === "w") {
            game_state.p1.paddle.keys.up = true;
            let counter = 40;
            while (game_state.p1.paddle.keys.up || counter > 0) {
                await sleep(20);
                paddle_moved();
                counter--;
            }
        }
        else if (event.key === "s") {
            game_state.p1.paddle.keys.down = true;
            let counter = 40;
            while (game_state.p1.paddle.keys.down || counter > 0) {
                await sleep(20);
                paddle_moved();
                counter--;
            }
        }
    }
});
window.addEventListener("keydown", (event) => {
    if (event.key == "Space" || event.key == ' ') {
        console.log("Shot");
        const e = {
            type: EventType.BallShot,
        };
        send_event(e);
    }
});
window.addEventListener("keyup", (event) => {
    if (event.key == "w" && game_state.status == GameStatus.Running) {
        game_state.p1.paddle.keys.up = false;
        const ev = {
            type: EventType.PaddleMoved,
            data: {
                y: game_state.p1.paddle.position.y,
            }
        };
        send_event(ev);
    }
    if (event.key == "s" && game_state.status == GameStatus.Running) {
        game_state.p1.paddle.keys.down = false;
        const ev = {
            type: EventType.PaddleMoved,
            data: {
                y: game_state.p1.paddle.position.y,
            }
        };
        send_event(ev);
    }
});
let lastTime = 0;
var deltaTime;
function animate(time) {
    if (!time) {
        lastTime = time;
        window.requestAnimationFrame(animate);
        return;
    }
    deltaTime = (time - lastTime) / 1000;
    update(deltaTime);
    draw_state(game_state);
    lastTime = time;
    raf = window.requestAnimationFrame(animate);
}
socket.addEventListener("open", () => {
    measure_latency();
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
        case EventType.Pong:
            handle_pong(event);
            break;
        case EventType.CreatedRoom:
            handle_room_created(event);
            break;
        case EventType.JoinedRoom:
            handle_joined(event);
            break;
        case EventType.PlayerJoinedRoom:
            handle_player_joined(event);
            break;
        case EventType.PaddleMoved:
            handle_paddle_move(event);
            break;
        case EventType.BallUpdate:
            handle_ball_update(event);
            break;
        case EventType.Goal:
            handle_goal(event);
            break;
        case EventType.PlayerDisconnected:
            handle_player_disconnect(event);
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
        case EventType.Ping:
            console.log("Could not ping");
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
function handle_player_disconnect(event) {
    if (event.data) {
        console.log("player " + event.data.username + " has left");
        game_state.p2.isConnected = false;
    }
}
function handle_pong(event) {
    if (event.data) {
        const endTime = performance.now();
        const responseTime = endTime - event.data.timestamp;
        ms = responseTime;
        setTimeout(measure_latency, 4000);
    }
}
function handle_joined(event) {
    if (event.data) {
        console.log(event);
        room_form.style.display = "none";
        var roomTitle = document.createElement("h1");
        roomTitle.classList.add("subtitle");
        roomTitle.classList.add("is-4");
        roomTitle.innerHTML = "Code: " + event.data.code;
        room_info_div.insertAdjacentElement("afterbegin", roomTitle);
        canvas.style.visibility = "visible";
        if (event.to) {
            game_state.p1.username = event.to;
        }
        game_state.p2.isConnected = true;
        game_state.p2.username = event.data.player;
        game_state.code = event.data.code;
        game_state.swap_players();
        game_state.status = GameStatus.Running;
    }
    else {
        console.log("no data", event);
    }
}
function handle_player_joined(event) {
    if (event.data) {
        game_state.p2.isConnected = true;
        game_state.p2.username = event.data.player;
        game_state.p2.label.x = game_state.width - game_state.p2.get_label_width(game_state.ctx) - 10;
        game_state.status = GameStatus.Running;
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
    game_state.p1.isConnected = true;
    room_form.style.display = "none";
    var roomTitle = document.createElement("h1");
    roomTitle.classList.add("subtitle");
    roomTitle.classList.add("is-4");
    roomTitle.innerHTML = "Code: " + event.data.code;
    navigator.clipboard.writeText(event.data.code);
    showNotification("code copied to clipboard");
    room_info_div.insertAdjacentElement("afterbegin", roomTitle);
    canvas.style.visibility = "visible";
}
function showNotification(message) {
    // Create the notification element
    const notification = document.createElement("div");
    notification.innerText = message;
    notification.style.position = "fixed";
    notification.style.bottom = "20px";
    notification.style.left = "50%";
    notification.style.transform = "translateX(-50%)";
    notification.style.backgroundColor = "#333";
    notification.style.color = "#fff";
    notification.style.padding = "10px 20px";
    notification.style.borderRadius = "5px";
    notification.style.opacity = "0";
    notification.style.transition = "opacity 0.3s";
    notification.style.fontSize = "12px";
    // Add the notification to the document body
    document.body.appendChild(notification);
    // Show the notification
    setTimeout(() => notification.style.opacity = "1", 10);
    // Hide and remove the notification after 2 seconds
    setTimeout(() => {
        notification.style.opacity = "0";
        setTimeout(() => notification.remove(), 300); // 300ms to match the transition
    }, 2000);
}
function handle_paddle_move(event) {
    if (event.data) {
        game_state.p2.paddle.move(event.data.y);
    }
}
function handle_ball_update(event) {
    if (event.data) {
        const pos = {
            x: event.data.x,
            y: event.data.y,
        };
        game_state.ball.move(pos);
        // console.log(game_state.ball.position.x, pos.x)
    }
}
function handle_goal(event) {
    if (event.data) {
        console.log(event.data);
        game_state.p1.score = event.data.player1_score;
        game_state.p2.score = event.data.player2_score;
        game_state.update_scores();
        game_state.ball.center(canvas_width, canvas_height);
    }
}
function create_room() {
    const event = {
        type: EventType.CreateRoom,
    };
    send_event(event);
}
function join_game() {
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
function paddle_moved() {
    const ev = {
        type: EventType.PaddleMoved,
        data: {
            y: game_state.p1.paddle.position.y,
        }
    };
    send_event(ev);
}
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
main();
//# sourceMappingURL=ponggame.js.map