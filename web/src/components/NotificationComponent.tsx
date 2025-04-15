'use client';

import React, { useState, useEffect } from 'react';
import { useWebSocket } from './WebSocketProvider';
import { Notification } from '@/lib/socket';

export default function NotificationComponent() {
  const { isConnected, notifications: wsNotifications, markAllAsRead } = useWebSocket();
  const [isOpen, setIsOpen] = useState(false);
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  const fetchNotifications = async () => {
    try {
      setLoading(true);
      const response = await fetch(`/api/notifications?userId=usr_uvwxy&page=${page}&limit=10`);
      if (!response.ok) {
        throw new Error('Failed to fetch notifications');
      }
      const data = await response.json();
      const newNotifications = data.notifications;

      if (page === 1) {
        setNotifications(newNotifications);
      } else {
        setNotifications(prev => [...prev, ...newNotifications]);
      }

      setHasMore(newNotifications.length === 10);
    } catch (error) {
      console.error('Error fetching notifications:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadMore = () => {
    if (!loading && hasMore) {
      setPage(prev => prev + 1);
    }
  };

  useEffect(() => {
    if (isOpen) {
      fetchNotifications();
    }
  }, [isOpen, page]);

  // Merge websocket notifications with fetched notifications
  useEffect(() => {
    if (wsNotifications.length > 0) {
      setNotifications(prev => [...wsNotifications, ...prev]);
    }
  }, [wsNotifications]);

  const toggleNotifications = () => {
    setIsOpen(!isOpen);
    if (!isOpen) {
      setPage(1); // Reset page when opening
    }
  };

  const handleMarkAllAsRead = () => {
    markAllAsRead();
  };

  const handleNotificationClick = (notification: Notification) => {
    if (notification.link) {
      window.open(notification.link, '_blank');
    }
  };

  return (
    <div className="relative">
      {/* Notification Bell */}
      <button
        onClick={toggleNotifications}
        className="relative p-2 text-gray-700 hover:bg-gray-100 rounded-full"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
          />
        </svg>
        {notifications.length > 0 && (
          <span className="absolute top-0 right-0 inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-white transform translate-x-1/2 -translate-y-1/2 bg-red-500 rounded-full">
            {notifications.length}
          </span>
        )}
      </button>

      {/* Connection Status Indicator */}
      <div className={`absolute -bottom-2 -right-2 h-3 w-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`}></div>

      {/* Notification Panel */}
      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-white rounded-md shadow-lg overflow-hidden z-20">
          <div className="py-2 px-3 bg-gray-100 flex justify-between items-center">
            <span className="text-sm font-medium">Notifications</span>
            {notifications.length > 0 && (
              <button
                onClick={handleMarkAllAsRead}
                className="text-xs text-blue-500 hover:text-blue-700"
              >
                Mark all as read
              </button>
            )}
          </div>
          <div className="max-h-96 overflow-y-auto">
            {notifications.length === 0 ? (
              <div className="py-4 px-3 text-center text-gray-500">
                {loading ? 'Loading notifications...' : 'No new notifications'}
              </div>
            ) : (
              <>
                <ul>
                  {notifications.map((notification, index) => (
                    <li
                      key={index}
                      className="border-b border-gray-200 last:border-0 hover:bg-gray-50 cursor-pointer"
                      onClick={() => handleNotificationClick(notification)}
                    >
                      <div className="py-3 px-4">
                        <div className="flex items-start">
                          <div className="flex-shrink-0 bg-blue-100 rounded-full p-1">
                            <svg className="h-4 w-4 text-blue-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                              <path d="M10 2a6 6 0 00-6 6v3.586l-.707.707A1 1 0 004 14h12a1 1 0 00.707-1.707L16 11.586V8a6 6 0 00-6-6zM10 18a3 3 0 01-3-3h6a3 3 0 01-3 3z" />
                            </svg>
                          </div>
                          <div className="ml-3 w-0 flex-1">
                            <p className="text-sm font-medium text-gray-900">{notification.title}</p>
                            <p className="text-sm text-gray-500">{notification.description}</p>
                            {notification.topic && (
                              <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800 mt-1">
                                {notification.topic}
                              </span>
                            )}
                          </div>
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
                {hasMore && (
                  <div className="py-2 px-3 text-center">
                    <button
                      onClick={loadMore}
                      className="text-sm text-blue-500 hover:text-blue-700"
                      disabled={loading}
                    >
                      {loading ? 'Loading...' : 'Load more'}
                    </button>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
