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

    PaddleUpPressed = 31,
    PaddleUpRelease = 32,

    PaddleDownPressed = 33,
    PaddleDownRelease = 34,

    PaddleMoved = 35,

    BallShot = 36,

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
    }

    draw_label(ctx: CanvasRenderingContext2D, x: number = this.label.x, y: number = this.label.y): void {
        ctx.font = this.label.font
        ctx.fillText(this.username, x, y);

    }

    get_label_width(ctx: CanvasRenderingContext2D): number {
        let old_font = ctx.font
        ctx.font = this.label.font
        ctx.font = old_font
        return ctx.measureText(this.username).width
    }
}

class Paddle {
    position: Point;
    length: number;
    width: number;
    speed: number;
    keys: Keys;
    pos_buffer: number[];
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
        this.pos_buffer = [];
        this.last_move_time = 0;
        this.last_interpolated_y = this.position.y;
    }

    draw(ctx: CanvasRenderingContext2D): void {
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

        this.position.y = Math.max(0, Math.min(this.position.y, canvas_height - this.length));
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
    radius: number;
    speed: number;
    direction: Direction;
    endAngle = 2 * Math.PI;

    constructor(
        position: Point,
        radius: number,
        speed: number,
        direction: Direction,
    ){
        this.position = position;
        this.radius = radius;
        this.speed = speed;
        this.direction = direction;
    }

    draw(ctx: CanvasRenderingContext2D): void {
        ctx.beginPath()
        ctx.arc(
            this.position.x,
            this.position.y,
            this.radius,
            0,
            this.endAngle
        );
        ctx.fill();
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

    offscreen: OffscreenCanvas;
    offscreen_ctx: OffscreenCanvasRenderingContext2D;


    constructor(
        code: string,
        ball: Ball,
        p1: Player,
        p2: Player,
        canvas: HTMLCanvasElement,
        offscreen: OffscreenCanvas,

        ctx: CanvasRenderingContext2D,
        offscreen_ctx: OffscreenCanvasRenderingContext2D,
    ){
        this._code = code;
        this.ball = ball;
        this.p1 = p1;
        this.p2 = p2;
        this.status = GameStatus.NotStarted;
        this.canvas = canvas
        this.offscreen = offscreen

        this.ctx = ctx
        this.offscreen_ctx = offscreen_ctx
    }

    set code(theCode: string) {
        this._code = theCode;
    }

    get code(): string {
        return this._code;
    }

    draw_fps(fps: number, x: number = 640, y:number = 360): void {
        this.ctx.font = "11px Arial";
        let fpsStr: string = fps.toString();
        let width = this.offscreen_ctx.measureText(fpsStr).width + 5;
        this.ctx?.fillText(fpsStr, x - width, y);
    }

    draw_code(x: number = 320, y: number = 20) {
        let width = this.offscreen_ctx.measureText(this._code).width + 5;
        this.ctx?.fillText(this._code, x - width, y);
    }

    swap_players() {
       const aux = this.p1.paddle
       this.p1.paddle = this.p2.paddle 
       this.p2.paddle = aux
       this.swap_players_position = true
    
       this.p1.paddle.speed = 300
       this.p2.paddle.speed = 200

       let label_aux = this.p1.label.x
       this.p1.label.x = this.p2.label.x
       this.p2.label.x = label_aux

       label_aux = this.p1.label.y
       this.p1.label.y = this.p2.label.y
       this.p2.label.y = label_aux
    }

    update_paddles(deltaTime: number) {
        this.p1.paddle.update(this.height, deltaTime)
        // this.p2.paddle.update(this.height, deltaTime)
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

function main(): void{
    const ctx = canvas.getContext("2d")
    
    const offscreen = new OffscreenCanvas(canvas_width, canvas_height);
    let offscreen_ctx = offscreen.getContext("2d")

    const paddle: Paddle = new Paddle({x: 30, y: canvas_height/2 }, 40, 4, 300)
    const paddle2: Paddle = new Paddle({x: canvas_width-30, y: canvas_height/2 }, 40, 4, 200)

    const player1: Player = new Player("User", paddle, true)
    const player2: Player = new Player("Waiting...", paddle2, false)
    
    const ball: Ball = new Ball({x: canvas_width/2, y: canvas_height/2}, 7, 6, Direction.Left); 

    if (ctx == null || offscreen_ctx == null) {
        console.log("Error: Could not put canvas context to work")
        return
    }
    game_state = new GameState("XXXX", ball, player1, player2, canvas, offscreen, ctx, offscreen_ctx)

    player1.label.x = 10
    player1.label.y = 20

    player2.label.x = game_state.width - player2.get_label_width(game_state.ctx) 
    player2.label.y = 20

    raf = requestAnimationFrame(animate)
}

function draw_state(state: GameState): void {
    state.ctx.clearRect(0, 0, state.width, state.height);

    state.ctx.fillStyle = "blue"

    state.p1.draw_label(state.ctx)
    state.p2.draw_label(state.ctx)

    state.draw_code()

    state.p1.paddle.draw(state.ctx)
    state.p2.paddle.draw(state.ctx)

    state.ctx.fillStyle = "yellow"
    state.ctx.strokeStyle = "1px black"


    state.ball.draw(state.ctx)

}

function update(deltaTime: any) {
    game_state.update_paddles(deltaTime)
}

window.addEventListener("keydown", async (event) => {
  if (game_state.status === GameStatus.Running) {
    if (event.key === "w") {
      game_state.p1.paddle.keys.up = true;
      let counter = 40
      while(game_state.p1.paddle.keys.up || counter > 0) {
          await sleep(ms)
          updatePaddlePosition();
          counter--
      }
    } else if (event.key === "s") {
      game_state.p1.paddle.keys.down = true;
      let counter = 40
      while(game_state.p1.paddle.keys.down || counter > 0) {
          await sleep(16)
          updatePaddlePosition();
          counter--
      }
    }

  }
});


let prevPaddleY = 360/2// Store initial position

function updatePaddlePosition() {
  // if (game_state.p1.paddle.position.y !== prevPaddleY) { 
    // Update event data only if position changed
    const ev: SocketEvent = {
      type: EventType.PaddleMoved,
      data: {
        y: game_state.p1.paddle.position.y,
      }
    };
    send_event(ev);
    // prevPaddleY = game_state.p1.paddle.position.y; // Store previous position
  }

  // Schedule next update (optional)
  // requestAnimationFrame(updatePaddlePosition);

// window.addEventListener("keypress", async (event) => {
//     if(event.key == "w" && game_state.status == GameStatus.Running) { 
//         game_state.p1.paddle.keys.up = true
//         const ev: SocketEvent = {
//             type: EventType.PaddleMoved,
//             data: {
//                 y: game_state.p1.paddle.position.y,
//             }
//         }
//         let counter = 20
//         while(game_state.p1.paddle.keys.up) {
//             if (counter == 0) {
//                 await sleep(200)
//             }
//             counter = 20
//             while(counter > 0) {
//                 ev.data.y = game_state.p1.paddle.position.y   
//                 send_event(ev)
//                 counter--;
//             }
//         }
//     }

//     if(event.key == "s" && game_state.status == GameStatus.Running) { 
//         game_state.p1.paddle.keys.down = true
//         const ev: SocketEvent = {
//             type: EventType.PaddleMoved,
//             data: {
//                 y: game_state.p1.paddle.position.y,
//             }
//         }
//         let counter = 20
//         while(game_state.p1.paddle.keys.down) {
//             if (counter == 0) {
//                 await sleep(200)
//             }
//             counter = 20
//             while(counter > 0) {
//                 ev.data.y = game_state.p1.paddle.position.y   
//                 send_event(ev)
//                 counter--;
//             }
//         }
//     }
// })

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

socket.addEventListener("open", (e) => {
    const msg_event: SocketEvent = {
        type: EventType.Message,
        data: {
            message: "Hello Server",
        }
    }
    send_event(msg_event)
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
        case EventType.PaddleUpPressed:
            handle_paddle_up_pressed(event)
            break;
        case EventType.PaddleDownPressed:
            handle_paddle_down_pressed(event)
            break;
        case EventType.PaddleUpRelease:
            handle_paddle_up_release(event)
            break;
        case EventType.PaddleDownRelease:
            handle_paddle_down_release(event)
            break;
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

function handle_pong(event: SocketEvent): void {
    if (event.data) {
        const endTime = performance.now();
        const responseTime = endTime - event.data.timestamp;
        console.clear()
        ms = responseTime
        console.log(`ms: ${ms}`);

        setTimeout(measure_latency, 4000);
    }
}

function handle_joined(event: SocketEvent): void { 
    if(event.data) {
        console.log(event)
        room_form.style.visibility = "hidden"
        var roomTitle = document.createElement("h1")
        roomTitle.classList.add("subtitle")
        roomTitle.classList.add("is-4")
        roomTitle.innerHTML = "Code: " + event.data.code    
        room_info_div.insertAdjacentElement("afterbegin", roomTitle)
        canvas.style.visibility = "visible"

        if (event.to) {
            game_state.p1.username = event.to
        }
        game_state.p2.isConnected = true
        game_state.p2.username = event.data.player        
        game_state.code = event.data.code
        game_state.swap_players()
        game_state.status = GameStatus.Running
    } else {
        console.log("no data", event)
    }
}

function handle_player_joined(event: SocketEvent): void {
    if (event.data) {
        game_state.p2.isConnected = true
        game_state.p2.username = event.data.player
        game_state.status = GameStatus.Running
    } else {
        console.log("no data", event)
    }
}

function handle_room_created(event: SocketEvent): void {
    game_state.code = event.data.code
    if (event.to) {
        game_state.p1.username = event.to
    }

    game_state.p1.isConnected = true
    room_form.style.visibility = "hidden"
    var roomTitle = document.createElement("h1")
    roomTitle.classList.add("subtitle")
    roomTitle.classList.add("is-4")
    roomTitle.innerHTML = "Code: " + event.data.code    
    room_info_div.insertAdjacentElement("afterbegin", roomTitle)
    canvas.style.visibility = "visible"
}

async function handle_paddle_up_pressed(event: SocketEvent): Promise<void> {
    // if (event.data) {
        // game_state.p2.paddle.position.y = event.data.y
    // }
    game_state.p2.paddle.keys.up = true
}


async function handle_paddle_down_pressed(event: SocketEvent): Promise<void> {
    // if (event.data){
        // game_state.p2.paddle.position.y = event.data.y
    // }
    game_state.p2.paddle.keys.down = true
}

async function handle_paddle_up_release(event: SocketEvent): Promise<void> {
    if (event.data){
        game_state.p2.paddle.position.y = event.data.y
    }
    game_state.p2.paddle.keys.up = false
}


async function handle_paddle_down_release(event: SocketEvent): Promise<void> {
    if (event.data){
        game_state.p2.paddle.position.y = event.data.y
    }
    game_state.p2.paddle.keys.down = false
}

function handle_paddle_move(event: SocketEvent): void {
    if(event.data){
        game_state.p2.paddle.move(event.data.y)
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
    const event: SocketEvent = {
        type: EventType.PaddleMoved,
        data: {
            y: game_state.p1.paddle.position.y,
        }
    }
    send_event(event);

}

function pressed_up(): void {
    const event: SocketEvent = {
        type: EventType.PaddleUpPressed,
        data: {
            y: game_state.p1.paddle.position.y,
        }
    }
    send_event(event);
}
function release_up(): void {
    const event: SocketEvent = {
        type: EventType.PaddleUpRelease,
        data: {
            y: game_state.p1.paddle.position.y,
        }
    }
    send_event(event);
}

function pressed_down(): void {
    const event: SocketEvent = {
        type: EventType.PaddleDownPressed,
        data: {
            y: game_state.p1.paddle.position.y,
        }
    }
    send_event(event);
}

function release_down(): void {
    const event: SocketEvent = {
        type: EventType.PaddleDownRelease,
        data: {
            y: game_state.p1.paddle.position.y,
        }
    }
    send_event(event);
}



function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

main()
