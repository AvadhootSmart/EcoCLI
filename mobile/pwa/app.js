const EcoApp = {
  client: null,
  
  elements: {},
  
  clipboardHistory: [],
  
  qrScanner: null,
  qrVideo: null,
  qrCanvas: null,
  qrCanvasCtx: null,
  
  init() {
    this.cacheElements();
    this.loadConfig();
    this.initClient();
    this.bindEvents();
    this.registerServiceWorker();
    this.requestNotificationPermission();
    this.checkQRParams();
    
    document.getElementById('deviceId').textContent = this.client.getDeviceId();
  },
  
  cacheElements() {
    this.elements = {
      serverUrl: document.getElementById('serverUrl'),
      secret: document.getElementById('secret'),
      connectBtn: document.getElementById('connectBtn'),
      scanQrBtn: document.getElementById('scanQrBtn'),
      closeQrBtn: document.getElementById('closeQrBtn'),
      copyBtn: document.getElementById('copyBtn'),
      pasteBtn: document.getElementById('pasteBtn'),
      clipboardPreview: document.getElementById('clipboardPreview'),
      notificationsList: document.getElementById('notificationsList'),
      eventsList: document.getElementById('eventsList'),
      connectionStatus: document.getElementById('connectionStatus'),
      toastContainer: document.getElementById('toastContainer'),
      qrScannerModal: document.getElementById('qrScannerModal'),
      qrVideo: document.getElementById('qrVideo')
    };
  },
  
  loadConfig() {
    const saved = localStorage.getItem('eco_config');
    if (saved) {
      const config = JSON.parse(saved);
      this.elements.serverUrl.value = config.serverUrl || '';
      this.elements.secret.value = config.secret || '';
    }
  },
  
  saveConfig() {
    const config = {
      serverUrl: this.elements.serverUrl.value,
      secret: this.elements.secret.value
    };
    localStorage.setItem('eco_config', JSON.stringify(config));
  },
  
  initClient() {
    this.client = new EcoClient({
      serverUrl: this.elements.serverUrl.value || 'ws://localhost:4949/ws',
      deviceId: localStorage.getItem('eco_device_id'),
      secret: this.elements.secret.value,
      deviceName: 'PWA'
    });
    
    this.setupClientHandlers();
  },
  
  setupClientHandlers() {
    this.client.on('*', (payload, msg) => {
      this.addEvent(msg.type, msg);
    });
    
    this.client.on('clipboard.changed', (payload) => {
      this.handleClipboardChange(payload);
    });
    
    this.client.on('clipboard.set', (payload) => {
      this.handleClipboardSet(payload);
    });
    
    this.client.on('notification.push', (payload) => {
      this.handleNotification(payload);
    });
    
    this.client.on('call.incoming', (payload) => {
      this.handleIncomingCall(payload);
    });
    
    this.client.on('device.hello', () => {
      this.showToast('Connected to server!', 'success');
      this.updateConnectionStatus('connected');
    });
    
    this.client.on('device.disconnect', () => {
      this.showToast('Disconnected from server', 'warning');
      this.updateConnectionStatus('disconnected');
    });
  },
  
  bindEvents() {
    this.elements.connectBtn.addEventListener('click', () => this.toggleConnection());
    
    this.elements.scanQrBtn.addEventListener('click', () => this.openQRScanner());
    this.elements.closeQrBtn.addEventListener('click', () => this.closeQRScanner());
    
    this.elements.qrScannerModal.addEventListener('click', (e) => {
      if (e.target === this.elements.qrScannerModal) {
        this.closeQRScanner();
      }
    });
    
    this.elements.copyBtn.addEventListener('click', () => this.copyFromDevice());
    this.elements.pasteBtn.addEventListener('click', () => this.pasteToDevice());
    
    this.elements.clipboardPreview.addEventListener('click', () => {
      if (this.clipboardHistory.length > 0) {
        this.copyToClipboard(this.clipboardHistory[0]);
      }
    });
  },
  
  async toggleConnection() {
    if (this.client.isConnected()) {
      this.client.disconnect();
      this.updateConnectionStatus('disconnected');
      this.elements.connectBtn.textContent = 'Connect';
      this.elements.copyBtn.disabled = true;
      this.elements.pasteBtn.disabled = true;
    } else {
      this.saveConfig();
      this.client.updateConfig({
        serverUrl: this.elements.serverUrl.value,
        secret: this.elements.secret.value
      });
      
      this.updateConnectionStatus('connecting');
      this.elements.connectBtn.disabled = true;
      this.elements.connectBtn.textContent = 'Connecting...';
      
      try {
        await this.client.connect();
        this.elements.connectBtn.textContent = 'Disconnect';
        this.elements.copyBtn.disabled = false;
        this.elements.pasteBtn.disabled = false;
        this.updateConnectionStatus('connected');
      } catch (error) {
        console.error('Connection failed:', error);
        this.updateConnectionStatus('error');
        this.showToast('Connection failed: ' + error.message, 'error');
        this.elements.connectBtn.textContent = 'Connect';
      } finally {
        this.elements.connectBtn.disabled = false;
      }
    }
  },
  
  updateConnectionStatus(status) {
    const dot = this.elements.connectionStatus.querySelector('.status-dot');
    const text = this.elements.connectionStatus.querySelector('.status-text');
    
    dot.className = 'status-dot';
    
    switch (status) {
      case 'connected':
        dot.classList.add('connected');
        text.textContent = 'Connected';
        break;
      case 'connecting':
        text.textContent = 'Connecting...';
        break;
      case 'error':
        dot.classList.add('error');
        text.textContent = 'Error';
        break;
      default:
        text.textContent = 'Disconnected';
    }
  },
  
  async copyFromDevice() {
    if (this.clipboardHistory.length > 0) {
      await this.copyToClipboard(this.clipboardHistory[0]);
      this.showToast('Copied to clipboard!', 'success');
    }
  },
  
  async pasteToDevice() {
    try {
      const text = await navigator.clipboard.readText();
      if (text) {
        await this.client.setClipboard(text);
        this.showToast('Sent to device!', 'success');
      }
    } catch (error) {
      this.showToast('Failed to read clipboard', 'error');
    }
  },
  
  async copyToClipboard(text) {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch (error) {
      console.error('Clipboard error:', error);
      return false;
    }
  },
  
  handleClipboardChange(payload) {
    if (payload && payload.data) {
      this.clipboardHistory.unshift(payload.data);
      if (this.clipboardHistory.length > 10) {
        this.clipboardHistory.pop();
      }
      this.renderClipboardPreview();
    }
  },
  
  handleClipboardSet(payload) {
    if (payload && payload.data) {
      this.copyToClipboard(payload.data);
      this.showToast('Clipboard updated from server', 'success');
    }
  },
  
  handleNotification(payload) {
    if (!payload) return;
    
    const notification = {
      app: payload.app || 'Unknown',
      title: payload.title || 'Notification',
      body: payload.body || '',
      timestamp: Date.now()
    };
    
    this.renderNotification(notification);
    this.showToast(`${notification.app}: ${notification.title}`, 'success');
    
    if (Notification.permission === 'granted') {
      new Notification(notification.title, {
        body: notification.body,
        icon: '/icons/icon-192.png',
        tag: 'eco-notification'
      });
    }
  },
  
  handleIncomingCall(payload) {
    if (!payload) return;
    
    this.showToast(`Incoming call from ${payload.number || 'unknown'}`, 'warning');
    
    if (Notification.permission === 'granted') {
      new Notification('Incoming Call', {
        body: `From: ${payload.number || 'unknown'}`,
        icon: '/icons/icon-192.png',
        tag: 'eco-call',
        requireInteraction: true
      });
    }
  },
  
  renderClipboardPreview() {
    if (this.clipboardHistory.length === 0) {
      this.elements.clipboardPreview.innerHTML = '<p class="placeholder">No clipboard data</p>';
      return;
    }
    
    const content = this.clipboardHistory[0];
    const truncated = content.length > 500 ? content.substring(0, 500) + '...' : content;
    this.elements.clipboardPreview.innerHTML = `<div class="content">${this.escapeHtml(truncated)}</div>`;
  },
  
  renderNotification(notification) {
    const placeholder = this.elements.notificationsList.querySelector('.placeholder');
    if (placeholder) {
      this.elements.notificationsList.innerHTML = '';
    }
    
    const item = document.createElement('div');
    item.className = 'notification-item';
    item.innerHTML = `
      <div class="app">${this.escapeHtml(notification.app)}</div>
      <div class="title">${this.escapeHtml(notification.title)}</div>
      <div class="body">${this.escapeHtml(notification.body)}</div>
    `;
    
    this.elements.notificationsList.insertBefore(item, this.elements.notificationsList.firstChild);
    
    while (this.elements.notificationsList.children.length > 20) {
      this.elements.notificationsList.removeChild(this.elements.notificationsList.lastChild);
    }
  },
  
  addEvent(type, msg) {
    const placeholder = this.elements.eventsList.querySelector('.placeholder');
    if (placeholder) {
      this.elements.eventsList.innerHTML = '';
    }
    
    const item = document.createElement('div');
    item.className = 'event-item';
    item.innerHTML = `
      <span class="type">${this.escapeHtml(type)}</span>
      <span class="time">${this.formatTime(Date.now())}</span>
    `;
    
    this.elements.eventsList.insertBefore(item, this.elements.eventsList.firstChild);
    
    while (this.elements.eventsList.children.length > 50) {
      this.elements.eventsList.removeChild(this.elements.eventsList.lastChild);
    }
  },
  
  showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.innerHTML = `<span class="toast-message">${this.escapeHtml(message)}</span>`;
    
    this.elements.toastContainer.appendChild(toast);
    
    setTimeout(() => {
      toast.style.opacity = '0';
      setTimeout(() => toast.remove(), 300);
    }, 3000);
  },
  
  formatTime(timestamp) {
    const date = new Date(timestamp);
    return date.toLocaleTimeString();
  },
  
  escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  },
  
  registerServiceWorker() {
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.register('sw.js')
        .then((registration) => {
          console.log('SW registered:', registration.scope);
        })
        .catch((error) => {
          console.log('SW registration failed:', error);
        });
    }
  },
  
  async requestNotificationPermission() {
    if ('Notification' in window && Notification.permission === 'default') {
      const permission = await Notification.requestPermission();
      console.log('Notification permission:', permission);
    }
  },
  
  checkQRParams() {
    const params = new URLSearchParams(window.location.search);
    const server = params.get('server');
    const secret = params.get('secret');
    
    if (server) {
      let wsUrl = server;
      if (!server.startsWith('ws://') && !server.startsWith('wss://')) {
        wsUrl = 'ws://' + server + '/ws';
      }
      this.elements.serverUrl.value = wsUrl;
    }
    
    if (secret) {
      this.elements.secret.value = secret;
    }
    
    // Clean URL
    if (server || secret) {
      window.history.replaceState({}, document.title, '/');
    }
  },
  
  openQRScanner() {
    this.elements.qrScannerModal.classList.add('active');
    this.startQRScanner();
  },
  
  closeQRScanner() {
    this.stopQRScanner();
    this.elements.qrScannerModal.classList.remove('active');
  },
  
  async startQRScanner() {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ 
        video: { facingMode: 'environment' } 
      });
      
      this.elements.qrVideo.srcObject = stream;
      this.elements.qrVideo.setAttribute('playsinline', true);
      
      this.qrScanner = setInterval(() => this.scanQRFrame(), 250);
    } catch (error) {
      console.error('Camera error:', error);
      this.showToast('Camera access denied. Please enable camera permissions.', 'error');
      this.closeQRScanner();
    }
  },
  
  stopQRScanner() {
    if (this.qrScanner) {
      clearInterval(this.qrScanner);
      this.qrScanner = null;
    }
    
    if (this.elements.qrVideo && this.elements.qrVideo.srcObject) {
      const tracks = this.elements.qrVideo.srcObject.getTracks();
      tracks.forEach(track => track.stop());
      this.elements.qrVideo.srcObject = null;
    }
  },
  
  scanQRFrame() {
    if (!this.elements.qrVideo.videoWidth) return;
    
    const canvas = document.createElement('canvas');
    canvas.width = this.elements.qrVideo.videoWidth;
    canvas.height = this.elements.qrVideo.videoHeight;
    
    const ctx = canvas.getContext('2d');
    ctx.drawImage(this.elements.qrVideo, 0, 0, canvas.width, canvas.height);
    
    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    
    if (typeof jsQR === 'undefined') {
      this.showToast('QR scanner library not loaded', 'error');
      return;
    }
    
    const code = jsQR(imageData.data, imageData.width, imageData.height);
    
    if (code) {
      this.handleQRCode(code.data);
    }
  },
  
  handleQRCode(data) {
    console.log('QR Code scanned:', data);
    
    // Try to parse as eco:// URL
    if (data.startsWith('eco://connect?')) {
      const params = new URLSearchParams(data.replace('eco://connect?', ''));
      const server = params.get('server');
      const secret = params.get('secret');
      
      if (server) {
        let wsUrl = server;
        if (!server.startsWith('ws://') && !server.startsWith('wss://')) {
          wsUrl = 'ws://' + server + '/ws';
        }
        this.elements.serverUrl.value = wsUrl;
      }
      
      if (secret) {
        this.elements.secret.value = secret;
      }
      
      this.showToast('QR Code scanned! Tap Connect to pair.', 'success');
      this.closeQRScanner();
      return;
    }
    
    // Try to parse as JSON
    try {
      const parsed = JSON.parse(data);
      if (parsed.serverUrl || parsed.server) {
        this.elements.serverUrl.value = parsed.serverUrl || parsed.server;
      }
      if (parsed.secret) {
        this.elements.secret.value = parsed.secret;
      }
      this.showToast('QR Code scanned! Tap Connect to pair.', 'success');
      this.closeQRScanner();
      return;
    } catch (e) {}
    
    // Plain URL
    if (data.startsWith('http://') || data.startsWith('https://')) {
      let wsUrl = data;
      if (!data.includes('/ws')) {
        wsUrl = data.replace('http://', 'ws://').replace('https://', 'wss://');
        if (!wsUrl.endsWith('/ws')) {
          wsUrl += '/ws';
        }
      }
      this.elements.serverUrl.value = wsUrl;
      this.showToast('QR Code scanned! Tap Connect to pair.', 'success');
      this.closeQRScanner();
    }
  }
};

document.addEventListener('DOMContentLoaded', () => {
  EcoApp.init();
});
