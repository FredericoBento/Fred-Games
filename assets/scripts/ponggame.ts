type Point = {
    x: number,
    y: number,
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

    BallShot = 35,

    PlayerDisconnected = 4,
}

class Player {
    username: string;
    paddle: Paddle;
    isConnected: boolean;
    usernameLabel: any;
    label_font: string; 


    constructor(
        username: string = "Not Connected...",
        paddle: Paddle,
        isConnected: boolean = false,
    ) {
        this.username = username;
        this.paddle = paddle;
        this.isConnected = isConnected;
        this.label_font = "19px Arial"; 
    }

    draw_label(ctx: CanvasRenderingContext2D, x: number, y: number): void {
        ctx.font = this.label_font
        ctx.fillText(this.username, x, y);

    }

    get_label_width(ctx: CanvasRenderingContext2D): number {
        let old_font = ctx.font
        ctx.font = this.label_font
        ctx.font = old_font
        return ctx.measureText(this.username).width
    }
}

class Paddle {
    position: Point;
    length: number;
    width: number;
    speed: number;

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
    }

    draw(ctx: CanvasRenderingContext2D): void {
        ctx.fillRect(
            this.position.x,
            this.position.y,
            this.width,
            this.length,
        );
    }
}

const endAngle = 2 * Math.PI

class Ball {
    position: Point;
    radius: number;
    speed: number;
    direction: Direction;

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
            // 2 * Math.PI
            endAngle
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


}

const canvas_width = 640
const canvas_height = 360

const socket = new WebSocket("ws://localhost:8080/ws/pong");

const room_form = document.getElementById("room-menu") as HTMLDivElement;
const room_info_div = document.getElementById("roomInfo") as HTMLDivElement;

const join_btn = document.getElementById("joinBtn") as HTMLButtonElement;
const create_btn = document.getElementById("createBtn") as HTMLButtonElement;

const code_input = document.getElementById("code") as HTMLInputElement;

join_btn?.addEventListener("click", join_game);
create_btn?.addEventListener("click", create_room)

let keys = {
    up: false,
    down: false,
    up2: false,
    down2: false
}

let swapped = false

const canvas = document.getElementById("gameCanvas") as HTMLCanvasElement;
let game_state: GameState;

let raf: number;

function main(): void{
    // Get Settings from server
    // Set Ball
    // Set Player1 Object
    // then create layer
    const ctx = canvas.getContext("2d")
    
    const offscreen = new OffscreenCanvas(canvas_width, canvas_height);
    let offscreen_ctx = offscreen.getContext("2d")

    const paddle: Paddle = new Paddle({x: 30, y: canvas_height/2 }, 40, 4, 5)
    const paddle2: Paddle = new Paddle({x: canvas_width-30, y: canvas_height/2 }, 40, 4, 5)

    const player1: Player = new Player("My Name", paddle, true)
    const player2: Player = new Player("To be connected...", paddle2, false)

    const ball: Ball = new Ball({x: canvas_width/2, y: canvas_height/2}, 7, 6, Direction.Left); 

    if (ctx == null || offscreen_ctx == null) {
        console.log("Error: Could not put canvas context to work")
        return
    }
    game_state = new GameState("CCDD", ball, player1, player2, canvas, offscreen, ctx, offscreen_ctx)


    raf = requestAnimationFrame(animate)
}

function draw_state(state: GameState): void {
    state.ctx.clearRect(0, 0, state.width, state.height);

    state.ctx.fillStyle = "blue"

    if (!swapped) {
        state.p1.draw_label(state.ctx, 10, 20)
        state.p2.draw_label(state.ctx, state.width - state.p2.get_label_width(state.ctx) - 10, 20)
    } else {
        state.p1.draw_label(state.ctx, state.width - state.p1.get_label_width(state.ctx) - 10, 20)
        state.p2.draw_label(state.ctx, 10, 20)
        
    }
    state.draw_code()

    state.p1.paddle.draw(state.ctx)
    state.p2.paddle.draw(state.ctx)

    state.ctx.fillStyle = "yellow"
    state.ctx.strokeStyle = "1px black"

    state.ball.draw(state.ctx)

}



function update(deltaTime: any) {
    // const speed = game_state.p1.paddle.speed
    const speed = 300

    if (keys.up) {
        game_state.p1.paddle.position.y -= speed * deltaTime
        // game_state.p1.paddle.position.y -= game_state.p1.paddle.speed
    }

    if (keys.down) {
        game_state.p1.paddle.position.y += speed * deltaTime
        // game_state.p1.paddle.position.y += game_state.p1.paddle.speed
    }

    game_state.p1.paddle.position.y = Math.max(0, Math.min(game_state.p1.paddle.position.y, canvas_height - game_state.p1.paddle.length));

    game_state.p1.paddle.position.y  = Math.round(game_state.p1.paddle.position.y)

}

window.addEventListener("keydown", (event) => {
    if(event.key == "w") { keys.up = true }
    if(event.key == "s") { keys.down = true }
})

window.addEventListener("keyup", (event) => {
    if(event.key == "w") { keys.up = false }
    if(event.key == "s") { keys.down = false }
})


let lastTime: number = 0;

function animate(time: any) {
    if (!time) {
        lastTime = time;
        window.requestAnimationFrame(animate)
        return
    }
    const deltaTime = (time - lastTime) / 1000;

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
        case EventType.CreatedRoom:
            handle_room_created(event)
            break;
        case EventType.JoinedRoom:
            handle_joined(event)
            break;
        case EventType.PlayerJoinedRoom:
            handle_player_joined(event)
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
        default:
            console.log("Unknown Event type: " + event.type)
            console.log(event)
            break
    }
}

function send_event(ev: SocketEvent): void {
    socket.send(JSON.stringify(ev))
}

function handle_other_player_joined(event: SocketEvent): void {
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
        game_state.p1.isConnected = true
        game_state.p2.isConnected = true
        game_state.p2.username = event.data.player        
        game_state.code = event.data.code
        swap_paddle_position()
    } else {
        console.log("no data", event)
    }
}

function handle_player_joined(event: SocketEvent): void {
    if (event.data) {
        console.log("player joined", event)
        game_state.p2.isConnected = true
        game_state.p2.username = event.data.player
    } else {
        console.log("no data", event)
    }
}

function handle_room_created(event: SocketEvent): void {
    game_state.code = event.data.code
    if (event.to) {
        game_state.p1.username = event.to
    }

    room_form.style.visibility = "hidden"
    var roomTitle = document.createElement("h1")
    roomTitle.classList.add("subtitle")
    roomTitle.classList.add("is-4")
    roomTitle.innerHTML = "Code: " + event.data.code    
    room_info_div.insertAdjacentElement("afterbegin", roomTitle)
    canvas.style.visibility = "visible"
}

function create_room(): void {
    const event = {
        type: EventType.CreateRoom,
    }
    send_event(event)    
}

function join_game(): void {
    console.log("Clicked joined")
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

main()

function swap_paddle_position() {
   const aux = game_state.p1.paddle
   game_state.p1.paddle = game_state.p2.paddle 
   game_state.p2.paddle = aux
   swapped = true
}

