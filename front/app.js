const API_URL = "http://88.218.67.103:8080";
let currentGame = { id: "", board: [] };

// Показать/Скрыть экраны
function showScreen(id) {
    ["menu", "load-input", "game-screen"].forEach(s => 
        document.getElementById(s).classList.add("hidden"));
    document.getElementById(id).classList.remove("hidden");
    document.getElementById("error-message").innerText = "";
}

function showMenu() { showScreen("menu"); }
function showLoadInput() { showScreen("load-input"); }

// 1. Создание новой игры
async function createNewGame() {
    try {
        const res = await fetch(`${API_URL}/game`, { method: 'POST' });
        currentGame = await res.json();
        renderGame();
    } catch (e) { showError("Ошибка сервера"); }
}

// 2. Загрузка по UUID
async function loadGame() {
    const uuid = document.getElementById("uuid-field").value;
    try {
        const res = await fetch(`${API_URL}/game/${uuid}`);
        if (!res.ok) throw new Error("Игра не найдена");
        currentGame = await res.json();
        renderGame();
    } catch (e) { showError(e.message); }
}

// 3. Обработка хода
async function makeMove(row, col) {
    if (currentGame.board[row][col] !== 0) return;

    // Клонируем доску и ставим наш ход (1 = X)
    let newBoard = JSON.parse(JSON.stringify(currentGame.board));
    newBoard[row][col] = 1;

    try {
        const res = await fetch(`${API_URL}/game/${currentGame.id}/move`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ board: newBoard })
        });
        
        if (!res.ok) {
            const err = await res.text();
            throw new Error(err);
        }
        
        currentGame = await res.json();
        renderGame();
    } catch (e) { showError(e.message); }
}

function renderGame() {
    showScreen("game-screen");
    document.getElementById("game-id-display").innerText = currentGame.id;
    const boardDiv = document.getElementById("board");
    boardDiv.innerHTML = "";

    currentGame.board.forEach((row, rIdx) => {
        row.forEach((cell, cIdx) => {
            const cellDiv = document.createElement("div");
            cellDiv.className = `cell ${cell !== 0 ? 'taken' : ''}`;
            cellDiv.innerText = cell === 1 ? "X" : cell === 2 ? "O" : "";
            cellDiv.onclick = () => makeMove(rIdx, cIdx);
            boardDiv.appendChild(cellDiv);
        });
    });
}

function showError(msg) { document.getElementById("error-message").innerText = msg; }