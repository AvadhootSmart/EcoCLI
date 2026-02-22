export type MessageType =
  | 'clipboard.changed'
  | 'clipboard.set'
  | 'notification.push'
  | 'call.incoming'
  | 'call.answer'
  | 'call.hangup'
  | 'device.hello'
  | 'device.ping'
  | 'device.disconnect';

export interface Message<T = unknown> {
  type: MessageType;
  device_id: string;
  secret: string;
  payload: T;
}

export interface ClipboardPayload {
  data: string;
}

export interface NotificationPayload {
  app: string;
  title: string;
  body: string;
}

export interface CallPayload {
  number: string;
}

export interface DevicePayload {
  device_name: string;
}

export type ConnectionState = 'disconnected' | 'connecting' | 'connected' | 'error';

export interface LogEntry {
  id: string;
  timestamp: number;
  type: MessageType;
  direction: 'sent' | 'received';
  payload: unknown;
}

export type ClipboardEvent = {
  text: string;
  timestamp: number;
  source: 'android';
};

export type NotificationEvent = {
  packageName: string;
  title: string;
  text: string;
  timestamp: number;
};

export type CallStateEvent = {
  state: 'ringing' | 'offhook' | 'idle';
  phoneNumber?: string;
  timestamp: number;
};
