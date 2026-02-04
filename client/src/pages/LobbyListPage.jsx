import { useState, useEffect } from 'react';
import { useGame } from '../contexts/GameContext';

export default function LobbyListPage() {
  const { state, actions } = useGame();
  const [lobbyName, setLobbyName] = useState('');
  const [showCreate, setShowCreate] = useState(false);

  useEffect(() => {
    // Refresh lobbies on mount and periodically
    actions.refreshLobbies();
    const interval = setInterval(() => actions.refreshLobbies(), 5000);
    return () => clearInterval(interval);
  }, []);

  const handleCreateLobby = (e) => {
    e.preventDefault();
    if (lobbyName.trim()) {
      actions.createLobby(lobbyName.trim());
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 p-4">
      <div className="max-w-2xl mx-auto">
        <div className="flex justify-between items-center mb-6">
          <div>
            <h1 className="text-2xl font-bold text-white">Game Lobbies</h1>
            <p className="text-gray-400">Welcome, {state.playerName}</p>
          </div>
          <button
            onClick={() => setShowCreate(!showCreate)}
            className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors"
          >
            {showCreate ? 'Cancel' : 'Create Lobby'}
          </button>
        </div>

        {showCreate && (
          <form onSubmit={handleCreateLobby} className="bg-gray-800 p-4 rounded-lg mb-6">
            <div className="flex gap-3">
              <input
                type="text"
                value={lobbyName}
                onChange={(e) => setLobbyName(e.target.value)}
                placeholder="Lobby name..."
                className="flex-1 px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-blue-500"
                required
              />
              <button
                type="submit"
                className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
              >
                Create
              </button>
            </div>
          </form>
        )}

        <div className="space-y-3">
          {state.lobbies.length === 0 ? (
            <div className="bg-gray-800 p-8 rounded-lg text-center">
              <p className="text-gray-400">No open lobbies. Create one to get started!</p>
            </div>
          ) : (
            state.lobbies.map((lobby) => (
              <div
                key={lobby.id}
                className="bg-gray-800 p-4 rounded-lg flex justify-between items-center"
              >
                <div>
                  <h3 className="text-white font-semibold">{lobby.name}</h3>
                  <p className="text-gray-400 text-sm">
                    {lobby.players.length}/{lobby.maxPlayers} players
                  </p>
                </div>
                <button
                  onClick={() => actions.joinLobby(lobby.id)}
                  className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                >
                  Join
                </button>
              </div>
            ))
          )}
        </div>

        <button
          onClick={() => actions.refreshLobbies()}
          className="mt-4 text-gray-400 hover:text-white text-sm"
        >
          Refresh lobbies
        </button>
      </div>
    </div>
  );
}
