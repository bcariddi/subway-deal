import { useGame } from '../contexts/GameContext';

export default function GameOverPage() {
  const { state, actions } = useGame();

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="bg-gray-800 p-8 rounded-lg text-center max-w-md w-full">
        <h1 className="text-4xl font-bold text-white mb-4">Game Over!</h1>
        <h2 className="text-2xl text-green-400 mb-2">{state.winner}</h2>
        <p className="text-gray-400 mb-8">completed 3 property sets!</p>
        <button
          onClick={() => {
            actions.setScreen('home');
          }}
          className="px-8 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition-colors"
        >
          Play Again
        </button>
      </div>
    </div>
  );
}
