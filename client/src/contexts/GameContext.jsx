import { createContext, useContext, useReducer, useEffect } from 'react';
import { useWebSocket } from '../hooks/useWebSocket';

const GameContext = createContext(null);

export function useGame() {
  const context = useContext(GameContext);
  if (!context) throw new Error('useGame must be used within GameProvider');
  return context;
}

const initialState = {
  screen: 'home',
  playerName: '',
  lobby: null,
  lobbies: [],
  gameState: null,
  yourHand: [],
  yourId: null,
  error: null,
  winner: null,
};

function reducer(state, action) {
  switch (action.type) {
    case 'SET_SCREEN':
      return { ...state, screen: action.payload, error: null };
    case 'SET_PLAYER_NAME':
      return { ...state, playerName: action.payload };
    case 'SET_LOBBY':
      return { ...state, lobby: action.payload, screen: 'lobby' };
    case 'UPDATE_LOBBY':
      return { ...state, lobby: action.payload };
    case 'SET_LOBBIES':
      return { ...state, lobbies: action.payload || [] };
    case 'SET_GAME_STATE':
      return {
        ...state,
        gameState: action.payload,
        yourHand: action.payload.yourHand || [],
        yourId: action.payload.yourId,
        screen: 'game',
      };
    case 'SET_ERROR':
      return { ...state, error: action.payload };
    case 'GAME_ENDED':
      return { ...state, screen: 'game-over', winner: action.payload };
    case 'RESET':
      return { ...initialState, playerName: state.playerName };
    default:
      return state;
  }
}

export function GameProvider({ children }) {
  const [state, dispatch] = useReducer(reducer, initialState);
  // In production (Docker), use relative WebSocket URL via nginx proxy
  // In development, use explicit localhost URL
  const serverUrl = import.meta.env.VITE_SERVER_URL ||
    (window.location.protocol === 'https:' ? 'wss:' : 'ws:') + '//' + window.location.host + '/ws';
  const { connected, send, subscribe } = useWebSocket(serverUrl);

  // Subscribe to server messages
  useEffect(() => {
    const unsubscribes = [
      subscribe('lobby:created', (data) => dispatch({ type: 'SET_LOBBY', payload: data })),
      subscribe('lobby:updated', (data) => dispatch({ type: 'UPDATE_LOBBY', payload: data })),
      subscribe('lobby:list', (data) => dispatch({ type: 'SET_LOBBIES', payload: data })),
      subscribe('game:state', (data) => dispatch({ type: 'SET_GAME_STATE', payload: data })),
      subscribe('game:ended', (data) => dispatch({ type: 'GAME_ENDED', payload: data.winner })),
      subscribe('error', (data) => dispatch({ type: 'SET_ERROR', payload: data.message })),
    ];

    return () => unsubscribes.forEach(unsub => unsub());
  }, [subscribe]);

  const actions = {
    setPlayerName: (name) => dispatch({ type: 'SET_PLAYER_NAME', payload: name }),
    setScreen: (screen) => dispatch({ type: 'SET_SCREEN', payload: screen }),

    createLobby: (lobbyName, maxPlayers = 5) => {
      send('lobby:create', {
        name: lobbyName,
        playerName: state.playerName,
        maxPlayers,
      });
    },

    joinLobby: (lobbyId) => {
      send('lobby:join', {
        lobbyId,
        playerName: state.playerName,
      });
    },

    leaveLobby: () => {
      send('lobby:leave', {});
      dispatch({ type: 'RESET' });
    },

    setReady: (ready) => {
      send('lobby:ready', { ready });
    },

    startGame: () => {
      send('lobby:start', {});
    },

    playAction: (actionType, data) => {
      send('game:action', { type: actionType, data });
    },

    endTurn: () => {
      send('game:action', { type: 'END_TURN', data: {} });
    },

    refreshLobbies: () => {
      send('lobby:list', {});
    },
  };

  return (
    <GameContext.Provider value={{ state, actions, connected }}>
      {children}
    </GameContext.Provider>
  );
}
