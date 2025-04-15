'use client';

import NotificationComponent from '@/components/NotificationComponent';
import UserOrderList from '@/components/UserOrderList';

export default function Home() {
  const userId = 'usr_uvwxy'; // Default user ID - in a real app, this would come from authentication

  return (
    <main className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="py-8">
          <div className="flex justify-between items-center pb-6 border-b border-gray-200">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
              <p className="mt-1 text-sm text-gray-500">Welcome to your application dashboard</p>
            </div>
			<NotificationComponent />
          </div>

          <div className="mt-8">
            <UserOrderList userId={userId} />
          </div>

          <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* Sample content cards */}
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-white shadow rounded-lg p-6">
                <h2 className="text-xl font-medium text-gray-900 mb-4">Content Card {i}</h2>
                <p className="text-gray-500">This is a sample content card for demonstration purposes.</p>
                <div className="mt-4 pt-4 border-t border-gray-200">
                  <button className="text-blue-600 hover:text-blue-800 text-sm font-medium">
                    Learn more
                  </button>
                </div>
              </div>
            ))}
          </div>

          <footer className="mt-10 pt-6 border-t border-gray-200 text-center text-sm text-gray-500">
            <p>WebSocket Notification Demo</p>
            <p>Notification WebSocket URL: ws://localhost:8382/notifications</p>
          </footer>
        </div>
      </div>
    </main>
  );
}
