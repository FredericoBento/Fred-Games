type Point = {
    x: number,
    y: number,
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

export class Player {
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
        ctx.fillStyle = "white"
        ctx.fillText(this.username, x, y);
    }

    get_label_width(ctx: CanvasRenderingContext2D): number {
        let old_font = ctx.font
        ctx.font = this.label_font
        ctx.font = old_font
        return ctx.measureText(this.username).width
    }
}

export class Paddle {
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
        ctx.fillStyle = "blue";
        ctx.fillRect(
            this.position.x,
            this.position.y,
            this.width,
            this.length,
        );
    }
}


export class Ball {
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
        ctx.fillStyle = "yellow";
        ctx.arc(
            this.position.x,
            this.position.y,
            this.radius,
            0,
            2 * Math.PI
        );
        ctx.fill();
    }
}

export class GameState {
    code: string;
    ball: Ball;
    p1: Player;
    p2: Player;
    status: GameStatus;

    constructor(
        code: string,
        ball: Ball,
        p1: Player,
        p2: Player,
    ){
        this.code = code;
        this.ball = ball;
        this.p1 = p1;
        this.p2 = p2;
        this.status = GameStatus.NotStarted;
    }

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

const canvas_width = 640
const canvas_height = 360

// const canvas_div = document.getElementById("canvasDiv") as HTMLDivElement;
const canvas = document.getElementById("gameCanvas") as HTMLCanvasElement;
const ctx = canvas.getContext("2d")

const socket = new WebSocket("ws://localhost:8080/ws/pong");

const room_form = document.getElementById("room-menu") as HTMLDivElement;
const room_info_div = document.getElementById("roomInfo") as HTMLDivElement;

const join_btn = document.getElementById("joinBtn") as HTMLButtonElement;
const create_btn = document.getElementById("createBtn") as HTMLButtonElement;

const code_input = document.getElementById("code") as HTMLInputElement;
join_btn?.addEventListener("click", join_game);

let keys = {
    up: false,
    down: false,
}

let game_state: GameState;

let lastTimestamp: any;

function main(): void{
    // Get Settings from server
    // Set Ball
    // Set Player1 Object
    // then create layer
    const paddle: Paddle = new Paddle({x: 30, y: canvas_height/2 }, 40, 4, 200)
    const paddle2: Paddle = new Paddle({x: canvas_width-30, y: canvas_height/2 }, 40, 4, 200)

    const player1: Player = new Player("My Name", paddle, true)
    const player2: Player = new Player("To be connected...", paddle2, false)

    const ball: Ball = new Ball({x: canvas_width/2, y: canvas_height/2}, 7, 6, Direction.Left); 

    game_state = new GameState("CCDD", ball, player1, player2)


    lastTimestamp = 0
    window.requestAnimationFrame(animate)

}

function draw_state(state: GameState): boolean{
    if (ctx == null) {
        return false
    }

    // window.requestAnimationFrame(update);

    ctx.clearRect(0, 0, canvas_width, canvas_height);
    ctx.fillStyle = "black"
    ctx.fillRect(0, 0, canvas_width, canvas_height);
    ctx.save();

    state.ball.draw(ctx)
    state.p1.paddle.draw(ctx)
    state.p2.paddle.draw(ctx)

    state.p1.draw_label(ctx, 5, 25)
    state.p2.draw_label(ctx, canvas_width - state.p2.get_label_width(ctx), 20)

    return true
}

function update(deltaTime: any) {

    if (keys.up) {
        game_state.p1.paddle.position.y -= game_state.p1.paddle.speed * deltaTime
    }

    if (keys.down) {
        game_state.p1.paddle.position.y += game_state.p1.paddle.speed * deltaTime
    }

        game_state.p1.paddle.position.y = Math.max(0, Math.min(game_state.p1.paddle.position.y, canvas_height - game_state.p1.paddle.length));
}

window.addEventListener("keydown", (event) => {
    if(event.key == "w") { keys.up = true }
    if(event.key == "s") { keys.down = true }
})

window.addEventListener("keyup", (event) => {
    if(event.key == "w") { keys.up = false }
    if(event.key == "s") { keys.down = false }
})

function animate(timestamp: any) {
    if (!lastTimestamp) {
        lastTimestamp = timestamp;
        window.requestAnimationFrame(animate);
        return;
    }

    let deltaTime: number = timestamp - lastTimestamp;
    deltaTime = deltaTime / 1000;
    
    update(deltaTime);
    draw_state(game_state);

    lastTimestamp = timestamp;
    window.requestAnimationFrame(animate);
}


function handle_event(): void {
}

function send_event(): void {
}


function join_game(): void {
    console.log("Clicked joined")
    let code: string = code_input?.value
    if (code == "") {
        console.log("abc");
    }
}

main()
