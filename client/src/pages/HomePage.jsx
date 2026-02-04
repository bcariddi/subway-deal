import { useState } from 'react';
import { useGame } from '../contexts/GameContext';

export default function HomePage() {
  const { state, actions } = useGame();
  const [name, setName] = useState(state.playerName);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (name.trim()) {
      actions.setPlayerName(name.trim());
      actions.setScreen('lobby-list');
      actions.refreshLobbies();
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900">
      <div className="bg-gray-800 p-8 rounded-lg shadow-xl max-w-md w-full mx-4">
        <h1 className="text-4xl font-bold text-center mb-2 bg-gradient-to-r from-blue-500 to-orange-500 bg-clip-text text-transparent">
          Subway Deal
        </h1>
        <p className="text-gray-400 text-center mb-8">NYC Transit Card Game</p>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="name" className="block text-sm font-medium text-gray-300 mb-2">
              Your Name
            </label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter your name"
              className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-blue-500"
              required
            />
          </div>

          <button
            type="submit"
            className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition-colors"
          >
            Enter Lobby
          </button>
        </form>

        <p className="text-gray-500 text-sm text-center mt-6">
          First to complete 3 property sets wins!
        </p>
      </div>
    </div>
  );
}
