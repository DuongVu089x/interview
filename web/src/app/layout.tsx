import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { Toaster } from "react-hot-toast";
import "./globals.css";
import { WebSocketProvider } from '@/components/WebSocketProvider';

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "WebSocket Demo",
  description: "Next.js WebSocket Client Demo",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <WebSocketProvider userId="usr_uvwxy">
          {children}
        </WebSocketProvider>
        <Toaster position="top-right" />
      </body>
    </html>
  );
}
