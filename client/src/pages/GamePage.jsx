import { useState } from 'react';
import { useGame } from '../contexts/GameContext';

// MTA color mapping
const colorMap = {
  brown: '#996633',
  blue: '#0039A6',
  gray: '#808183',
  orange: '#FF6319',
  red: '#EE352E',
  yellow: '#FCCC0A',
  green: '#00933C',
  darkblue: '#2A344D',
  black: '#000000',
  lightgreen: '#6CBE45',
  railroad: '#000000',
  utility: '#6CBE45',
};

export default function GamePage() {
  const { state, actions } = useGame();
  const { gameState, yourHand, yourId } = state;
  const [selectedCard, setSelectedCard] = useState(null);

  if (!gameState) {
    return <div className="min-h-screen bg-gray-900 flex items-center justify-center text-white">Loading game...</div>;
  }

  const isYourTurn = gameState.currentPlayer === yourId;
  const currentPlayer = gameState.players?.find(p => p.id === gameState.currentPlayer);
  const actionsLeft = gameState.maxActionsPerTurn - gameState.actionsPlayedThisTurn;
  const hasPendingAction = gameState.pendingAction && gameState.pendingAction.targets?.includes(yourId);

  const handlePlayProperty = (card) => {
    actions.playAction('PLAY_PROPERTY', { cardId: card.id });
    setSelectedCard(null);
  };

  const handleBankCard = (card) => {
    actions.playAction('PLAY_MONEY', { cardId: card.id });
    setSelectedCard(null);
  };

  const handlePlayRent = (card) => {
    // For simplicity, we'll just use the first available color
    const color = card.colors?.[0] || 'blue';
    actions.playAction('PLAY_RENT', { cardId: card.id, color });
    setSelectedCard(null);
  };

  const handleAccept = () => {
    actions.playAction('ACCEPT', {});
  };

  const handlePlayFareEvasion = () => {
    const fareEvasionCard = yourHand.find(c => c.name === 'Fare Evasion');
    if (fareEvasionCard) {
      actions.playAction('PLAY_FARE_EVASION', { cardId: fareEvasionCard.id });
    } else {
      alert("You don't have a Fare Evasion card!");
    }
  };

  // Response phase overlay
  if (hasPendingAction) {
    const sourcePlayer = gameState.players.find(p => p.id === gameState.pendingAction.sourcePlayer);
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
        <div className="bg-gray-800 p-8 rounded-lg max-w-md w-full border-2 border-red-600">
          <h2 className="text-2xl font-bold text-red-500 mb-4">You Must Respond!</h2>
          <p className="text-gray-300 mb-4">
            {sourcePlayer?.name} is demanding payment
          </p>
          <div className="bg-gray-700 p-4 rounded-lg mb-6">
            <p className="text-white">
              <strong>Amount:</strong> ${gameState.pendingAction.rentAmount}
            </p>
            {gameState.pendingAction.rentColor && (
              <p className="text-white">
                <strong>Color:</strong> {gameState.pendingAction.rentColor}
              </p>
            )}
          </div>
          <div className="flex gap-4">
            <button
              onClick={handleAccept}
              className="flex-1 py-3 bg-red-600 hover:bg-red-700 text-white font-semibold rounded-lg"
            >
              Accept & Pay
            </button>
            <button
              onClick={handlePlayFareEvasion}
              className="flex-1 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg"
            >
              Fare Evasion
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 p-4">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="bg-gray-800 rounded-lg p-4 mb-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-white">Subway Deal</h1>
          <div className="text-lg">
            {isYourTurn ? (
              <span className="text-green-400 font-bold">Your Turn ({actionsLeft} actions left)</span>
            ) : (
              <span className="text-gray-400">{currentPlayer?.name}'s Turn</span>
            )}
          </div>
        </div>

        {/* Players Grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          {gameState.players?.map((player) => (
            <div
              key={player.id}
              className={`bg-gray-800 p-4 rounded-lg ${
                player.id === yourId ? 'ring-2 ring-blue-500' : ''
              } ${player.id === gameState.currentPlayer ? 'ring-2 ring-green-500' : ''}`}
            >
              <h3 className="font-bold text-white">
                {player.name} {player.id === yourId && '(You)'}
              </h3>
              <div className="text-sm text-gray-400 mt-2 space-y-1">
                <div>Cards: {player.handCount}</div>
                <div>Bank: ${player.bankValue}</div>
                <div>Sets: {player.completeSets}/3</div>
              </div>
              {/* Properties */}
              {player.properties?.length > 0 && (
                <div className="mt-2 flex flex-wrap gap-1">
                  {player.properties.map((p, i) => (
                    <span
                      key={i}
                      className="inline-block px-2 py-1 rounded text-xs text-white"
                      style={{ backgroundColor: p.complete ? '#00933C' : colorMap[p.color] || '#666' }}
                    >
                      {p.color}: {p.cardCount}
                    </span>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Your Hand */}
        <div className="bg-gray-800 rounded-lg p-4 mb-4">
          <h2 className="text-xl font-bold text-white mb-4">Your Hand ({yourHand.length} cards)</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3">
            {yourHand.map((card) => (
              <div
                key={card.id}
                onClick={() => setSelectedCard(selectedCard?.id === card.id ? null : card)}
                className={`bg-gray-700 rounded-lg p-3 border-2 cursor-pointer transition-all ${
                  selectedCard?.id === card.id ? 'border-blue-500 scale-105' : 'border-gray-600 hover:border-gray-500'
                }`}
              >
                {/* Color bar */}
                <div
                  className="h-2 rounded mb-2"
                  style={{
                    background: card.colors?.length === 2
                      ? `linear-gradient(90deg, ${colorMap[card.colors[0]] || '#666'} 50%, ${colorMap[card.colors[1]] || '#666'} 50%)`
                      : colorMap[card.color] || colorMap[card.colors?.[0]] || '#666'
                  }}
                />
                <div className="font-bold text-sm text-white">{card.name}</div>
                <div className="text-xs text-gray-400 capitalize">{card.type}</div>
                {card.effect && <div className="text-xs text-gray-500 mt-1">{card.effect}</div>}
                <div className="text-xs text-green-400 mt-1">${card.value}</div>

                {/* Action buttons */}
                {isYourTurn && actionsLeft > 0 && selectedCard?.id === card.id && (
                  <div className="mt-2 flex gap-1 flex-wrap">
                    {(card.type === 'property' || card.type === 'wildcard') && (
                      <button
                        onClick={(e) => { e.stopPropagation(); handlePlayProperty(card); }}
                        className="text-xs bg-blue-600 hover:bg-blue-700 px-2 py-1 rounded text-white"
                      >
                        Play
                      </button>
                    )}
                    {card.value > 0 && (
                      <button
                        onClick={(e) => { e.stopPropagation(); handleBankCard(card); }}
                        className="text-xs bg-green-600 hover:bg-green-700 px-2 py-1 rounded text-white"
                      >
                        Bank
                      </button>
                    )}
                    {card.type === 'rent' && (
                      <button
                        onClick={(e) => { e.stopPropagation(); handlePlayRent(card); }}
                        className="text-xs bg-orange-600 hover:bg-orange-700 px-2 py-1 rounded text-white"
                      >
                        Rent
                      </button>
                    )}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* End Turn Button */}
        {isYourTurn && (
          <div className="bg-gray-800 rounded-lg p-4">
            <button
              onClick={() => actions.endTurn()}
              className="px-6 py-3 bg-red-600 hover:bg-red-700 text-white font-bold rounded-lg"
            >
              End Turn
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
