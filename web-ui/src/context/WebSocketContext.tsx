import React, { createContext, useContext, useRef } from "react";
import { WebSocketService } from "../services/websocketService";

const WebSocketContext = createContext(null);

export const WebSocketProvider = ({ children }: { children: React.ReactNode }) => {
    const wsServiceRef = useRef(null);

    const connectWebSocket = (url: string) => {
        if (!wsServiceRef.current) {
            console.log("Connecting to WebSocket");
            const wsService = WebSocketService.getInstance(url);
            wsService.connect();
            wsServiceRef.current = wsService;
        }
    };

    const getWebSocketService = () => wsServiceRef.current;

    return (
        <WebSocketContext.Provider value={{ connectWebSocket, getWebSocketService }}>
            {children}
        </WebSocketContext.Provider>
    );
};

export const useWebSocket = () => useContext(WebSocketContext);

