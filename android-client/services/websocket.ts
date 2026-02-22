import { DeviceEventEmitter } from 'react-native';
import type { Message, ConnectionState, LogEntry, MessageType } from '@/types';

type MessageHandler = (message: Message) => void;
type StateHandler = (state: ConnectionState) => void;

const generateId = () => Math.random().toString(36).substring(2, 15);

export class EcoWebSocket {
  private ws: WebSocket | null = null;
  private url: string;
  private deviceId: string;
  private secret: string;
  private deviceName: string;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private pingInterval: ReturnType<typeof setInterval> | null = null;
  private messageHandlers: Set<MessageHandler> = new Set();
  private stateHandlers: Set<StateHandler> = new Set();
  private logHandlers: Set<(log: LogEntry) => void> = new Set();
  private _state: ConnectionState = 'disconnected';
  private isConnecting = false;
  private shouldReconnect = false;

  constructor(url: string, deviceId: string, secret: string, deviceName: string) {
    this.url = url;
    this.deviceId = deviceId;
    this.secret = secret;
    this.deviceName = deviceName;
  }

  get state(): ConnectionState {
    return this._state;
  }

  private setState(newState: ConnectionState) {
    this._state = newState;
    this.stateHandlers.forEach((h) => h(newState));
  }

  private addLog(type: MessageType, direction: 'sent' | 'received', payload: unknown) {
    const log: LogEntry = {
      id: generateId(),
      timestamp: Date.now(),
      type,
      direction,
      payload,
    };
    this.logHandlers.forEach((h) => h(log));
  }

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.isConnecting) {
        console.log('Already connecting, skipping...');
        resolve();
        return;
      }
      
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        console.log('Already connected');
        resolve();
        return;
      }

      this.isConnecting = true;
      this.shouldReconnect = true;
      this.setState('connecting');

      try {
        this.ws = new WebSocket(this.url);

        this.ws.onopen = () => {
          console.log('WebSocket opened');
          this.isConnecting = false;
          this.reconnectAttempts = 0;
          this.sendHello();
          this.startPing();
          this.setState('connected');
          this.setupNativeListeners();
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data) as Message;
            this.addLog(message.type, 'received', message.payload);
            this.messageHandlers.forEach((h) => h(message));
          } catch (e) {
            console.error('Failed to parse message:', e);
          }
        };

        this.ws.onerror = () => {
          console.log('WebSocket error');
          this.isConnecting = false;
          this.setState('error');
        };

        this.ws.onclose = () => {
          console.log('WebSocket closed');
          this.isConnecting = false;
          this.stopPing();
          this.setState('disconnected');
          
          if (this.shouldReconnect && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.attemptReconnect();
          }
        };
      } catch (error) {
        this.isConnecting = false;
        this.setState('error');
        reject(error instanceof Error ? error : new Error('Unknown WebSocket error'));
      }
    });
  }

  private sendHello() {
    this.send('device.hello', { device_name: this.deviceName });
  }

  private startPing() {
    this.pingInterval = setInterval(() => {
      if (this._state === 'connected') {
        this.send('device.ping', {});
      }
    }, 30000);
  }

  private stopPing() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
      setTimeout(() => this.connect(), delay);
    }
  }

  send<T>(type: MessageType, payload: T): boolean {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      return false;
    }

    const message: Message<T> = {
      type,
      device_id: this.deviceId,
      secret: this.secret,
      payload,
    };

    try {
      this.ws.send(JSON.stringify(message));
      this.addLog(type, 'sent', payload);
      return true;
    } catch (e) {
      console.error('Failed to send message:', e);
      return false;
    }
  }

  private setupNativeListeners() {
    try {
      DeviceEventEmitter.addListener('EcoCallStateChanged', (event) => {
        console.log('Call state changed event received:', event);
        this.send('call.incoming', { number: event.phoneNumber || '' });
      });

      DeviceEventEmitter.addListener('EcoNotificationPosted', (event) => {
        console.log('Notification posted event received:', event);
        this.send('notification.push', {
          app: event.packageName,
          title: event.title,
          body: event.text,
        });
      });
      console.log('Native event listeners setup complete');
    } catch (e) {
      console.error('Error setting up native listeners:', e);
    }
  }

  onMessage(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

  onStateChange(handler: StateHandler): () => void {
    this.stateHandlers.add(handler);
    return () => this.stateHandlers.delete(handler);
  }

  onLog(handler: (log: LogEntry) => void): () => void {
    this.logHandlers.add(handler);
    return () => this.logHandlers.delete(handler);
  }

  disconnect() {
    console.log('WebSocket disconnect called');
    this.shouldReconnect = false;
    this.isConnecting = false;
    this.reconnectAttempts = 0;
    this.stopPing();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.setState('disconnected');
    DeviceEventEmitter.removeAllListeners('EcoCallStateChanged');
    DeviceEventEmitter.removeAllListeners('EcoNotificationPosted');
  }
}
