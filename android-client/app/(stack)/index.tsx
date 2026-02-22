import React, { useState } from "react";
import {
  StyleSheet,
  View,
  TextInput,
  TouchableOpacity,
  Alert,
  Platform,
} from "react-native";
import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { useEco } from "@/context";

const getDefaultUrl = () => {
  if (Platform.OS === "android") {
    return "ws://192.168.1.11:4949/ws";
  }
  return "ws://localhost:4949/ws";
};

export default function HomeScreen() {
  const {
    state,
    connect,
    disconnect,
    serverUrl,
    deviceId,
    secret,
    setConfig,
    nativeAvailable,
  } = useEco();
  const [url, setUrl] = useState(serverUrl || getDefaultUrl());
  const [id, setId] = useState(deviceId);
  const [sec, setSec] = useState(secret);
  const [deviceName] = useState("Eco Android");
  const [error, setError] = useState<string | null>(null);

  const handleConnect = async () => {
    if (!id || !sec) {
      Alert.alert("Error", "Please enter Device ID and Secret");
      return;
    }
    setError(null);
    setConfig(url, id, sec);
    try {
      await connect(url, id, sec, deviceName);
    } catch (e) {
      const message = e instanceof Error ? e.message : "Connection failed";
      setError(message);
      Alert.alert("Connection Failed", message);
    }
  };

  const handleDisconnect = () => {
    disconnect();
  };

  const getStatusColor = () => {
    switch (state) {
      case "connected":
        return "#4CAF50";
      case "connecting":
        return "#FFC107";
      case "error":
        return "#F44336";
      default:
        return "#9E9E9E";
    }
  };

  return (
    <ThemedView style={styles.container}>
      <View style={styles.header}>
        <ThemedText type="title">Eco</ThemedText>
        <View
          style={[styles.statusDot, { backgroundColor: getStatusColor() }]}
        />
      </View>

      <ThemedView style={styles.card}>
        <ThemedText type="subtitle">Connection</ThemedText>

        <View style={styles.inputGroup}>
          <ThemedText style={styles.label}>Server URL</ThemedText>
          <TextInput
            style={styles.input}
            value={url}
            onChangeText={setUrl}
            placeholder="ws://10.0.2.2:4949/ws (emulator) or ws://YOUR_IP:4949/ws (device)"
            placeholderTextColor="#666"
            autoCapitalize="none"
            autoCorrect={false}
            editable={state !== "connected"}
          />
        </View>

        <View style={styles.inputGroup}>
          <ThemedText style={styles.label}>Device ID</ThemedText>
          <TextInput
            style={styles.input}
            value={id}
            onChangeText={setId}
            placeholder="device-id"
            placeholderTextColor="#666"
            autoCapitalize="none"
            autoCorrect={false}
            editable={state !== "connected"}
          />
        </View>

        <View style={styles.inputGroup}>
          <ThemedText style={styles.label}>Secret</ThemedText>
          <TextInput
            style={styles.input}
            value={sec}
            onChangeText={setSec}
            placeholder="shared-secret"
            placeholderTextColor="#666"
            secureTextEntry
            autoCapitalize="none"
            autoCorrect={false}
            editable={state !== "connected"}
          />
        </View>

        <TouchableOpacity
          style={[
            styles.button,
            state === "connected"
              ? styles.disconnectButton
              : styles.connectButton,
          ]}
          onPress={state === "connected" ? handleDisconnect : handleConnect}
        >
          <ThemedText style={styles.buttonText}>
            {state === "connected"
              ? "Disconnect"
              : state === "connecting"
                ? "Connecting..."
                : "Connect"}
          </ThemedText>
        </TouchableOpacity>
      </ThemedView>

      <ThemedView style={styles.statusCard}>
        <ThemedText type="subtitle">Status</ThemedText>
        <ThemedText style={styles.statusText}>State: {state}</ThemedText>
        <ThemedText style={styles.statusText}>
          Native Modules:{" "}
          {nativeAvailable ? "Available" : "Unavailable (use expo run:android)"}
        </ThemedText>
        {error && <ThemedText style={styles.errorText}>{error}</ThemedText>}
      </ThemedView>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    marginBottom: 24,
  },
  statusDot: {
    width: 12,
    height: 12,
    borderRadius: 6,
  },
  card: {
    padding: 16,
    borderRadius: 12,
    marginBottom: 16,
  },
  statusCard: {
    padding: 16,
    borderRadius: 12,
  },
  inputGroup: {
    marginBottom: 16,
  },
  label: {
    marginBottom: 8,
    fontSize: 14,
    opacity: 0.7,
  },
  input: {
    backgroundColor: "#1a1a1a",
    borderRadius: 8,
    padding: 12,
    color: "#fff",
    fontSize: 16,
  },
  button: {
    padding: 16,
    borderRadius: 8,
    alignItems: "center",
    marginTop: 8,
  },
  connectButton: {
    backgroundColor: "#0a7ea4",
  },
  disconnectButton: {
    backgroundColor: "#F44336",
  },
  buttonText: {
    color: "#fff",
    fontWeight: "600",
    fontSize: 16,
  },
  statusText: {
    marginTop: 8,
    fontSize: 14,
    opacity: 0.8,
  },
  errorText: {
    marginTop: 8,
    fontSize: 12,
    color: "#F44336",
  },
});
