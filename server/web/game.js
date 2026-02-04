// Game State
let gameState = null;

// API Functions
async function api(endpoint, method = 'GET', data = null) {
    const options = {
        method,
        headers: { 'Content-Type': 'application/json' },
    };
    if (data) {
        options.body = JSON.stringify(data);
    }

    const response = await fetch(`/api/game${endpoint}`, options);
    if (!response.ok) {
        const error = await response.text();
        throw new Error(error);
    }
    return response.json();
}

async function createGame(playerNames) {
    try {
        gameState = await api('/new', 'POST', { playerNames });
        showScreen('game-screen');
        renderGame();
    } catch (e) {
        alert('Failed to create game: ' + e.message);
    }
}

async function refreshState() {
    try {
        gameState = await api('/state');
        renderGame();
    } catch (e) {
        console.error('Failed to refresh state:', e);
    }
}

async function executeAction(type, playerId, data = {}) {
    try {
        gameState = await api('/action', 'POST', { type, playerId, data });
        renderGame();

        if (gameState.actionResult && !gameState.actionResult.success) {
            alert(gameState.actionResult.error || 'Action failed');
        }
    } catch (e) {
        alert(e.message);
    }
}

// Screen Management
function showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
    document.getElementById(screenId).classList.add('active');
}

// Setup Screen
function initSetup() {
    const addBtn = document.getElementById('add-player-btn');
    const startBtn = document.getElementById('start-game-btn');
    const inputs = document.getElementById('player-inputs');

    addBtn.addEventListener('click', () => {
        const count = inputs.querySelectorAll('input').length;
        if (count < 5) {
            const input = document.createElement('input');
            input.type = 'text';
            input.className = 'player-name';
            input.placeholder = `Player ${count + 1}`;
            input.dataset.index = count;
            inputs.appendChild(input);
        }
        if (count >= 4) {
            addBtn.disabled = true;
        }
    });

    startBtn.addEventListener('click', () => {
        const names = Array.from(inputs.querySelectorAll('input'))
            .map(i => i.value.trim())
            .filter(n => n.length > 0);

        if (names.length < 2) {
            alert('Need at least 2 player names');
            return;
        }

        createGame(names);
    });

    // Allow Enter key to start game
    inputs.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            startBtn.click();
        }
    });
}

// Game Rendering
function renderGame() {
    if (!gameState) return;

    // Check for game over
    if (gameState.winner) {
        document.getElementById('winner-name').textContent = gameState.winner.name;
        showScreen('gameover-screen');
        return;
    }

    // Check for response phase
    if (gameState.pendingAction && gameState.pendingAction.targets && gameState.pendingAction.targets.length > 0) {
        renderResponsePhase();
        return;
    }

    // Hide response overlay if no pending action
    document.getElementById('response-overlay').classList.add('hidden');

    const currentPlayer = gameState.players.find(p => p.id === gameState.currentPlayer);
    const otherPlayers = gameState.players.filter(p => p.id !== gameState.currentPlayer);

    // Update header
    document.getElementById('turn-indicator').textContent = `${currentPlayer.name}'s Turn`;
    const actionsLeft = gameState.maxActionsPerTurn - gameState.actionsPlayedThisTurn;
    document.getElementById('actions-remaining').textContent = `${actionsLeft} actions remaining`;

    // Render other players
    renderOtherPlayers(otherPlayers);

    // Render current player
    renderCurrentPlayer(currentPlayer);
}

function renderOtherPlayers(players) {
    const container = document.getElementById('other-players-list');
    container.innerHTML = players.map(p => `
        <div class="opponent-card">
            <h3>${escapeHtml(p.name)}</h3>
            <div class="opponent-stats">
                <div>Hand: ${p.hand.length} cards</div>
                <div>Bank: $${p.bankTotal}</div>
                <div>Sets: ${p.completeSets}/3</div>
            </div>
            <div class="opponent-properties">
                ${p.properties.map(prop => `
                    <div class="property-pip color-${prop.color} ${prop.isComplete ? 'complete' : ''}"
                         title="${prop.color}: ${prop.cards.length}/${prop.setSize}">
                        ${prop.cards.length}
                    </div>
                `).join('')}
            </div>
        </div>
    `).join('');
}

function renderCurrentPlayer(player) {
    document.getElementById('current-player-name').textContent = player.name;
    document.getElementById('current-player-stats').innerHTML = `
        <span>Bank: $${player.bankTotal}</span>
        <span>Complete Sets: ${player.completeSets}/3</span>
    `;

    // Render properties
    const propsContainer = document.getElementById('current-player-properties');
    if (player.properties.length === 0) {
        propsContainer.innerHTML = '<p class="no-properties">No properties yet</p>';
    } else {
        propsContainer.innerHTML = player.properties.map(prop => `
            <div class="property-set color-${prop.color} ${prop.isComplete ? 'complete' : ''}">
                <div class="property-set-header">
                    <span class="property-set-name">${prop.color}</span>
                    <span class="property-set-count">${prop.cards.length}/${prop.setSize}</span>
                </div>
                <div class="property-set-cards">
                    ${prop.cards.map(c => escapeHtml(c.name)).join(', ')}
                </div>
                <div class="property-set-rent">Rent: $${prop.rent}</div>
                ${prop.improvements && prop.improvements.length > 0 ? `
                    <div class="property-set-improvements">
                        + ${prop.improvements.join(', ')}
                    </div>
                ` : ''}
            </div>
        `).join('');
    }

    // Render hand
    document.getElementById('hand-count').textContent = `(${player.hand.length} cards)`;
    const handContainer = document.getElementById('current-player-hand');
    handContainer.innerHTML = player.hand.map(card => renderCard(card, player.id)).join('');

    // Attach card click handlers
    handContainer.querySelectorAll('.card').forEach(cardEl => {
        const cardId = cardEl.dataset.cardId;
        const card = player.hand.find(c => c.id === cardId);

        cardEl.querySelectorAll('.btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                handleCardAction(btn.dataset.action, card, player.id);
            });
        });
    });
}

// Map color names to CSS variable values
const colorToHex = {
    brown: '#996633',
    blue: '#0039A6',
    gray: '#808183',
    pink: '#808183',
    orange: '#FF6319',
    red: '#EE352E',
    yellow: '#FCCC0A',
    green: '#00933C',
    darkblue: '#2A344D',
    railroad: '#000000',
    black: '#000000',
    utility: '#6CBE45',
    lightgreen: '#6CBE45',
};

function renderCard(card, playerId) {
    const actionsLeft = gameState.maxActionsPerTurn - gameState.actionsPlayedThisTurn;
    const canAct = actionsLeft > 0;

    // Determine color bar styling
    let colorBarClass = 'card-color-bar';
    let colorBarStyle = '';
    let colorClass = '';

    if (card.type === 'wildcard' && card.colors && card.colors.length > 0) {
        if (card.colors.length >= 10) {
            // "Any color" wildcard (Fulton Center)
            colorBarClass += ' any-color';
        } else if (card.colors.length === 2) {
            // Dual-color wildcard
            colorBarClass += ' dual-color';
            const color1 = colorToHex[card.colors[0]] || '#444';
            const color2 = colorToHex[card.colors[1]] || '#444';
            colorBarStyle = `style="--color-1: ${color1}; --color-2: ${color2}"`;
        }
    } else if (card.type === 'rent' && card.colors && card.colors.length === 2) {
        // Dual-color rent card
        colorBarClass += ' dual-color';
        const color1 = colorToHex[card.colors[0]] || '#444';
        const color2 = colorToHex[card.colors[1]] || '#444';
        colorBarStyle = `style="--color-1: ${color1}; --color-2: ${color2}"`;
    } else if (card.type === 'rent' && card.isWildRent) {
        // Wild rent - any color
        colorBarClass += ' any-color';
    } else if (card.color) {
        colorClass = `color-${card.color}`;
    } else if (card.currentColor) {
        colorClass = `color-${card.currentColor}`;
    }

    let actions = '';
    if (canAct) {
        if (card.type === 'property' || card.type === 'wildcard') {
            actions += `<button class="btn btn-primary" data-action="play-property">Play</button>`;
        }
        if (card.value > 0) {
            actions += `<button class="btn btn-secondary" data-action="bank">Bank $${card.value}</button>`;
        }
        if (card.type === 'rent') {
            actions += `<button class="btn btn-success" data-action="play-rent">Charge Rent</button>`;
        }
        if (card.type === 'action' && card.name !== 'Fare Evasion' && card.name !== 'Rush Hour') {
            actions += `<button class="btn btn-primary" data-action="play-action">Use</button>`;
        }
    }

    // Show which colors the wildcard can be
    let colorsInfo = '';
    if (card.type === 'wildcard' && card.colors && card.colors.length > 0 && card.colors.length < 10) {
        colorsInfo = `<div class="card-effect">${card.colors.join(' / ')}</div>`;
    } else if (card.type === 'wildcard' && card.colors && card.colors.length >= 10) {
        colorsInfo = `<div class="card-effect">Any color</div>`;
    }

    return `
        <div class="card ${colorClass}" data-card-id="${card.id}">
            <div class="${colorBarClass}" ${colorBarStyle}></div>
            <div class="card-name">${escapeHtml(card.name)}</div>
            <div class="card-type">${card.type}</div>
            ${card.effect ? `<div class="card-effect">${escapeHtml(card.effect)}</div>` : ''}
            ${colorsInfo}
            <div class="card-value">$${card.value}</div>
            <div class="card-actions">${actions}</div>
        </div>
    `;
}

function handleCardAction(action, card, playerId) {
    switch (action) {
        case 'play-property':
            executeAction('PLAY_PROPERTY', playerId, { cardId: card.id });
            break;
        case 'bank':
            executeAction('PLAY_MONEY', playerId, { cardId: card.id });
            break;
        case 'play-rent':
            handlePlayRent(card, playerId);
            break;
        case 'play-action':
            handlePlayAction(card, playerId);
            break;
    }
}

function handlePlayRent(card, playerId) {
    const currentPlayer = gameState.players.find(p => p.id === playerId);

    // Choose color
    let color;
    if (card.isWildRent) {
        color = prompt('Enter color to charge rent on (e.g., blue, red, brown):');
        if (!color) return;
    } else if (card.colors && card.colors.length === 1) {
        color = card.colors[0];
    } else if (card.colors && card.colors.length > 1) {
        color = prompt(`Choose color (${card.colors.join(' or ')}):`);
        if (!color) return;
    } else {
        alert('Cannot determine rent color');
        return;
    }

    const data = { cardId: card.id, color: color.toLowerCase() };

    // For wild rent, choose target
    if (card.isWildRent) {
        const others = gameState.players.filter(p => p.id !== playerId);
        const targetChoice = prompt(`Choose target player:\n${others.map((p, i) => `${i+1}. ${p.name} (Bank: $${p.bankTotal})`).join('\n')}`);
        const targetIndex = parseInt(targetChoice) - 1;
        if (isNaN(targetIndex) || targetIndex < 0 || targetIndex >= others.length) {
            alert('Invalid choice');
            return;
        }
        data.targetPlayerId = others[targetIndex].id;
    }

    executeAction('PLAY_RENT', playerId, data);
}

function handlePlayAction(card, playerId) {
    const data = { cardId: card.id };
    let actionType;

    switch (card.name) {
        case 'Swipe In':
            actionType = 'SWIPE_IN';
            break;

        case 'Power Broker':
            actionType = 'POWER_BROKER';
            const pbTarget = selectTargetProperty('Select property to steal (from incomplete set):');
            if (!pbTarget) return;
            Object.assign(data, pbTarget);
            break;

        case 'Line Closure':
            actionType = 'LINE_CLOSURE';
            const lcTarget = selectTargetCompleteSet('Select complete set to steal:');
            if (!lcTarget) return;
            Object.assign(data, lcTarget);
            break;

        case 'Service Change':
            actionType = 'SERVICE_CHANGE';
            const currentPlayer = gameState.players.find(p => p.id === playerId);

            // Select own property to give
            const ownProp = selectOwnProperty(currentPlayer, 'Select YOUR property to give:');
            if (!ownProp) return;

            // Select target property to take
            const scTarget = selectTargetProperty('Select property to take:');
            if (!scTarget) return;

            data.playerCardId = ownProp.cardId;
            data.playerColor = ownProp.color;
            data.targetPlayerId = scTarget.targetPlayerId;
            data.targetCardId = scTarget.targetCardId;
            data.color = scTarget.color;
            break;

        case 'Missed Your Train':
            actionType = 'MISSED_YOUR_TRAIN';
            const mtTarget = selectTargetPlayer();
            if (!mtTarget) return;
            data.targetPlayerId = mtTarget;
            break;

        case "It's My Stop!":
            actionType = 'ITS_MY_STOP';
            break;

        case 'Express Service':
            actionType = 'EXPRESS_SERVICE';
            const esColor = selectOwnCompleteSet('Select complete set to add Express Service:');
            if (!esColor) return;
            data.color = esColor;
            break;

        case 'New Station':
            actionType = 'NEW_STATION';
            const nsColor = selectOwnCompleteSetWithExpress('Select complete set with Express to add New Station:');
            if (!nsColor) return;
            data.color = nsColor;
            break;

        default:
            alert(`Action "${card.name}" not implemented`);
            return;
    }

    executeAction(actionType, playerId, data);
}

function selectTargetPlayer() {
    const currentPlayer = gameState.players.find(p => p.id === gameState.currentPlayer);
    const others = gameState.players.filter(p => p.id !== currentPlayer.id);

    if (others.length === 0) {
        alert('No other players');
        return null;
    }

    const choice = prompt(`Select target player:\n${others.map((p, i) => `${i+1}. ${p.name} (Bank: $${p.bankTotal})`).join('\n')}`);
    const index = parseInt(choice) - 1;
    if (isNaN(index) || index < 0 || index >= others.length) {
        alert('Invalid choice');
        return null;
    }
    return others[index].id;
}

function selectTargetProperty(message) {
    const currentPlayer = gameState.players.find(p => p.id === gameState.currentPlayer);
    const others = gameState.players.filter(p => p.id !== currentPlayer.id);

    // Find stealable properties (not in complete sets)
    let options = [];
    others.forEach(p => {
        p.properties.forEach(prop => {
            if (!prop.isComplete) {
                prop.cards.forEach(card => {
                    options.push({ playerId: p.id, playerName: p.name, card, color: prop.color });
                });
            }
        });
    });

    if (options.length === 0) {
        alert('No stealable properties (all sets are complete or empty)');
        return null;
    }

    const choice = prompt(`${message}\n${options.map((o, i) => `${i+1}. ${o.playerName}'s ${o.card.name} (${o.color})`).join('\n')}`);
    const index = parseInt(choice) - 1;
    if (isNaN(index) || index < 0 || index >= options.length) {
        alert('Invalid choice');
        return null;
    }

    const selected = options[index];
    return { targetPlayerId: selected.playerId, targetCardId: selected.card.id, color: selected.color };
}

function selectOwnProperty(player, message) {
    let options = [];
    player.properties.forEach(prop => {
        if (!prop.isComplete) {
            prop.cards.forEach(card => {
                options.push({ card, color: prop.color });
            });
        }
    });

    if (options.length === 0) {
        alert('You have no properties to trade (or all are in complete sets)');
        return null;
    }

    const choice = prompt(`${message}\n${options.map((o, i) => `${i+1}. ${o.card.name} (${o.color})`).join('\n')}`);
    const index = parseInt(choice) - 1;
    if (isNaN(index) || index < 0 || index >= options.length) {
        alert('Invalid choice');
        return null;
    }

    const selected = options[index];
    return { cardId: selected.card.id, color: selected.color };
}

function selectTargetCompleteSet(message) {
    const currentPlayer = gameState.players.find(p => p.id === gameState.currentPlayer);
    const others = gameState.players.filter(p => p.id !== currentPlayer.id);

    let options = [];
    others.forEach(p => {
        p.properties.forEach(prop => {
            if (prop.isComplete) {
                options.push({ playerId: p.id, playerName: p.name, color: prop.color, cards: prop.cards });
            }
        });
    });

    if (options.length === 0) {
        alert('No complete sets to steal');
        return null;
    }

    const choice = prompt(`${message}\n${options.map((o, i) => `${i+1}. ${o.playerName}'s ${o.color} set (${o.cards.map(c => c.name).join(', ')})`).join('\n')}`);
    const index = parseInt(choice) - 1;
    if (isNaN(index) || index < 0 || index >= options.length) {
        alert('Invalid choice');
        return null;
    }

    const selected = options[index];
    return { targetPlayerId: selected.playerId, color: selected.color };
}

function selectOwnCompleteSet(message) {
    const currentPlayer = gameState.players.find(p => p.id === gameState.currentPlayer);
    const completeSets = currentPlayer.properties.filter(p =>
        p.isComplete && p.color !== 'railroad' && p.color !== 'utility'
    );

    if (completeSets.length === 0) {
        alert('No eligible complete sets (cannot add improvements to Railroad/Utility)');
        return null;
    }

    const choice = prompt(`${message}\n${completeSets.map((p, i) => `${i+1}. ${p.color} (Rent: $${p.rent})`).join('\n')}`);
    const index = parseInt(choice) - 1;
    if (isNaN(index) || index < 0 || index >= completeSets.length) {
        alert('Invalid choice');
        return null;
    }

    return completeSets[index].color;
}

function selectOwnCompleteSetWithExpress(message) {
    const currentPlayer = gameState.players.find(p => p.id === gameState.currentPlayer);
    const setsWithExpress = currentPlayer.properties.filter(p =>
        p.isComplete &&
        p.improvements &&
        p.improvements.includes('express') &&
        !p.improvements.includes('station')
    );

    if (setsWithExpress.length === 0) {
        alert('No complete sets with Express Service (and without New Station already)');
        return null;
    }

    const choice = prompt(`${message}\n${setsWithExpress.map((p, i) => `${i+1}. ${p.color}`).join('\n')}`);
    const index = parseInt(choice) - 1;
    if (isNaN(index) || index < 0 || index >= setsWithExpress.length) {
        alert('Invalid choice');
        return null;
    }

    return setsWithExpress[index].color;
}

// Response Phase
function renderResponsePhase() {
    const overlay = document.getElementById('response-overlay');
    overlay.classList.remove('hidden');

    const pending = gameState.pendingAction;
    const sourcePlayer = gameState.players.find(p => p.id === pending.sourcePlayer);
    const targetId = pending.targets[0]; // First unresponded target
    const targetPlayer = gameState.players.find(p => p.id === targetId);

    document.getElementById('response-title').textContent =
        `${targetPlayer.name} must respond!`;

    let description = '';
    let info = '';

    switch (pending.type) {
        case 'PLAY_RENT':
            description = `${sourcePlayer.name} is charging rent`;
            info = `<strong>Color:</strong> ${pending.rentColor}<br>
                    <strong>Amount:</strong> $${pending.rentAmount}<br>
                    <strong>Your Bank:</strong> $${targetPlayer.bankTotal}`;
            break;
        case 'MISSED_YOUR_TRAIN':
            description = `${sourcePlayer.name} played Missed Your Train`;
            info = `<strong>Amount:</strong> $${pending.rentAmount}<br>
                    <strong>Your Bank:</strong> $${targetPlayer.bankTotal}`;
            break;
        case 'ITS_MY_STOP':
            description = `${sourcePlayer.name} played It's My Stop!`;
            info = `<strong>Amount:</strong> $${pending.rentAmount}<br>
                    <strong>Your Bank:</strong> $${targetPlayer.bankTotal}`;
            break;
        case 'POWER_BROKER':
            description = `${sourcePlayer.name} is trying to steal a property`;
            info = '<strong>Effect:</strong> They want to take one of your properties!';
            break;
        case 'LINE_CLOSURE':
            description = `${sourcePlayer.name} is trying to steal a complete set`;
            info = '<strong>Effect:</strong> They want to take your complete property set!';
            break;
        case 'SERVICE_CHANGE':
            description = `${sourcePlayer.name} wants to swap properties`;
            info = '<strong>Effect:</strong> Forced property exchange!';
            break;
        default:
            description = `${sourcePlayer.name} played an action against you`;
            info = '';
    }

    document.getElementById('response-description').textContent = description;
    document.getElementById('response-info').innerHTML = info;

    // Check if target has Fare Evasion
    const hasFareEvasion = targetPlayer.hand.some(c => c.name === 'Fare Evasion');
    const fareEvasionBtn = document.getElementById('fare-evasion-btn');
    fareEvasionBtn.style.display = 'inline-block'; // Always show (no info leak)

    // Set up handlers
    document.getElementById('accept-btn').onclick = () => {
        executeAction('ACCEPT', targetId, {});
    };

    fareEvasionBtn.onclick = () => {
        if (!hasFareEvasion) {
            alert("You don't have a Fare Evasion card! Accepting instead...");
            executeAction('ACCEPT', targetId, {});
            return;
        }
        const feCard = targetPlayer.hand.find(c => c.name === 'Fare Evasion');
        executeAction('PLAY_FARE_EVASION', targetId, { cardId: feCard.id });
    };
}

// End Turn
function initGameControls() {
    document.getElementById('end-turn-btn').addEventListener('click', () => {
        const currentPlayer = gameState.players.find(p => p.id === gameState.currentPlayer);
        executeAction('END_TURN', currentPlayer.id, {});
    });

    document.getElementById('new-game-btn').addEventListener('click', () => {
        // Reset the form
        const inputs = document.getElementById('player-inputs');
        inputs.innerHTML = `
            <input type="text" class="player-name" placeholder="Player 1" data-index="0">
            <input type="text" class="player-name" placeholder="Player 2" data-index="1">
        `;
        document.getElementById('add-player-btn').disabled = false;
        showScreen('setup-screen');
    });
}

// Utility
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    initSetup();
    initGameControls();
});
