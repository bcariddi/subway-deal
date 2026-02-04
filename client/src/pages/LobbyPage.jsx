import { useGame } from '../contexts/GameContext';

export default function LobbyPage() {
  const { state, actions } = useGame();
  const { lobby, playerName } = state;

  if (!lobby) {
    return <div className="min-h-screen bg-gray-900 flex items-center justify-center text-white">Loading...</div>;
  }

  const isHost = lobby.players.find(p => p.name === playerName)?.clientId === lobby.hostId;
  const myPlayer = lobby.players.find(p => p.name === playerName);
  const allReady = lobby.players.length >= 2 && lobby.players.every(p => p.ready);

  return (
    <div className="min-h-screen bg-gray-900 p-4">
      <div className="max-w-2xl mx-auto">
        <div className="bg-gray-800 rounded-lg p-6 mb-6">
          <div className="flex justify-between items-start mb-4">
            <div>
              <h1 className="text-2xl font-bold text-white">{lobby.name}</h1>
              <p className="text-gray-400">
                {lobby.players.length}/{lobby.maxPlayers} players
              </p>
            </div>
            <button
              onClick={() => actions.leaveLobby()}
              className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
            >
              Leave
            </button>
          </div>

          <div className="space-y-3">
            {lobby.players.map((player, index) => (
              <div
                key={player.clientId}
                className={`flex justify-between items-center p-3 rounded-lg ${
                  player.name === playerName ? 'bg-blue-900/30 border border-blue-700' : 'bg-gray-700'
                }`}
              >
                <div className="flex items-center gap-3">
                  <span className="text-white font-medium">{player.name}</span>
                  {player.clientId === lobby.hostId && (
                    <span className="px-2 py-0.5 bg-yellow-600 text-xs text-white rounded">Host</span>
                  )}
                  {player.name === playerName && (
                    <span className="text-gray-400 text-sm">(You)</span>
                  )}
                </div>
                <div className={`px-3 py-1 rounded-full text-sm ${
                  player.ready ? 'bg-green-600 text-white' : 'bg-gray-600 text-gray-300'
                }`}>
                  {player.ready ? 'Ready' : 'Not Ready'}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="flex gap-4">
          {myPlayer && !myPlayer.ready ? (
            <button
              onClick={() => actions.setReady(true)}
              className="flex-1 py-3 bg-green-600 hover:bg-green-700 text-white font-semibold rounded-lg transition-colors"
            >
              Ready Up
            </button>
          ) : myPlayer && myPlayer.ready ? (
            <button
              onClick={() => actions.setReady(false)}
              className="flex-1 py-3 bg-gray-600 hover:bg-gray-700 text-white font-semibold rounded-lg transition-colors"
            >
              Cancel Ready
            </button>
          ) : null}

          {isHost && (
            <button
              onClick={() => actions.startGame()}
              disabled={!allReady}
              className={`flex-1 py-3 font-semibold rounded-lg transition-colors ${
                allReady
                  ? 'bg-blue-600 hover:bg-blue-700 text-white'
                  : 'bg-gray-700 text-gray-500 cursor-not-allowed'
              }`}
            >
              {allReady ? 'Start Game' : 'Waiting for players...'}
            </button>
          )}
        </div>

        {!isHost && (
          <p className="text-gray-400 text-center mt-4">
            Waiting for the host to start the game...
          </p>
        )}
      </div>
    </div>
  );
}
