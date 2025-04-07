'use client';

import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import toast from 'react-hot-toast';
import { initializeSocket, closeSocket, Notification, WebSocketMessage } from '@/lib/socket';

type WebSocketContextType = {
  isConnected: boolean;
  notifications: Notification[];
  markAllAsRead: () => void;
};

const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

export const WebSocketProvider = ({
  children,
  autoConnect = true,
  url = 'ws://localhost/notifications',
  userId = 'U0001' // Default user ID
}: {
  children: ReactNode;
  autoConnect?: boolean;
  url?: string;
  userId?: string;
}) => {
  const [isConnected, setIsConnected] = useState(false);
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const connect = () => {
    const wsClient = initializeSocket(url);

    wsClient.onConnect(() => {
      setIsConnected(true);
      toast.success('Connected to notification server');

      // Send AUTHORIZATION message with user_id after successful connection
      const authMessage: WebSocketMessage = {
        topic: 'AUTHORIZATION',
        content: {
          user_id: userId
        }
      };

      // Send the authorization message
      wsClient.send(authMessage);
      console.log('Sent authorization message with user_id:', userId);
    });

    wsClient.onDisconnect(() => {
      setIsConnected(false);
      toast.error('Disconnected from notification server');
    });

    // Listen for notification messages
    wsClient.onMessage((data) => {
      console.log('Received message:', data);
      if (data.topic === 'ANNOUNCEMENT') {
        const notification = data.content as Notification;
        setNotifications(prev => [notification, ...prev]);

        // Display toast notification
        toast.custom(
          (t) => (
            <div
              className={`${
                t.visible ? 'animate-enter' : 'animate-leave'
              } max-w-md w-full bg-white shadow-lg rounded-lg pointer-events-auto flex ring-1 ring-black ring-opacity-5`}
            >
              <div className="flex-1 w-0 p-4">
                <div className="flex items-start">
                  <div className="flex-shrink-0 pt-0.5">
                    <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                      <svg className="h-6 w-6 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                      </svg>
                    </div>
                  </div>
                  <div className="ml-3 flex-1">
                    <p className="text-sm font-medium text-gray-900">
                      {notification.title}
                    </p>
                    <p className="mt-1 text-sm text-gray-500">
                      {notification.description}
                    </p>
                    {notification.topic && (
                      <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800 mt-1">
                        {notification.topic}
                      </span>
                    )}
                  </div>
                </div>
              </div>
              <div className="flex border-l border-gray-200">
                <button
                  onClick={() => {
                    toast.dismiss(t.id);
                    if (notification.link) {
                      window.open(notification.link, '_blank');
                    }
                  }}
                  className="w-full border border-transparent rounded-none rounded-r-lg p-4 flex items-center justify-center text-sm font-medium text-blue-600 hover:text-blue-500 focus:outline-none"
                >
                  View
                </button>
              </div>
            </div>
          ),
          { duration: 5000 }
        );
      }
    });
  };

  const markAllAsRead = () => {
    setNotifications([]);
  };

  useEffect(() => {
    if (autoConnect) {
      connect();
    }

    return () => {
      closeSocket();
      setIsConnected(false);
    };
  }, [autoConnect, url]);

  return (
    <WebSocketContext.Provider value={{ isConnected, notifications, markAllAsRead }}>
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (context === undefined) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
};
