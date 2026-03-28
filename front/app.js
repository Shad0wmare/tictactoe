const API_URL = "http://88.218.67.103:8080"; 
let currentGame = { id: "", board: [] };

function showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(s => s.classList.add('hidden'));
    
    const target = document.getElementById(screenId);
    if (target) target.classList.remove('hidden');

    const errorEl = document.getElementById("error-message");
    if (errorEl) errorEl.innerText = "";

    if (screenId === 'screen-main') {
        const line = document.getElementById("winning-line");
        if (line) line.style.display = "none";
    }
}


async function createNewGame() {
    try {
        const res = await fetch(`${API_URL}/game`, { method: 'POST' });
        if (!res.ok) throw new Error("Ошибка при создании игры");
        currentGame = await res.json();
        renderGame();
    } catch (e) {
        showError("Сервер недоступен: " + e.message);
    }
}

async function loadGame() {
    const uuid = document.getElementById("uuid-field").value.trim();
    if (!uuid) return showError("Введите UUID игры");

    try {
        const res = await fetch(`${API_URL}/game/${uuid}`);
        if (!res.ok) throw new Error("Игра с таким ID не найдена");
        currentGame = await res.json();
        renderGame();
    } catch (e) {
        showError(e.message);
    }
}

async function makeMove(row, col) {
    if (currentGame.board[row][col] !== 0) return;

    let newBoard = JSON.parse(JSON.stringify(currentGame.board));
    newBoard[row][col] = 1;

    try {
        const res = await fetch(`${API_URL}/game/${currentGame.id}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ board: newBoard })
        });
        
        if (!res.ok) {
            const errMsg = await res.text();
            throw new Error(errMsg);
        }
        
        currentGame = await res.json();
        renderGame();
    } catch (e) {
        showError(e.message);
    }
}

function renderGame() {
    showScreen('screen-game');
    
    document.getElementById("game-id-display").innerText = currentGame.id;
    const boardDiv = document.getElementById("board");
    const statusMsg = document.getElementById("status-message");
    
    const oldCells = boardDiv.querySelectorAll('.cell');
    oldCells.forEach(c => c.remove());
    statusMsg.innerText = "";

    currentGame.board.forEach((row, rIdx) => {
        row.forEach((cell, cIdx) => {
            const cellDiv = document.createElement("div");
            cellDiv.className = `cell ${cell !== 0 ? 'taken' : ''}`;
            cellDiv.innerText = cell === 1 ? "X" : cell === 2 ? "O" : "";
            
            if (cell === 0) {
                cellDiv.onclick = () => makeMove(rIdx, cIdx);
            }
            boardDiv.appendChild(cellDiv);
        });
    });

    const winner = checkAndDrawWinningLine(currentGame.board);
    if (winner) {
        statusMsg.innerText = (winner === 1 ? "Вы победили!" : "Победил ИИ!");
        statusMsg.style.color = (winner === 1 ? "#28a745" : "#dc3545");
    } else if (!currentGame.board.flat().includes(0)) {
        statusMsg.innerText = "Ничья!";
        statusMsg.style.color = "#007bff";
    }
}

const WINNING_COMBOS = [
    { cells: [0, 1, 2], type: 'h', lineIdx: 0 },
    { cells: [3, 4, 5], type: 'h', lineIdx: 1 },
    { cells: [6, 7, 8], type: 'h', lineIdx: 2 },
    { cells: [0, 3, 6], type: 'v', lineIdx: 0 },
    { cells: [1, 4, 7], type: 'v', lineIdx: 1 },
    { cells: [2, 5, 8], type: 'v', lineIdx: 2 },
    { cells: [0, 4, 8], type: 'd', lineIdx: 0 },
    { cells: [2, 4, 6], type: 'd', lineIdx: 1 }
];

function checkAndDrawWinningLine(board) {
    const flatBoard = board.flat();
    const lineEl = document.getElementById("winning-line");
    lineEl.style.display = "none";

    for (let combo of WINNING_COMBOS) {
        const [a, b, c] = combo.cells;
        if (flatBoard[a] && flatBoard[a] === flatBoard[b] && flatBoard[a] === flatBoard[c]) {
            drawVisualLine(combo);
            return flatBoard[a];
        }
    }
    return null;
}

function drawVisualLine(combo) {
    const lineEl = document.getElementById("winning-line");
    const boardSize = 310;
    const cellStep = 105;
    const center = 50;

    lineEl.style.display = "block";
    lineEl.className = "winning-line";
    lineEl.style.width = "0";
    lineEl.style.height = "0";
    lineEl.style.transform = "none";

    if (combo.type === 'h') {
        lineEl.style.top = `${combo.lineIdx * cellStep + center - 3}px`;
        lineEl.style.left = "5px";
        lineEl.style.height = "6px";
        setTimeout(() => lineEl.style.width = "300px", 10);
    } else if (combo.type === 'v') {
        lineEl.style.left = `${combo.lineIdx * cellStep + center - 3}px`;
        lineEl.style.top = "5px";
        lineEl.style.width = "6px";
        setTimeout(() => lineEl.style.height = "300px", 10);
    } else if (combo.type === 'd') {
        lineEl.style.width = "6px";
        lineEl.style.top = "0";
        lineEl.style.transformOrigin = "50% 0";
        const diagLen = Math.sqrt(2) * boardSize - 20;
        if (combo.lineIdx === 0) {
            lineEl.style.left = "0";
            lineEl.style.transform = "rotate(-45deg)";
        } else {
            lineEl.style.left = `${boardSize}px`;
            lineEl.style.transform = "rotate(45deg)";
        }
        setTimeout(() => lineEl.style.height = `${diagLen}px`, 10);
    }
}

function showError(msg) {
    const errEl = document.getElementById("error-message");
    if (errEl) errEl.innerText = msg;
}