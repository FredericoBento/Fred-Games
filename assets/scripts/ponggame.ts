type Point = {
    x: number,
    y: number,
}

type Keys = {
    up: boolean,
    down: boolean,
}

type SocketEvent = {
    type: EventType,
    data?: any,
    room_code?: string,
    from?: string,
    to?: string,
    isError?: Boolean,
}

enum Direction {
    Up = 1,
    Down,
    Left,
    Right,
}

enum GameStatus {
    Running,
    Paused,
    NotStarted,
}

enum EventType {
    Ping = 98,
    Pong = 99,
    
    GameSettings = 0,
    Message = 1,

    CreateRoom = 21,
    CreatedRoom = 22,

    JoinRoom = 23,
    JoinedRoom = 24,
    PlayerJoinedRoom = 25,

    PaddleMoved = 35,

    BallShot = 36,
    BallUpdate = 37,
    Goal = 38,

    PlayerDisconnected = 4,
}

type Label = {
    x: number,
    y: number,
    content: string,
    font: string
}

class Player {
    username: string;
    paddle: Paddle;
    isConnected: boolean;
    label: Label;
    score: number;

    constructor(
        username: string = "Not Connected...",
        paddle: Paddle,
        isConnected: boolean = false,
    ) {
        this.username = username;
        this.paddle = paddle;
        this.isConnected = isConnected;
        this.label = {
            x: 0,
            y: 0,
            content: this.username,  
            font: "19px Arial"
        }
        this.score = 0
    }

    draw_label(ctx: CanvasRenderingContext2D, x: number = this.label.x, y: number = this.label.y): void {
        ctx.font = this.label.font
        if (this.isConnected) {
            ctx.fillStyle = "white"
        } else {
            ctx.fillStyle = "red"
        }
        ctx.fillText(this.username, x, y);

    }

    get_label_width(ctx: CanvasRenderingContext2D): number {
        let old_font = ctx.font
        ctx.font = this.label.font
        const w = ctx.measureText(this.username).width
        ctx.font = old_font
        return w
    }
}

class Paddle {
    position: Point;
    length: number;
    width: number;
    speed: number;
    keys: Keys;
    last_move_time: number;
    last_interpolated_y: number;

    constructor(
        position: Point,
        length: number,      
        width: number,
        speed: number,
    ) {
        this.position = position;
        this.length = length;
        this.width = width;
        this.speed = speed;
        this.keys = {
            up: false,
            down: false,
        }
        this.last_move_time = 0;
        this.last_interpolated_y = this.position.y;
    }

    draw(ctx: CanvasRenderingContext2D): void {
        ctx.fillStyle = "white"
        ctx.fillRect(
            this.position.x,
            this.position.y,
            this.width,
            this.length,
        );
    }

    update(canvas_height: number, deltaTime: number, speed: number = this.speed, keys: Keys = this.keys) {
        if (keys.up) {
            this.position.y -= speed * deltaTime
        }

        if (keys.down) {
            this.position.y += speed * deltaTime
        }

        this.position.y = Math.max(25, Math.min(this.position.y, (canvas_height - 25) - this.length));
        this.position.y  = Math.round(this.position.y)
    }

    move(y: number){
        const factor = 0.9
        const interpolatedY = this.lerp(this.last_interpolated_y, y, factor)
        this.last_interpolated_y = interpolatedY

        this.position.y = Math.max(0, Math.min(interpolatedY, canvas_height - this.length));
    }

    lerp(current_y: number, target_y: number, interpolation_factor: number): number {
        return current_y + (target_y - current_y) * interpolation_factor
    }
}


class Ball {
    position: Point;
    last_interpolated_position: Point;
    radius: number;
    endAngle = 2 * Math.PI;

    constructor(
        position: Point,
        radius: number,
        speed: number,
    ){
        this.last_interpolated_position = position;
        this.position = position;
        this.radius = radius;
    }

    draw(ctx: CanvasRenderingContext2D): void {
        ctx.strokeStyle = "black";
        ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.arc(
            this.position.x,
            this.position.y,
            this.radius,
            0,
            this.endAngle
        );
        ctx.closePath();
        ctx.fill();
    }

    center(width: number, height: number) {
        this.position.x = width / 2
        this.position.y = height/ 2
    }

    move(position: Point): void {
        const factor = 0.9
        const interpolatedY = this.lerp(this.last_interpolated_position.y, position.y, factor)
        this.last_interpolated_position.y = interpolatedY

        const interpolatedX = this.lerp(this.last_interpolated_position.x, position.x, factor)
        this.last_interpolated_position.x = interpolatedX
    }

    lerp(current: number, target: number, interpolation_factor: number): number {
        return current + (target - current) * interpolation_factor
    }
}

class GameState {
    _code: string;
    ball: Ball;
    p1: Player;
    p2: Player;
    status: GameStatus;
    swap_players_position: boolean = false;

    fps: number = 0;

    width: number = 640;
    height: number = 360;

    canvas: HTMLCanvasElement;
    ctx: CanvasRenderingContext2D;

    score_canvas: OffscreenCanvas;
    score_canvas_ctx: OffscreenCanvasRenderingContext2D;

    constructor(
        code: string,
        ball: Ball,
        p1: Player,
        p2: Player,
        canvas: HTMLCanvasElement,
        ctx: CanvasRenderingContext2D,
    ){
        this._code = code;
        this.ball = ball;
        this.p1 = p1;
        this.p2 = p2;
        this.status = GameStatus.NotStarted;
        this.canvas = canvas
        this.ctx = ctx

        this.score_canvas = new OffscreenCanvas(125, 25)
        const score_ctx = this.score_canvas.getContext("2d")
        if (score_ctx) {
            this.score_canvas_ctx = score_ctx
            this.update_scores()
        } else {
            throw new Error("Could not setup score_canvas context")
        }
    }

    set code(theCode: string) {
        this._code = theCode;
    }

    get code(): string {
        return this._code;
    }

    update_scores(): void{
        this.score_canvas_ctx.fillStyle = "black"
        this.score_canvas_ctx.fillRect(0, 0, 25, 25)
        this.score_canvas_ctx.fillRect(100, 0, 25, 25)

        // this.score_canvas_ctx.strokeStyle = "white"
        this.score_canvas_ctx.fillStyle = "white"

        this.score_canvas_ctx.font = "16px Arial"
        this.score_canvas_ctx.fillText(this.p1.score.toString(), 8, 17)
        this.score_canvas_ctx.fillText(this.p2.score.toString(), 109, 17)
    }

    draw_scores(): void {
        this.ctx.drawImage(this.score_canvas, (this.width / 2) - 65, 0)
    }

    draw_fps(fps: number, x: number = 640, y:number = 360): void {
        this.ctx.font = "11px Arial";
        let fpsStr: string = fps.toString();
        let width = this.ctx.measureText(fpsStr).width + 5;
        this.ctx?.fillText(fpsStr, x - width, y);
    }

    draw_ms(ms: number, x: number = 640, y:number = 360): void {
        this.ctx.font = "11px Arial";
        let fpsStr: string = ms.toFixed(0) + " ms";
        let width = this.ctx.measureText(fpsStr).width + 5;
        this.ctx?.fillText(fpsStr, x - width, y);
    }

    draw_code(x: number = 320, y: number = 20) {
        let width = this.ctx.measureText(this._code).width + 5;
        this.ctx?.fillText(this._code, x - width, y);
    }

    swap_players() {
       const aux = this.p1.paddle
       this.p1.paddle = this.p2.paddle 
       this.p2.paddle = aux
       this.swap_players_position = true

       let label_aux = this.p1.label.x
       this.p1.label.x = this.width - (this.p1.get_label_width(this.ctx) + 10)
       this.p2.label.x = label_aux

       label_aux = this.p1.label.y
       this.p1.label.y = this.p2.label.y
       this.p2.label.y = label_aux
    }

    update_paddle_p1(deltaTime: number) {
        this.p1.paddle.update(this.height, deltaTime)
    }

}

const canvas_width = 640
const canvas_height = 360

const socket = new WebSocket("ws://localhost:8080/ws/pong");

const room_form = document.getElementById("room-menu") as HTMLDivElement;
const room_info_div = document.getElementById("roomInfo") as HTMLDivElement;

const join_btn = document.getElementById("joinBtn") as HTMLButtonElement;
const create_btn = document.getElementById("createBtn") as HTMLButtonElement;

const code_input = document.getElementById("code") as HTMLInputElement;

const canvas = document.getElementById("gameCanvas") as HTMLCanvasElement;

let game_state: GameState;
let raf: number;
let ms: number = 52;

join_btn?.addEventListener("click", join_game);
create_btn?.addEventListener("click", create_room)


function measure_latency() {
    const event_ping = {
        type: EventType.Ping,
        data: {
            timestamp: performance.now().toString(),
        }
    }
    send_event(event_ping)
}

function main(): void {
    const ctx = canvas.getContext("2d")
    
    const paddle_y = (canvas_height / 2) - 20
    const paddle: Paddle = new Paddle({x: 30, y: paddle_y }, 40, 4, 300)
    const paddle2: Paddle = new Paddle({x: canvas_width-30, y: paddle_y }, 40, 4, 300)

    const player1: Player = new Player("", paddle, true)
    const player2: Player = new Player("", paddle2, false)
    
    const ball: Ball = new Ball({x: canvas_width/2, y: (canvas_height)/2}, 7, 6); 

    if (ctx == null) {
        console.log("Error: Could not put canvas context to work")
        return
    }
    game_state = new GameState("XXXX", ball, player1, player2, canvas, ctx)

    player1.label.x = 10
    player1.label.y = 20

    player2.label.x = game_state.width - (player2.get_label_width(game_state.ctx) + 10)
    player2.label.y = 20

    raf = requestAnimationFrame(animate)
}

function draw_state(state: GameState): void {
    state.ctx.fillStyle = "#36454F" 
    state.ctx.fillRect(0, 0, state.width, 25);

    state.ctx.fillRect(0, state.height - 25, state.width, 25);

    state.ctx.fillStyle = "#91A3B0"
    state.ctx.fillRect(0, 25, state.width, state.height - 50);

    state.ctx.fillStyle = "white"
    state.ctx.strokeStyle = "white"

    state.ctx.setLineDash([10, 5]);
    state.ctx.beginPath()
    state.ctx.lineTo(state.width / 2, 25)
    state.ctx.lineTo(state.width / 2, state.height - 25)
    state.ctx.stroke()


    state.p1.draw_label(state.ctx)
    state.p2.draw_label(state.ctx)

    state.draw_scores()

    // state.draw_code()

    state.p1.paddle.draw(state.ctx)
    state.p2.paddle.draw(state.ctx)

    state.ctx.fillStyle = "yellow"

    state.ball.draw(state.ctx)

    state.draw_ms(ms, state.width, state.height - 2)

}

function update(deltaTime: any) {
    game_state.update_paddle_p1(deltaTime)
}

window.addEventListener("keydown", async (event) => {
  if (game_state.status === GameStatus.Running) {
    if (event.key === "w") {
      game_state.p1.paddle.keys.up = true;
      let counter = 40
      while(game_state.p1.paddle.keys.up || counter > 0) {
          await sleep(20)
          paddle_moved();
          counter--
      }
    } else if (event.key === "s") {
      game_state.p1.paddle.keys.down = true;
      let counter = 40
      while(game_state.p1.paddle.keys.down || counter > 0) {
          await sleep(20)
          paddle_moved();
          counter--
      }
    }

  }
});

window.addEventListener("keydown", (event) => {
   if (event.key == "Space" || event.key == ' ') {
       console.log("Shot")
       const e: SocketEvent = {
           type:EventType.BallShot,
           
       }
       send_event(e)
   } 
});

window.addEventListener("keyup", (event) => {
    if(event.key == "w" && game_state.status == GameStatus.Running) { 
        game_state.p1.paddle.keys.up = false
        const ev: SocketEvent = {
            type: EventType.PaddleMoved,
            data: {
                y: game_state.p1.paddle.position.y,
            }
        }
        send_event(ev)
    }
    if(event.key == "s" && game_state.status == GameStatus.Running) { 
        game_state.p1.paddle.keys.down = false
        const ev: SocketEvent = {
            type: EventType.PaddleMoved,
            data: {
                y: game_state.p1.paddle.position.y,
            }
        }
        send_event(ev)
    }
})


let lastTime: number = 0;
var deltaTime: number;

function animate(time: any) {
    if (!time) {
        lastTime = time;
        window.requestAnimationFrame(animate)
        return
    }
    deltaTime = (time - lastTime) / 1000;

    update(deltaTime)
    draw_state(game_state)

    
    lastTime = time
    raf = window.requestAnimationFrame(animate);
}

socket.addEventListener("open", () => {
    measure_latency()
});

socket.addEventListener("message", (e) => {
    const event = parse_event(e.data)
    if (event.isError !== undefined) {
        if (event.isError == true) {
            handle_event_error(event)
            return
        }
    }       
    handle_event(event)
});

function parse_event(data: any): SocketEvent {
    const event_data = JSON.parse(data);
    const event: SocketEvent = event_data as SocketEvent;

    if (!event.type) {
        throw new Error("Receive event without type")
    }

    return event
}

function handle_event(event: SocketEvent): void {
    switch (event.type) {
        case EventType.Pong:
            handle_pong(event)
            break;
        case EventType.CreatedRoom:
            handle_room_created(event)
            break;
        case EventType.JoinedRoom:
            handle_joined(event)
            break;
        case EventType.PlayerJoinedRoom:
            handle_player_joined(event)
            break;
        case EventType.PaddleMoved:
            handle_paddle_move(event)
            break;
        case EventType.BallUpdate:
            handle_ball_update(event)
            break;
        case EventType.Goal:
            handle_goal(event);
            break;
        case EventType.PlayerDisconnected:
            handle_player_disconnect(event)
            break
        default:
            console.log("Unknown Event type: " + event.type)
            console.log(event)
            break
    }
}

function handle_event_error(event: SocketEvent): void {
    switch (event.type) {
        case EventType.JoinRoom:
            join_game_error(event)
            break
        case EventType.Ping:
            console.log("Could not ping")
            break
        default:
            console.log("Unknown Event type: " + event.type)
            console.log(event)
            break
    }
}

function send_event(ev: SocketEvent): void {
    socket.send(JSON.stringify(ev))
}

function handle_player_disconnect(event: SocketEvent): void {
    if (event.data) {
        console.log("player " + event.data.username + " has left")
        game_state.p2.isConnected = false
    }
}

function handle_pong(event: SocketEvent): void {
    if (event.data) {
        const endTime = performance.now();
        const responseTime = endTime - event.data.timestamp;
        ms = responseTime

        setTimeout(measure_latency, 4000);
    }
}

function handle_joined(event: SocketEvent): void { 
    if(event.data) {
        console.log(event)
        room_form.style.display = "none"
        var roomTitle = document.createElement("h1")
        roomTitle.classList.add("subtitle")
        roomTitle.classList.add("is-4")
        roomTitle.innerHTML = "Code: " + event.data.code    
        room_info_div.insertAdjacentElement("afterbegin", roomTitle)
        canvas.style.visibility = "visible"

        game_state.p1.username = event.data.username

        game_state.p2.isConnected = true
        game_state.p2.username = event.data.player        
        game_state.code = event.data.code
        if (event.data.is_player_1) {
            game_state.swap_players()
        } else {
            game_state.p2.label.x = game_state.width - game_state.p2.get_label_width(game_state.ctx) - 10
        }
        game_state.status = GameStatus.Running
    } else {
        console.log("no data", event)
    }
}

function handle_player_joined(event: SocketEvent): void {
    if (event.data) {
        game_state.p2.isConnected = true
        game_state.p2.username = event.data.player
        if (!event.data.is_player_1) {
            game_state.p2.label.x = game_state.width - game_state.p2.get_label_width(game_state.ctx) - 10
        }
        game_state.status = GameStatus.Running
    } else {
        console.log("no data", event)
    }
}

function handle_room_created(event: SocketEvent): void {
    game_state.code = event.data.code
    if (event.data) {
        game_state.p1.username = event.data.username
    }

    game_state.p1.isConnected = true
    room_form.style.display = "none"
    var roomTitle = document.createElement("h1")
    roomTitle.classList.add("subtitle")
    roomTitle.classList.add("is-4")
    roomTitle.innerHTML = "Code: " + event.data.code    
    navigator.clipboard.writeText(event.data.code);
    showNotification("code copied to clipboard")
    room_info_div.insertAdjacentElement("afterbegin", roomTitle)
    canvas.style.visibility = "visible"
}

function showNotification(message: string) {
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

function handle_paddle_move(event: SocketEvent): void {
    if(event.data){
        game_state.p2.paddle.move(event.data.y)
    }
}

function handle_ball_update(event: SocketEvent): void {
    if(event.data) {
        const pos: Point = {
            x: event.data.x,
            y: event.data.y,
        }
        game_state.ball.move(pos)
        // console.log(game_state.ball.position.x, pos.x)
    }    
}

function handle_goal(event: SocketEvent): void {
    if(event.data) {
        console.log(event.data)
        game_state.p1.score = event.data.player1_score
        game_state.p2.score = event.data.player2_score
        game_state.update_scores()
        game_state.ball.center(canvas_width, canvas_height)
    }
}

function create_room(): void {
    const event = {
        type: EventType.CreateRoom,
    }
    send_event(event)    
}

function join_game(): void {
    let code: string = code_input?.value
    if (code == "") {
        console.log("Empty code");
        return
    }
    const event: SocketEvent = {
        type: EventType.JoinRoom,
        data: {
            code: code,
        }
    }
    send_event(event)
}

function join_game_error(event: SocketEvent) {
    if (event.data){
        alert("Could not join room: " + event.data.message)
        return
    }
    alert("Could not join room")
}

function paddle_moved(): void {
    const ev: SocketEvent = {
        type: EventType.PaddleMoved,
        data: {
            y: game_state.p1.paddle.position.y,
        }
    };
    send_event(ev);
}

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

main()
