class EcoClient {
  constructor(config = {}) {
    this.serverUrl = config.serverUrl || 'ws://localhost:4949/ws';
    this.deviceId = config.deviceId || this.generateDeviceId();
    this.secret = config.secret || '';
    this.deviceName = config.deviceName || 'PWA';

    this.ws = null;
    this.connected = false;
    this.reconnecting = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 10;
    this.baseReconnectDelay = 1000;
    this.maxReconnectDelay = 60000;

    this.messageQueue = [];
    this.handlers = new Map();

    this.pingInterval = null;
    this.pingTimeout = null;
    this.lastPingTime = 0;

    this.heartbeatInterval = null;
  }

  generateDeviceId() {
    const stored = localStorage.getItem('eco_device_id');
    if (stored) return stored;

    const id = 'pwa_' + Math.random().toString(36).substr(2, 9) + '_' + Date.now();
    localStorage.setItem('eco_device_id', id);
    return id;
  }

  async connect() {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return;
    }

    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(this.serverUrl);

        this.ws.onopen = () => {
          console.log('[Eco] Connected to server');
          this.connected = true;
          this.reconnectAttempts = 0;
          this.reconnecting = false;

          this.sendHello();
          this.flushQueue();
          this.startHeartbeat();

          resolve();
        };

        this.ws.onmessage = (event) => {
          this.handleMessage(event.data);
        };

        this.ws.onclose = (event) => {
          console.log('[Eco] Connection closed', event.code, event.reason);
          this.connected = false;
          this.stopHeartbeat();
          this.handleDisconnect();
        };

        this.ws.onerror = (error) => {
          console.error('[Eco] WebSocket error:', error);
          reject(error);
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  disconnect() {
    this.stopHeartbeat();
    
    if (this.ws) {
      this.send({
        type: 'device.disconnect',
        device_id: this.deviceId,
        secret: this.secret,
        payload: null
      });
      this.ws.close(1000, 'Client disconnect');
      this.ws = null;
    }
    this.connected = false;
  }

  sendHello() {
    this.send({
      type: 'device.hello',
      device_id: this.deviceId,
      secret: this.secret,
      payload: {
        device_name: this.deviceName
      }
    });
  }

  send(data) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      this.messageQueue.push(data);
      return false;
    }

    try {
      this.ws.send(JSON.stringify(data));
      return true;
    } catch (error) {
      console.error('[Eco] Send error:', error);
      this.messageQueue.push(data);
      return false;
    }
  }

  flushQueue() {
    while (this.messageQueue.length > 0) {
      const msg = this.messageQueue.shift();
      this.send(msg);
    }
  }

  handleMessage(data) {
    try {
      const msg = JSON.parse(data);
      console.log('[Eco] Received:', msg.type);

      if (msg.type === 'device.ping') {
        this.sendPong();
        return;
      }

      const handler = this.handlers.get(msg.type);
      if (handler) {
        handler(msg.payload, msg);
      }

      const wildcardHandler = this.handlers.get('*');
      if (wildcardHandler) {
        wildcardHandler(msg.payload, msg);
      }
    } catch (error) {
      console.error('[Eco] Parse error:', error);
    }
  }

  sendPong() {
    this.send({
      type: 'device.ping',
      device_id: this.deviceId,
      secret: this.secret,
      payload: null
    });
  }

  startHeartbeat() {
    this.stopHeartbeat();
    
    this.heartbeatInterval = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        try {
          this.ws.send(JSON.stringify({ type: 'device.ping', device_id: this.deviceId, secret: this.secret, payload: null }));
        } catch (e) {}
      }
    }, 30000);
  }

  stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }

  handleDisconnect() {
    if (this.reconnecting) return;
    this.reconnecting = true;

    const delay = Math.min(
      this.baseReconnectDelay * Math.pow(2, this.reconnectAttempts),
      this.maxReconnectDelay
    );

    console.log(`[Eco] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts + 1})`);

    setTimeout(() => {
      this.reconnectAttempts++;
      
      if (this.reconnectAttempts > this.maxReconnectAttempts) {
        console.error('[Eco] Max reconnect attempts reached');
        this.reconnecting = false;
        return;
      }

      this.connect().catch((err) => {
        console.error('[Eco] Reconnect failed:', err);
      });
    }, delay);
  }

  on(type, handler) {
    this.handlers.set(type, handler);
  }

  off(type) {
    this.handlers.delete(type);
  }

  setClipboard(text) {
    return navigator.clipboard.writeText(text).then(() => {
      this.send({
        type: 'clipboard.changed',
        device_id: this.deviceId,
        secret: this.secret,
        payload: { data: text }
      });
      return true;
    }).catch((err) => {
      console.error('[Eco] Clipboard error:', err);
      return false;
    });
  }

  getClipboard() {
    return navigator.clipboard.readText();
  }

  isConnected() {
    return this.connected;
  }

  getDeviceId() {
    return this.deviceId;
  }

  updateConfig(config) {
    if (config.serverUrl) this.serverUrl = config.serverUrl;
    if (config.secret) this.secret = config.secret;
    if (config.deviceName) this.deviceName = config.deviceName;
  }
}

if (typeof module !== 'undefined' && module.exports) {
  module.exports = EcoClient;
}
