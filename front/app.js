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
        const res = await fetch(`${API_URL}/game/${currentGame.id}`, {
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
    
    const gameIdDisplay = document.getElementById("game-id-display");
    if (gameIdDisplay) gameIdDisplay.innerText = currentGame.id;

    const boardDiv = document.getElementById("board");
    if (!boardDiv) {
        console.error("Контейнер #board не найден в HTML!");
        return;
    }

    // 1. ОЧИСТКА: Удаляем только ячейки, НЕ трогая линию зачеркивания
    const cells = boardDiv.querySelectorAll('.cell');
    cells.forEach(cell => cell.remove());

    const statusMsg = document.getElementById("status-message");
    if (statusMsg) statusMsg.innerText = "";

    // 2. ОТРИСОВКА ЯЧЕЕК
    if (currentGame.board && currentGame.board.length > 0) {
        currentGame.board.forEach((row, rIdx) => {
            row.forEach((cell, cIdx) => {
                const cellDiv = document.createElement("div");
                cellDiv.className = `cell ${cell !== 0 ? 'taken' : ''}`;
                
                // Отображаем X (1) или O (2)
                if (cell === 1) cellDiv.innerText = "X";
                if (cell === 2) cellDiv.innerText = "O";

                // Обработка клика: только если клетка пуста (0)
                if (cell === 0) {
                    cellDiv.onclick = () => makeMove(rIdx, cIdx);
                }
                
                boardDiv.appendChild(cellDiv);
            });
        });
    }

    // 3. ЛОГИКА ЗАВЕРШЕНИЯ (Линия и статус)
    const winner = checkAndDrawWinningLine(currentGame.board);
    if (winner) {
        if (statusMsg) {
            statusMsg.innerText = (winner === 1 ? "Вы победили!" : "Победил ИИ!");
            statusMsg.style.color = (winner === 1 ? "green" : "red");
        }
        return; 
    }

    // Проверка на ничью
    const isDraw = !currentGame.board.flat().includes(0);
    if (isDraw && statusMsg) {
        statusMsg.innerText = "Ничья!";
        statusMsg.style.color = "blue";
    }
}

function showError(msg) { document.getElementById("error-message").innerText = msg; }

const WINNING_COMBOS = [
    // Строки (type: 'h', index: 0,1,2)
    { cells: [0, 1, 2], type: 'h', lineIdx: 0 },
    { cells: [3, 4, 5], type: 'h', lineIdx: 1 },
    { cells: [6, 7, 8], type: 'h', lineIdx: 2 },
    // Столбцы (type: 'v', index: 0,1,2)
    { cells: [0, 3, 6], type: 'v', lineIdx: 0 },
    { cells: [1, 4, 7], type: 'v', lineIdx: 1 },
    { cells: [2, 5, 8], type: 'v', lineIdx: 2 },
    // Диагонали (type: 'd', index: 0-главная, 1-побочная)
    { cells: [0, 4, 8], type: 'd', lineIdx: 0 },
    { cells: [2, 4, 6], type: 'd', lineIdx: 1 }
];

// Функция проверки победителя и отрисовки линии
function checkAndDrawWinningLine(board) {
    const flatBoard = board.flat(); // Превращаем [3][3] в [9]
    const lineEl = document.getElementById("winning-line");
    
    // Сбрасываем старую линию
    lineEl.style.display = 'none';
    lineEl.className = 'winning-line'; 
    lineEl.style.width = '';
    lineEl.style.height = '';
    lineEl.style.transform = '';

    for (let combo of WINNING_COMBOS) {
        const [a, b, c] = combo.cells;
        
        // Если все три клетки заняты одним игроком (1 или 2) и не пустые (0)
        if (flatBoard[a] && flatBoard[a] === flatBoard[b] && flatBoard[a] === flatBoard[c]) {
            drawSVGLine(combo); // Рисуем линию
            return flatBoard[a]; // Возвращаем ID победителя (1 или 2)
        }
    }
    return null; // Победителя нет
}

// Функция позиционирования линии
function drawSVGLine(combo) {
    const lineEl = document.getElementById("winning-line");
    if (!lineEl) return;

    const cellSize = 100;
    const gap = 5;
    const boardSize = 310;
    const center = cellSize / 2; // 50px

    // 1. Полный сброс всех стилей перед отрисовкой
    lineEl.style.display = "block";
    lineEl.style.width = "0";
    lineEl.style.height = "0";
    lineEl.style.top = "auto";
    lineEl.style.left = "auto";
    lineEl.style.right = "auto";
    lineEl.style.transform = "none";
    lineEl.style.transformOrigin = "center";

    if (combo.type === 'h') {
        // ГОРИЗОНТАЛЬ
        const y = combo.lineIdx * (cellSize + gap) + center;
        lineEl.style.top = `${y - 3}px`; // -3 для центровки толщины (6px)
        lineEl.style.left = "5px";
        lineEl.style.height = "6px";
        setTimeout(() => lineEl.style.width = "300px", 10);

    } else if (combo.type === 'v') {
        // ВЕРТИКАЛЬ
        const x = combo.lineIdx * (cellSize + gap) + center;
        lineEl.style.left = `${x - 3}px`;
        lineEl.style.top = "5px";
        lineEl.style.width = "6px";
        setTimeout(() => lineEl.style.height = "300px", 10);

    } else if (combo.type === 'd') {
        // ДИАГОНАЛИ
        lineEl.style.width = "6px";
        lineEl.style.transformOrigin = "50% 0"; // Вращаем от центра верхней точки
        const diagLength = Math.sqrt(2) * boardSize - 20; // Чуть короче полной диагонали

        if (combo.lineIdx === 0) {
            // Главная (0,0 -> 2,2)
            lineEl.style.top = "0";
            lineEl.style.left = "0";
            lineEl.style.transform = "rotate(-45deg)";
        } else {
            // Побочная (0,2 -> 2,0)
            lineEl.style.top = "0";
            lineEl.style.left = `${boardSize}px`;
            lineEl.style.transform = "rotate(45deg)";
        }
        setTimeout(() => lineEl.style.height = `${diagLength}px`, 10);
    }
}