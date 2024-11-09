"use strict";
class TicPlayer {
    name;
    wins;
    connected;
    constructor(name, connected) {
        this.name = name;
        this.connected = connected;
        this.wins = 0;
    }
}
class Board {
    board;
    constructor() {
        this.board = [
            [0, 0, 0],
            [0, 0, 0],
            [0, 0, 0]
        ];
    }
}
class State {
    player1;
    player2;
    board;
    status;
    ties;
    constructor() {
        this.player1 = new TicPlayer("", false);
        this.player2 = new TicPlayer("", false);
        this.board = new Board();
        this.ties = 0;
        this.status = 0;
    }
}
var TTTEventType;
(function (TTTEventType) {
    TTTEventType[TTTEventType["CreateGame"] = 1] = "CreateGame";
    TTTEventType[TTTEventType["JoinGame"] = 2] = "JoinGame";
    TTTEventType[TTTEventType["JoinedGame"] = 22] = "JoinedGame";
    TTTEventType[TTTEventType["OtherPlayerJoined"] = 23] = "OtherPlayerJoined";
    TTTEventType[TTTEventType["MakePlay"] = 3] = "MakePlay";
    TTTEventType[TTTEventType["PlayerDisconnected"] = 5] = "PlayerDisconnected";
    TTTEventType[TTTEventType["PlayerReconnected"] = 6] = "PlayerReconnected";
    TTTEventType[TTTEventType["BoardCellUpdate"] = 8] = "BoardCellUpdate";
    TTTEventType[TTTEventType["StateUpdate"] = 9] = "StateUpdate";
    TTTEventType[TTTEventType["Tie"] = 10] = "Tie";
    TTTEventType[TTTEventType["Victory"] = 11] = "Victory";
    TTTEventType[TTTEventType["Defeat"] = 12] = "Defeat";
})(TTTEventType || (TTTEventType = {}));
document.addEventListener("DOMContentLoaded", () => {
    const ttt_socket = new WebSocket("ws://localhost:8080/ws/tictactoe");
    const ttt_create_btn = document.getElementById("tictactoe_create_btn");
    const ttt_join_btn = document.getElementById("tictactoe_join_btn");
    const ttt_code_label = document.getElementById("ttt_code_label");
    const scoreboard = document.getElementById("scoreboard");
    const player1_label = document.getElementById("player1_label");
    const player1_wins = document.getElementById("player1_wins");
    const player2_label = document.getElementById("player2_label");
    const player2_wins = document.getElementById("player2_wins");
    const ties_label = document.getElementById("ties");
    ties_label.style.color = "green";
    const board_el = document.getElementById("ttt_board");
    const cells = document.querySelectorAll('#ttt_board .cell');
    ttt_create_btn.addEventListener("click", ttt_create_game);
    ttt_join_btn.addEventListener("click", ttt_join_game);
    const state = new State();
    let clicked_create = false;
    function ttt_create_game() {
        const create_event = {
            type: TTTEventType.CreateGame,
        };
        ttt_send_event(create_event);
        clicked_create = true;
    }
    function ttt_join_game() {
        const code_input = document.getElementById("ttt_code");
        const join_event = {
            type: TTTEventType.JoinGame,
            data: {
                code: code_input?.value
            }
        };
        ttt_send_event(join_event);
        clicked_create = true;
    }
    function ttt_send_event(ev) {
        ttt_socket.send(JSON.stringify(ev));
    }
    ttt_socket.addEventListener("message", (e) => {
        const event = ttt_parse_event(e.data);
        if (event.isError !== undefined) {
            ttt_handle_event_error(event);
            return;
        }
        ttt_handle_event(event);
    });
    function ttt_handle_event(event) {
        switch (event.type) {
            case TTTEventType.JoinedGame:
                handle_ttt_joined_game(event);
                break;
            case TTTEventType.OtherPlayerJoined:
                handle_ttt_other_player_joined_game(event);
                break;
            case TTTEventType.StateUpdate:
                handle_ttt_state_update(event);
                break;
            case TTTEventType.BoardCellUpdate:
                handle_ttt_board_cell_update(event);
                break;
            case TTTEventType.Tie:
                handle_ttt_tie(event);
                break;
            case TTTEventType.Victory:
                handle_ttt_victory(event);
                break;
            case TTTEventType.Defeat:
                handle_ttt_defeat(event);
                break;
            case TTTEventType.PlayerDisconnected:
                handle_ttt_player_disconnected(event);
                break;
            case TTTEventType.PlayerReconnected:
                handle_ttt_player_reconnected(event);
                break;
            default:
                console.log("Unknown event");
                console.log(event);
                break;
        }
    }
    function ttt_handle_event_error(event) {
        console.log("Event gave an error: " + event.type);
        if (event.data) {
            console.log(event.data.message);
            alert(event.data.message);
        }
    }
    function ttt_parse_event(data) {
        const event_data = JSON.parse(data);
        const event = event_data;
        if (!event.type) {
            throw new Error("Receive event without type");
        }
        return event;
    }
    function handle_ttt_board_cell_update(event) {
        if (event.data) {
            const row = event.data.row;
            const col = event.data.col;
            const value = event.data.value;
            update_cell(row, col, value);
        }
        if (state.status == 2) {
            clear_board();
            update_scoreboard();
            state.status = 1;
        }
    }
    async function handle_ttt_tie(event) {
        if (event.data) {
            state.ties = event.data.ties;
            await update_cell(event.data.row, event.data.col, event.data.value);
            await show_game_result_message("tie", event.data.row, event.data.col).then(() => {
                clear_board();
            });
            clear_board();
        }
        update_scoreboard();
    }
    async function handle_ttt_victory(event) {
        if (event.data) {
            state.player1.wins = event.data.player1.wins;
            state.player1.name = event.data.player1.username;
            state.player2.wins = event.data.player2.wins;
            state.player2.name = event.data.player2.username;
            await update_cell(event.data.row, event.data.col, event.data.value);
            await show_game_result_message("win", event.data.row, event.data.col).then(() => {
                clear_board();
            });
            clear_board();
        }
        update_scoreboard();
    }
    async function handle_ttt_defeat(event) {
        if (event.data) {
            state.player1.wins = event.data.player1.wins;
            state.player1.name = event.data.player1.username;
            state.player2.wins = event.data.player2.wins;
            state.player2.name = event.data.player2.username;
            await update_cell(event.data.row, event.data.col, event.data.value);
            await show_game_result_message("lose", event.data.row, event.data.col).then(() => {
                clear_board();
            });
        }
        update_scoreboard();
    }
    function handle_ttt_player_reconnected(event) {
        if (event.data) {
            let player;
            if (event.data.username == state.player1.name) {
                player = state.player1;
            }
            else if (event.data.username == state.player2.name) {
                player = state.player2;
            }
            else {
                return;
            }
            player.connected = true;
            player.wins = event.data.wins;
            update_scoreboard();
        }
    }
    function handle_ttt_player_disconnected(event) {
        if (event.data) {
            let player;
            if (event.data.username == state.player1.name) {
                player = state.player1;
            }
            else if (event.data.username == state.player2.name) {
                player = state.player2;
            }
            else {
                return;
            }
            player.connected = false;
            player.wins = event.data.wins;
            update_scoreboard();
        }
    }
    function handle_ttt_state_update(event) {
        if (event.data) {
            state.player2.name = event.data.player;
            state.board.board = event.data.board;
            state.status = event.data.status;
            state.ties = event.data.ties;
            state.player1.wins = event.data.player1.wins;
            state.player1.name = event.data.player1.username;
            state.player1.connected = event.data.player1.connected;
            state.player2.wins = event.data.player2.wins;
            state.player2.name = event.data.player2.username;
            state.player2.connected = event.data.player2.connected;
        }
        update_scoreboard();
        update_board();
    }
    function handle_ttt_other_player_joined_game(event) {
        if (event.data) {
            state.player2.name = event.data.username;
            state.player2.connected = event.data.connected;
            state.player2.wins = event.data.wins;
            update_scoreboard();
        }
    }
    function handle_ttt_joined_game(event) {
        const form = document.getElementById("room-menu");
        form.style.display = "none";
        ttt_code_label.innerText = event.data.code;
        if (clicked_create) {
            ttt_show_notification(event.data.code);
            navigator.clipboard.writeText(event.data.code);
        }
        if (event.data.player1) {
            state.player1.name = event.data.player1.username;
            if (event.data.player1.connected == true) {
                state.player1.connected = true;
            }
            else {
                state.player1.connected = false;
            }
        }
        if (event.data.player2) {
            state.player2.name = event.data.player2.username;
            if (event.data.player2.connected == true) {
                state.player2.connected = true;
            }
            else {
                state.player2.connected = false;
            }
        }
        state.ties = event.data.ties;
        state.board.board = event.data.board;
        state.status - event.data.status;
        update_scoreboard();
        update_board();
    }
    function ttt_show_notification(message) {
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
    const p = document.createElement("p");
    p.innerText = "Waiting for 2nd player...";
    p.id = "waiting_label";
    function update_scoreboard() {
        if (state.player1.name != "" && state.player2.name == "") {
            board_el.insertAdjacentElement("afterend", p);
            scoreboard.style.visibility = "hidden";
            return;
        }
        else if (state.player1.name == "" && state.player2.name == "") {
            p.style.visibility = "none";
            p.innerText = "";
            scoreboard.style.visibility = "hidden";
        }
        else {
            scoreboard.style.visibility = "visible";
        }
        player1_label.innerHTML = state.player1.name.toUpperCase() + " (X)";
        player1_wins.innerHTML = state.player1.wins.toString();
        player2_label.innerHTML = state.player2.name.toUpperCase() + " (O)";
        player2_wins.innerHTML = state.player2.wins.toString();
        if (state.player1.connected == false) {
            player1_label.style.color = "red";
        }
        else {
            player1_label.style.color = "green";
        }
        if (state.player2.connected == false) {
            player2_label.style.color = "red";
        }
        else {
            player2_label.style.color = "green";
        }
        ties_label.innerHTML = state.ties.toString();
    }
    async function update_board() {
        for (let row = 0; row < 3; row++) {
            for (let col = 0; col < 3; col++) {
                const cellIndex = row * 3 + col;
                const cell = cells[cellIndex];
                if (state.board.board[row][col] === 1) {
                    cell.classList.add('x');
                }
                else if (state.board.board[row][col] === 2) {
                    cell.classList.add('o');
                }
                else if (state.board.board[row][col] == 0) {
                    cell.classList.remove('o', 'x');
                }
            }
        }
    }
    function clear_board() {
        state.board.board = [[0, 0, 0], [0, 0, 0], [0, 0, 0]];
        update_board();
    }
    async function update_cell(row, col, value) {
        state.board.board[row][col] = value;
        const cell = cells[row * 3 + col];
        if (value == 1) {
            cell.classList.add('x');
        }
        else if (value == 2) {
            cell.classList.add('o');
        }
        else if (value == 0) {
            cell.classList.remove('o', 'x');
        }
    }
    update_scoreboard();
    cells.forEach((cell, index) => {
        cell.addEventListener('click', () => {
            const row = Math.floor(index / 3);
            const col = index % 3;
            const play_event = {
                type: TTTEventType.MakePlay,
                data: {
                    row: row,
                    col: col,
                }
            };
            ttt_send_event(play_event);
        });
    });
    async function show_game_result_message(result, winning_row, winning_col) {
        return new Promise((resolve) => {
            const messageElement = document.getElementById("game-result-overlay");
            const resultTextElement = document.getElementById("result-text");
            const winningCells = getWinningCells(winning_row, winning_col); // Retorna um array de coordenadas [row, col]
            // Define a mensagem com base no resultado
            if (result === "win") {
                resultTextElement.textContent = "You Win!";
            }
            else if (result === "lose") {
                resultTextElement.textContent = "You Lose!";
            }
            else if (result === "tie") {
                resultTextElement.textContent = "It's a Tie!";
            }
            // Exibe a mensagem sobre o resultado
            messageElement.classList.remove("is-hidden");
            // Destaca as células vencedoras, se houver um vencedor
            if (winningCells && result !== "tie") {
                winningCells.forEach(([row, col]) => {
                    const cellIndex = row * 3 + col; // Calcula o índice da célula com base em row e col
                    const cell = document.querySelectorAll('.cell')[cellIndex];
                    cell.classList.add('winning-cell');
                });
            }
            // Opcional: oculta a mensagem e redefine o tabuleiro após alguns segundos
            setTimeout(() => {
                messageElement.classList.add("is-hidden");
                if (winningCells) {
                    winningCells.forEach(([row, col]) => {
                        const cellIndex = row * 3 + col;
                        const cell = document.querySelectorAll('.cell')[cellIndex];
                        cell.classList.remove('winning-cell');
                    });
                }
                resolve();
            }, 1500); // Ajuste o tempo conforme necessário
        });
    }
    function getWinningCells(lastRow, lastCol) {
        const player = state.board.board[lastRow][lastCol];
        const winningCells = [];
        const board = state.board.board;
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
        return winningCells;
    }
});
//# sourceMappingURL=tictactoe.js.map