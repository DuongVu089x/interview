'use client';

import React, { useState, useEffect } from 'react';
import { useWebSocket } from './WebSocketProvider';

type Message = {
  id: string;
  content: string;
  timestamp: string;
};

export default function WebSocketDemo() {
  const { socket, isConnected } = useWebSocket();
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState<string>('');

  useEffect(() => {
    if (!socket) return;

    // Listen for incoming messages
    socket.on('message', (message: Message) => {
      setMessages((prev) => [...prev, message]);
    });

    // Clean up the event listener when component unmounts
    return () => {
      socket.off('message');
    };
  }, [socket]);

  const sendMessage = () => {
    if (!socket || !inputMessage.trim()) return;

    // Send message to server
    socket.emit('message', {
      content: inputMessage
    });

    setInputMessage('');
  };

  return (
    <div className="p-4 max-w-2xl mx-auto">
      <div className="mb-4">
        <div className={`px-3 py-1 rounded-full text-sm inline-block mb-4 ${
          isConnected ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
        }`}>
          {isConnected ? 'Connected' : 'Disconnected'}
        </div>

        <div className="border rounded-lg p-4 h-80 overflow-y-auto mb-4 bg-gray-50">
          {messages.length === 0 ? (
            <p className="text-gray-500 text-center mt-32">No messages yet</p>
          ) : (
            messages.map((msg) => (
              <div key={msg.id} className="mb-2 p-2 rounded bg-white shadow-sm">
                <p>{msg.content}</p>
                <p className="text-xs text-gray-500">{new Date(msg.timestamp).toLocaleTimeString()}</p>
              </div>
            ))
          )}
        </div>

        <div className="flex">
          <input
            type="text"
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
            placeholder="Type a message..."
            className="flex-1 border rounded-l-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={!isConnected}
          />
          <button
            onClick={sendMessage}
            disabled={!isConnected}
            className="bg-blue-500 text-white px-4 py-2 rounded-r-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-300"
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
}
