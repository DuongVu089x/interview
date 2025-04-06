# WebSocket Notification Client

This is a Next.js client that connects to a WebSocket notification server running at `ws://localhost:8382/notifications`.

## Features

-   Real-time WebSocket notifications
-   Clean, modern UI with Tailwind CSS
-   TypeScript support
-   Connection status indicator
-   Notification badge counter
-   Clickable notifications with links
-   "Mark all as read" functionality

## Getting Started

1. Make sure your WebSocket server is running at `ws://localhost:8382/notifications`
2. Install dependencies:

```bash
npm install
```

3. Run the development server:

```bash
npm run dev
```

4. Open [http://localhost:3000](http://localhost:3000) with your browser to see the dashboard with notification functionality.

## How It Works

The client connects to the WebSocket server and listens for notification messages sent from the backend Kafka consumer. When an order is created, the Kafka consumer processes the message and triggers a notification that gets sent to the connected WebSocket clients.

## Configuration

The WebSocket server URL is configured in the `WebSocketProvider` component. You can change it by modifying the `url` prop:

```jsx
<WebSocketProvider url="ws://your-server-url/notifications">
    <NotificationComponent />
</WebSocketProvider>
```

## Project Structure

-   `src/lib/socket.ts` - WebSocket utility functions
-   `src/components/WebSocketProvider.tsx` - WebSocket context provider that manages notifications
-   `src/components/NotificationComponent.tsx` - Notification bell and dropdown component
-   `src/app/page.tsx` - Main dashboard page component

## Expected Notification Format

The server is expected to send notifications in this format:

```json
{
    "topic": "ANNOUNCEMENT",
    "content": {
        "topic": "order-created",
        "title": "Order Created",
        "description": "Order created successfully",
        "link": "localhost:8081/order/123"
    }
}
```

## Dependencies

-   Next.js
-   React
-   Socket.io-client
-   Tailwind CSS
