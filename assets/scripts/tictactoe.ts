async function ttt_init() {

setTimeout(() => { 
    class TicPlayer {
        name: string  
        wins: number
        connected: boolean

        constructor(name: string, connected: boolean) {
          this.name = name
          this.connected = connected
          this.wins = 0
        }
    }

    class Board {
        board: number[][]

        constructor() {
            this.board = [
              [0, 0, 0],  
              [0, 0, 0],  
              [0, 0, 0]  
            ];
        }  
    }

    class State {
        player1: TicPlayer
        player2: TicPlayer  
        board: Board
        status: number
        ties: number

        constructor() {
           this.player1 = new TicPlayer("", false)   
           this.player2 = new TicPlayer("", false)   
           this.board = new Board()
           this.ties = 0
           this.status = 0
        }
    }

    enum TTTEventType {
        CreateGame = 1,
        JoinGame = 2,
        JoinedGame = 22,
        OtherPlayerJoined = 23,

        MakePlay = 3,

        PlayerDisconnected = 5,
        PlayerReconnected = 6,

        BoardCellUpdate = 8,
        StateUpdate = 9,

        Tie = 10,
        Victory = 11,
        Defeat = 12,
    }

    type TTTEvent = {
        type: TTTEventType,
        data?: any,
        isError?: boolean
    }
    let host = window.location.host
    let ttt_socket = new WebSocket("ws://"+host+"/ws/tictactoe")

    let ttt_create_btn = document.getElementById("tictactoe_create_btn") as HTMLButtonElement
    let ttt_join_btn = document.getElementById("tictactoe_join_btn") as HTMLButtonElement

    let ttt_code_label = document.getElementById("ttt_code_label") as HTMLParagraphElement

    let scoreboard = document.getElementById("scoreboard") as HTMLDivElement

    let player1_label = document.getElementById("player1_label") as HTMLParagraphElement
    let player1_wins = document.getElementById("player1_wins") as HTMLParagraphElement

    let player2_label = document.getElementById("player2_label") as HTMLParagraphElement
    let player2_wins = document.getElementById("player2_wins") as HTMLParagraphElement

    let ties_label = document.getElementById("ties") as HTMLParagraphElement
    ties_label.style.color = "green"

    let board_el = document.getElementById("ttt_board") as HTMLDivElement
    let cells = document.querySelectorAll('#ttt_board .cell');

    ttt_create_btn.addEventListener("click", ttt_create_game)
    ttt_join_btn.addEventListener("click", ttt_join_game)

    let state = new State()

    let clicked_create = false

    function ttt_create_game(): void {
        const create_event: TTTEvent = {
            type: TTTEventType.CreateGame,
        }
        ttt_send_event(create_event)
        clicked_create = true
    }

    function ttt_join_game(): void {
        const code_input = document.getElementById("ttt_code") as HTMLInputElement
        const join_event: TTTEvent = {
            type: TTTEventType.JoinGame,
            data: {
                code: code_input?.value            
            }
        }
        ttt_send_event(join_event)
        clicked_create = true
    }

    function ttt_send_event(ev: TTTEvent): void {
        ttt_socket.send(JSON.stringify(ev))
    }

    ttt_socket.addEventListener("message", (e) => {
        console.log(e)
        const event = ttt_parse_event(e.data)
        if (event.isError !== undefined) {
            ttt_handle_event_error(event)
            return
        }  

        if (event.data.type == "pong") {
            console.log("got pong")
            return
        }
        ttt_handle_event(event)
    })


    function ttt_handle_event(event: TTTEvent): void {
        switch (event.type) {
            case TTTEventType.JoinedGame:
                handle_ttt_joined_game(event)
                break
            case TTTEventType.OtherPlayerJoined:
                handle_ttt_other_player_joined_game(event)
                break
            case TTTEventType.StateUpdate:
                handle_ttt_state_update(event)
                break
            case TTTEventType.BoardCellUpdate:
                handle_ttt_board_cell_update(event)
                break
            case TTTEventType.Tie:
                handle_ttt_tie(event)
                break
            case TTTEventType.Victory:
                handle_ttt_victory(event)
                break
            case TTTEventType.Defeat:
                handle_ttt_defeat(event)
                break
            case TTTEventType.PlayerDisconnected:
                handle_ttt_player_disconnected(event)
                break
            case TTTEventType.PlayerReconnected:
                handle_ttt_player_reconnected(event)
                break
            default:
                console.log("Unknown event")
                console.log(event)
                break
            
        }
    }

    function ttt_handle_event_error(event: TTTEvent): void {
        console.log("Event gave an error: " + event.type)
        if (event.data) {
            console.log(event.data.message)
            alert(event.data.message)
        }
    }

    function ttt_parse_event(data: any): TTTEvent {
        const event_data = JSON.parse(data)
        const event: TTTEvent = event_data as TTTEvent

        if (!event.type) {
            throw new Error("Receive event without type")
        }
        return event
    }

    function handle_ttt_board_cell_update(event: TTTEvent): void {
        if(event.data) {
            const row: number = event.data.row
            const col: number = event.data.col
            const value: number = event.data.value

            update_cell(row, col, value)
        }
        if (state.status == 2) {
            clear_board()
            update_scoreboard()
            state.status = 1
        }
    }

    async function handle_ttt_tie(event: TTTEvent): Promise<void> {
        if (event.data) {
            state.ties = event.data.ties

            await update_cell(event.data.row, event.data.col, event.data.value)
            await show_game_result_message("tie", event.data.row, event.data.col).then(() => {
                clear_board()
            })
            clear_board()
        }
        update_scoreboard()
    }


    function close_web_socket() {
        if (ttt_socket) {
            ttt_socket.close();
        }
    }
    async function handle_ttt_victory(event: TTTEvent): Promise<void> {
        if (event.data) {
            state.player1.wins = event.data.player1.wins
            state.player1.name = event.data.player1.username

            state.player2.wins = event.data.player2.wins
            state.player2.name = event.data.player2.username

            await update_cell(event.data.row, event.data.col, event.data.value)
            await show_game_result_message("win", event.data.row, event.data.col).then(() => {
                clear_board()
            })
            clear_board()
        }
        update_scoreboard()
    }

    async function handle_ttt_defeat(event: TTTEvent): Promise<void> {
        if(event.data) {
            state.player1.wins = event.data.player1.wins
            state.player1.name = event.data.player1.username

            state.player2.wins = event.data.player2.wins
            state.player2.name = event.data.player2.username

            await update_cell(event.data.row, event.data.col, event.data.value)
            await show_game_result_message("lose", event.data.row, event.data.col).then(() => {
                clear_board()
            })
        }
        update_scoreboard()
    }

    function handle_ttt_player_reconnected(event: TTTEvent): void {
        if (event.data){
            let player: TicPlayer 
            if (event.data.username == state.player1.name) {
                player = state.player1
            } else if (event.data.username == state.player2.name) {
                player = state.player2
            } else {
                return
            }
            player.connected = true
            player.wins = event.data.wins
            update_scoreboard()

        }
    }

    function handle_ttt_player_disconnected(event: TTTEvent): void {
        if (event.data){
            let player: TicPlayer 
            if (event.data.username == state.player1.name) {
                player = state.player1
            } else if (event.data.username == state.player2.name) {
                player = state.player2
            } else {
                return
            }
            player.connected = false
            player.wins = event.data.wins
            update_scoreboard()
        }
    }

    function handle_ttt_state_update(event: TTTEvent): void {
        if (event.data) {
            state.player2.name = event.data.player
            state.board.board = event.data.board

            state.status = event.data.status
            state.ties = event.data.ties

            state.player1.wins = event.data.player1.wins
            state.player1.name = event.data.player1.username
            state.player1.connected = event.data.player1.connected

            state.player2.wins = event.data.player2.wins
            state.player2.name = event.data.player2.username
            state.player2.connected = event.data.player2.connected

        }
        update_scoreboard()
        update_board()
    }

    function handle_ttt_other_player_joined_game(event: TTTEvent): void {
        if (event.data) {
            state.player2.name = event.data.username
            state.player2.connected = event.data.connected
            state.player2.wins = event.data.wins
            update_scoreboard()
        }
    }

    function handle_ttt_joined_game(event: TTTEvent): void {
        const form = document.getElementById("room-menu") as HTMLDivElement
        form.style.display = "none"

        ttt_code_label.innerText = event.data.code
        if (clicked_create) {
            ttt_show_notification(event.data.code)
            navigator.clipboard.writeText(event.data.code);
        }

        if (event.data.player1) {
            state.player1.name = event.data.player1.username
            if (event.data.player1.connected == true) {
                state.player1.connected = true
            } else {
                state.player1.connected = false
            }
        }

        if (event.data.player2) {
            state.player2.name = event.data.player2.username

            if (event.data.player2.connected == true) {
                state.player2.connected = true
            } else {
                state.player2.connected = false
            }
        }
        state.ties = event.data.ties
        state.board.board = event.data.board
        state.status - event.data.status
    
        update_scoreboard()
        update_board()
    }

    function ttt_show_notification(message: string) {
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


        document.body.appendChild(notification);
        setTimeout(() => notification.style.opacity = "1", 10);

        setTimeout(() => {
            notification.style.opacity = "0";
            setTimeout(() => notification.remove(), 300); // 300ms to match the transition
        }, 2000);
    }

    const p = document.createElement("p") as HTMLParagraphElement
    p.innerText = "Waiting for 2nd player..."
    p.id = "waiting_label"

    function update_scoreboard():void {
        if (state.player1.name != "" && state.player2.name == "") {
            board_el.insertAdjacentElement("afterend", p)
            scoreboard.style.visibility = "hidden"
            return
        } else if (state.player1.name == "" && state.player2.name == "") {
            p.style.visibility = "none"
            p.innerText = ""
            scoreboard.style.visibility = "hidden"
        } else {
            scoreboard.style.visibility = "visible"
        }

        player1_label.innerHTML = state.player1.name.toUpperCase() + " (X)"
        player1_wins.innerHTML = state.player1.wins.toString()

        player2_label.innerHTML = state.player2.name.toUpperCase() + " (O)"
        player2_wins.innerHTML = state.player2.wins.toString()

        if (state.player1.connected == false) {
            player1_label.style.color = "red"
        } else {
            player1_label.style.color = "green"
        }
        if (state.player2.connected == false) {
            player2_label.style.color = "red"
        } else {
            player2_label.style.color = "green"
        }

        ties_label.innerHTML = state.ties.toString()
    
    }

    async function update_board(): Promise<void> {

        for (let row:number = 0; row < 3; row++) {
            for (let col:number = 0; col < 3; col++) {
                const cellIndex = row * 3 + col;
                const cell = cells[cellIndex];

                if (state.board.board[row][col] === 1) {
                    cell.classList.add('x');
                } else if (state.board.board[row][col] === 2) {
                    cell.classList.add('o');
                } else if (state.board.board[row][col] == 0) {
                    cell.classList.remove('o', 'x')
                }
            }
        }
    }

    function clear_board(): void {
        state.board.board = [ [0, 0, 0], [0,0,0], [0,0,0] ]
        update_board()
    }

    async function update_cell(row: number, col: number, value: number): Promise<void> {
        state.board.board[row][col] = value
        const cell = cells[row * 3 + col]

        if (value == 1) {
            cell.classList.add('x');
        } else if (value == 2) {
            cell.classList.add('o');
        } else if (value == 0) {
            cell.classList.remove('o', 'x')
        }
    
    }

    update_scoreboard()

    cells.forEach((cell, index) => {
        cell.addEventListener('click', () => {
            const row: number = Math.floor(index / 3);
            const col: number  = index % 3;

        
            const play_event: TTTEvent = {
                type: TTTEventType.MakePlay,
                data: {
                    row: row,
                    col: col,
                }
            }
            ttt_send_event(play_event)
        });
    });

    async function show_game_result_message(result: string, winning_row: number, winning_col: number): Promise<void> {
        return new Promise((resolve) => {
            const messageElement = document.getElementById("game-result-overlay") as HTMLDivElement;
            const resultTextElement = document.getElementById("result-text") as HTMLParagraphElement;
            const winningCells = get_winning_cells(winning_row, winning_col); // Retorna um array de coordenadas [row, col]

            if (result === "win") {
                resultTextElement.textContent = "You Win!";
            } else if (result === "lose") {
                resultTextElement.textContent = "You Lose!";
            } else if (result === "tie") {
                resultTextElement.textContent = "It's a Tie!";
            }

            messageElement.classList.remove("is-hidden");

            if (winningCells && result !== "tie") {
                winningCells.forEach(([row, col]) => {
                    const cellIndex = row * 3 + col; // Calcula o índice da célula com base em row e col
                    const cell = document.querySelectorAll('.cell')[cellIndex] as HTMLElement;
                    cell.classList.add('winning-cell');
                });
            }

            setTimeout(() => {
                messageElement.classList.add("is-hidden");
                if (winningCells) {
                    winningCells.forEach(([row, col]) => {
                        const cellIndex = row * 3 + col;
                        const cell = document.querySelectorAll('.cell')[cellIndex] as HTMLElement;
                        cell.classList.remove('winning-cell');
                    });
                }
                resolve()
            }, 1500);
        })
    }

    function get_winning_cells(lastRow: number, lastCol: number): number[][] {
        const player = state.board.board[lastRow][lastCol];
        const winningCells = [];
        const board = state.board.board

        // Check row
        if (board[lastRow][0] === player && board[lastRow][1] === player && board[lastRow][2] === player) {
            winningCells.push([lastRow, 0], [lastRow, 1], [lastRow, 2]);
        }

        // Check column
        if (board[0][lastCol] === player && board[1][lastCol] === player && board[2][lastCol] === player) {
            winningCells.push([0, lastCol], [1, lastCol], [2, lastCol]);
        }

        // Check diagonal (top-left to bottom-right)
        if (lastRow === lastCol && board[0][0] === player && board[1][1] === player && board[2][2] === player) {
            winningCells.push([0, 0], [1, 1], [2, 2]);
        }

        // Check diagonal (top-right to bottom-left)
        if (lastRow + lastCol === 2 && board[0][2] === player && board[1][1] === player && board[2][0] === player) {
            winningCells.push([0, 2], [1, 1], [2, 0]);
        }

        // If winning cells were found, return them
        return winningCells
    }


        document.body.addEventListener('htmx:afterOnLoad', function(event) {
            ttt_socket.close()
        });

    }, 200);
}
