import { GameProvider, useGame } from './contexts/GameContext';
import HomePage from './pages/HomePage';
import LobbyListPage from './pages/LobbyListPage';
import LobbyPage from './pages/LobbyPage';
import GamePage from './pages/GamePage';
import GameOverPage from './pages/GameOverPage';
import './index.css';

function AppContent() {
  const { state, connected } = useGame();

  if (!connected) {
    return (
      <div className="min-h-screen bg-gray-900 text-white flex items-center justify-center">
        <div className="text-center">
          <div className="text-xl mb-2">Connecting to server...</div>
          <div className="text-gray-400 text-sm">Make sure the server is running on port 8081</div>
        </div>
      </div>
    );
  }

  const screens = {
    home: <HomePage />,
    'lobby-list': <LobbyListPage />,
    lobby: <LobbyPage />,
    game: <GamePage />,
    'game-over': <GameOverPage />,
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {state.error && (
        <div className="bg-red-600 text-white p-4 text-center">
          {state.error}
        </div>
      )}
      {screens[state.screen] || <HomePage />}
    </div>
  );
}

export default function App() {
  return (
    <GameProvider>
      <AppContent />
    </GameProvider>
  );
}
