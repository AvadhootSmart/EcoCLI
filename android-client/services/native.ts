import { NativeModules, NativeEventEmitter, Platform } from 'react-native';
import * as Clipboard from 'expo-clipboard';

interface EcoCallInterface {
  startListening(): Promise<boolean>;
  stopListening(): Promise<boolean>;
  answerCall(): Promise<boolean>;
  rejectCall(): Promise<boolean>;
  isListening(): Promise<boolean>;
  addListener(eventName: string): void;
  removeListeners(count: number): void;
}

const { EcoCall } = NativeModules;

console.log('NativeModules.EcoCall:', NativeModules.EcoCall);

export const CallModule = EcoCall as EcoCallInterface | undefined;

console.log('CallModule available:', !!CallModule);

export const ecoEventEmitter = Platform.OS === 'android' && CallModule
  ? new NativeEventEmitter(CallModule)
  : null;

export const isNativeAvailable = Platform.OS === 'android' && !!CallModule;

console.log('isNativeAvailable:', isNativeAvailable);

export const NativeClipboard = {
  setText: async (text: string): Promise<boolean> => {
    try {
      await Clipboard.setStringAsync(text);
      return true;
    } catch (e) {
      console.error('Clipboard setText error:', e);
      return false;
    }
  },
  getText: async (): Promise<string> => {
    try {
      const text = await Clipboard.getStringAsync();
      return text;
    } catch (e) {
      console.error('Clipboard getText error:', e);
      return '';
    }
  },
  startListening: async (): Promise<boolean> => {
    return true;
  },
  stopListening: async (): Promise<boolean> => {
    return true;
  },
  isListening: async (): Promise<boolean> => {
    return true;
  },
  ping: async (): Promise<string> => {
    return 'pong';
  },
};

export const NativeCall = {
  startListening: async (): Promise<boolean> => {
    if (!CallModule) {
      console.log('EcoCall module not available');
      return false;
    }
    return CallModule.startListening();
  },
  stopListening: async (): Promise<boolean> => {
    if (!CallModule) return false;
    return CallModule.stopListening();
  },
  answerCall: async (): Promise<boolean> => {
    if (!CallModule) return false;
    return CallModule.answerCall();
  },
  rejectCall: async (): Promise<boolean> => {
    if (!CallModule) return false;
    return CallModule.rejectCall();
  },
  isListening: async (): Promise<boolean> => {
    if (!CallModule) return false;
    return CallModule.isListening();
  },
};

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
  phoneNumber: string;
  timestamp: number;
};