import React, { createContext, useContext, useState, useEffect, useCallback, useRef } from 'react';
import { DeviceEventEmitter, Platform, AppState, AppStateStatus } from 'react-native';
import { EcoWebSocket, NativeClipboard, NativeCall, ecoEventEmitter, isNativeAvailable } from '@/services';
import { requestAndroidPermissions } from '@/services/permissions';
import type { ConnectionState, LogEntry, Message, ClipboardPayload, ClipboardEvent, CallStateEvent, NotificationEvent } from '@/types';

const getDefaultUrl = () => {
  if (Platform.OS === 'android') {
    return 'ws://10.0.2.2:4949/ws';
  }
  return 'ws://localhost:4949/ws';
};

interface EcoContextType {
  state: ConnectionState;
  logs: LogEntry[];
  connect: (url: string, deviceId: string, secret: string, deviceName: string) => Promise<void>;
  disconnect: () => void;
  sendClipboard: (text: string) => boolean;
  serverUrl: string;
  deviceId: string;
  secret: string;
  setConfig: (url: string, deviceId: string, secret: string) => void;
  nativeEnabled: boolean;
  nativeAvailable: boolean;
  enableNativeListeners: () => Promise<void>;
  disableNativeListeners: () => Promise<void>;
}

const EcoContext = createContext<EcoContextType | null>(null);

export function EcoProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<ConnectionState>('disconnected');
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [serverUrl, setServerUrl] = useState(getDefaultUrl());
  const [deviceId, setDeviceId] = useState('');
  const [secret, setSecret] = useState('');
  const [nativeEnabled, setNativeEnabled] = useState(false);
  const wsRef = useRef<EcoWebSocket | null>(null);

  const setConfig = useCallback((url: string, devId: string, sec: string) => {
    setServerUrl(url);
    setDeviceId(devId);
    setSecret(sec);
  }, []);

  const addLog = useCallback((type: LogEntry['type'], direction: LogEntry['direction'], payload: unknown) => {
    const log: LogEntry = {
      id: Math.random().toString(36).substring(2, 15),
      timestamp: Date.now(),
      type,
      direction,
      payload,
    };
    setLogs((prev) => [log, ...prev].slice(0, 100));
  }, []);

  const enableNativeListeners = useCallback(async () => {
    console.log('isNativeAvailable:', isNativeAvailable);

    if (Platform.OS === 'android') {
      console.log('Requesting Android permissions...');
      const perms = await requestAndroidPermissions();
      console.log('Permissions granted:', perms);
    }
    
    if (!isNativeAvailable) {
      console.log('Native modules not available, only using expo-clipboard');
      setNativeEnabled(true);
      return;
    }
    
    try {
      console.log('Starting call listener...');
      const callResult = await NativeCall.startListening();
      console.log('Call listener result:', callResult);
      setNativeEnabled(true);
      console.log('Native listeners enabled');
    } catch (e: unknown) {
      console.error('Failed to enable native listeners:', e);
      if (e instanceof Error) {
        console.error('Error name:', e.name);
        console.error('Error message:', e.message);
      }
    }
  }, []);

  const disableNativeListeners = useCallback(async () => {
    if (!isNativeAvailable) return;
    try {
      await NativeCall.stopListening();
      setNativeEnabled(false);
    } catch (e) {
      console.error('Failed to disable native listeners:', e);
    }
  }, []);

  const connect = useCallback(async (url: string, devId: string, sec: string, deviceName: string) => {
    console.log('Connecting to:', url);
    console.log('Device ID:', devId);
    console.log('Secret:', sec ? '***' : 'empty');
    
    if (wsRef.current) {
      console.log('Disconnecting existing WebSocket...');
      wsRef.current.disconnect();
      wsRef.current = null;
    }

    const ws = new EcoWebSocket(url, devId, sec, deviceName);
    wsRef.current = ws;

    ws.onStateChange((newState) => {
      console.log('WebSocket state changed to:', newState);
      setState(newState);
    });
    ws.onLog((log) => setLogs((prev) => [log, ...prev].slice(0, 100)));
    ws.onMessage((message) => {
      console.log('Received message:', message.type);
      if (message.type === 'clipboard.set') {
        const payload = message.payload as ClipboardPayload;
        NativeClipboard.setText(payload.data);
      }
    });

    console.log('Starting WebSocket connection...');
    await ws.connect();
    console.log('WebSocket connected, enabling native listeners...');
    await enableNativeListeners();
    console.log('Connection complete');
  }, [enableNativeListeners]);

  const disconnect = useCallback(() => {
    console.log('Disconnecting...');
    if (wsRef.current) {
      wsRef.current.disconnect();
      wsRef.current = null;
    }
    setState('disconnected');
    disableNativeListeners();
  }, [disableNativeListeners]);

  const sendClipboard = useCallback((text: string): boolean => {
    if (!wsRef.current || state !== 'connected') return false;
    return wsRef.current.send('clipboard.changed', { data: text });
  }, [state]);

  const lastClipboardRef = useRef<string>('');

  useEffect(() => {
    if (state !== 'connected') return;

    const checkClipboard = async () => {
      try {
        const text = await NativeClipboard.getText();
        if (text && text !== lastClipboardRef.current && text.length > 0) {
          lastClipboardRef.current = text;
          console.log('Clipboard changed:', text.substring(0, 50));
          addLog('clipboard.changed', 'sent', { data: text });
          if (wsRef.current) {
            wsRef.current.send('clipboard.changed', { data: text });
          }
        }
      } catch (e) {
        console.error('Clipboard poll error:', e);
      }
    };

    const interval = setInterval(checkClipboard, 1000);
    return () => clearInterval(interval);
  }, [state, addLog]);

  useEffect(() => {
    if (!isNativeAvailable || !ecoEventEmitter) return;

    const callSub = ecoEventEmitter.addListener('EcoCallStateChanged', (event: CallStateEvent) => {
      addLog('call.incoming', 'sent', { number: event.phoneNumber ?? '' });
      if (wsRef.current && state === 'connected') {
        wsRef.current.send('call.incoming', { number: event.phoneNumber ?? '' });
      }
    });

    const notifSub = DeviceEventEmitter.addListener('EcoNotificationPosted', (event: NotificationEvent) => {
      addLog('notification.push', 'sent', { app: event.packageName, title: event.title, body: event.text });
      if (wsRef.current && state === 'connected') {
        wsRef.current.send('notification.push', {
          app: event.packageName,
          title: event.title,
          body: event.text,
        });
      }
    });

    return () => {
      callSub.remove();
      notifSub.remove();
    };
  }, [state, addLog]);

  useEffect(() => {
    const handleAppStateChange = (nextAppState: AppStateStatus) => {
      console.log('AppState changed to:', nextAppState);
      if (nextAppState === 'active' && state === 'disconnected' && wsRef.current) {
        console.log('App became active, attempting to reconnect...');
      }
    };

    const subscription = AppState.addEventListener('change', handleAppStateChange);
    return () => subscription.remove();
  }, [state]);

  return (
    <EcoContext.Provider
      value={{
        state,
        logs,
        connect,
        disconnect,
        sendClipboard,
        serverUrl,
        deviceId,
        secret,
        setConfig,
        nativeEnabled,
        nativeAvailable: isNativeAvailable,
        enableNativeListeners,
        disableNativeListeners,
      }}
    >
      {children}
    </EcoContext.Provider>
  );
}

export function useEco() {
  const context = useContext(EcoContext);
  if (!context) {
    throw new Error('useEco must be used within an EcoProvider');
  }
  return context;
}
